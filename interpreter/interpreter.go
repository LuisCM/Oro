// Copyright 2011 The LuisCM. All rights reserved.
// Use of this source code is license that can be found in the LICENSE file.

// Package interpreter implements functions to interpreter.
package interpreter

import (
	"fmt"
	"github.com/luiscm/oro/ast"
	"github.com/luiscm/oro/lexer"
	"github.com/luiscm/oro/parser"
	"github.com/luiscm/oro/rerror"
	"github.com/luiscm/oro/runtime"
	"github.com/luiscm/oro/runtime/stdlib"
	"github.com/luiscm/oro/token"
	"github.com/luiscm/oro/util"
	"io/ioutil"
	"math"
	"path/filepath"
	"strings"
)

type Interpreter struct {
	modules     map[string]*runtime.TModule
	moduleCache map[string]map[string]runtime.Data
	useCache    map[string]runtime.Data
	immutable   map[string]*ast.Identifier
	useStdlib   bool
}

func New() *Interpreter {
	return &Interpreter{
		modules:     map[string]*runtime.TModule{},
		moduleCache: map[string]map[string]runtime.Data{},
		useCache:    map[string]runtime.Data{},
		immutable:   map[string]*ast.Identifier{},
		useStdlib:   false,
	}
}

func (i *Interpreter) Interpreter(ni ast.Node, sc *runtime.Scope) runtime.Data {
	if err := i.useStdLibModules(sc); err != nil {
		i.interpreterError(ni, err.Error())
		return nil
	}
	switch ni := ni.(type) {
	case *ast.Program:
		return i.Program(ni, sc)
	case *ast.Boolean:
		return i.nativeToBoolean(ni.Value)
	case *ast.String:
		return &runtime.TString{Value: ni.Value}
	case *ast.Integer:
		return &runtime.TInteger{Value: ni.Value}
	case *ast.Float:
		return &runtime.TFloat{Value: ni.Value}
	case *ast.Array:
		return i.Array(ni, sc)
	case *ast.Dictionary:
		return i.Dictionary(ni, sc)
	case *ast.Symbol:
		return &runtime.TSymbol{Value: ni.Value}
	case *ast.Nil:
		return &runtime.TNil{}
	case *ast.Val:
		return i.Val(ni, sc)
	case *ast.Var:
		return i.Var(ni, sc)
	case *ast.Is:
		return i.Is(ni, sc)
	case *ast.As:
		return i.As(ni, sc)
	case *ast.If:
		return i.If(ni, sc)
	case *ast.Repeat:
		return i.Repeat(ni, sc)
	case *ast.Match:
		return i.Match(ni, sc)
	case *ast.Break:
		return &runtime.TBreak{}
	case *ast.Continue:
		return &runtime.TContinue{}
	case *ast.Return:
		return &runtime.TReturn{Value: i.Interpreter(ni.Value, sc)}
	case *ast.Pipe:
		return i.Pipe(ni, sc)
	case *ast.PlaceHolder:
		return &runtime.TPlaceHolder{}
	case *ast.Function:
		return &runtime.TFunction{
			Parameters: ni.Parameters,
			Body:       ni.Body,
			ReturnType: ni.ReturnType,
			Variadic:   ni.Variadic,
			Scope:      runtime.NewScopeFrom(sc),
		}
	case *ast.FunctionCall:
		return i.Function(ni, sc)
	case *ast.Module:
		return i.Module(ni, sc)
	case *ast.ModuleAccess:
		return i.ModuleAccess(ni, sc)
	case *ast.Subscript:
		return i.Subscript(ni, sc)
	case *ast.Use:
		return i.Use(ni, sc)
	case *ast.Assign:
		return i.Assign(ni, sc)
	case *ast.BlockStatement:
		return i.BlockStatement(ni, sc)
	case *ast.ExpressionStatement:
		return i.Interpreter(ni.Expression, sc)
	case *ast.Identifier:
		return i.Identifier(ni, sc)
	case *ast.PrefixExpression:
		return i.PrefixExpression(ni, sc)
	case *ast.InfixExpression:
		return i.InfixExpression(ni, sc)
	}
	return nil
}

func (i *Interpreter) useStdLibModules(sc *runtime.Scope) error {
	if i.useStdlib {
		return nil
	}
	i.useStdlib = true
	for _, module := range stdlib.Modules {
		lex := lexer.New([]byte(module))
		if rerror.HasErrors() {
			return rerror.ErrorFmt("Problem reading Standard Library module")
		}
		parse := parser.New(lex)
		program := parse.Parse()
		if rerror.HasErrors() {
			return rerror.ErrorFmt("Problem parsing Standard Library module")
		}
		i.Interpreter(program, sc)
	}
	return nil
}

func (i *Interpreter) Program(np *ast.Program, sc *runtime.Scope) runtime.Data {
	var result runtime.Data
	for _, statement := range np.Statements {
		result = i.Interpreter(statement, sc)
	}
	return result
}

func (i *Interpreter) Val(nl *ast.Val, sc *runtime.Scope) runtime.Data {
	data := i.Interpreter(nl.Value, sc)
	if data == nil {
		return nil
	}
	if _, ok := sc.Read(nl.Name.Value); ok {
		i.interpreterError(nl, fmt.Sprintf("Identifier '%s' already declared", nl.Name.Value))
		return nil
	}
	sc.Write(nl.Name.Value, data)
	if _, ok := i.immutable[nl.Name.Value]; !ok {
		i.immutable[nl.Name.Value] = nl.Name
	}
	return data
}

func (i *Interpreter) Var(nv *ast.Var, sc *runtime.Scope) runtime.Data {
	data := i.Interpreter(nv.Value, sc)
	if data == nil {
		return nil
	}
	if _, ok := sc.Read(nv.Name.Value); ok {
		i.interpreterError(nv, fmt.Sprintf("Identifier '%s' already declared", nv.Name.Value))
		return nil
	}
	sc.Write(nv.Name.Value, data)
	return data
}

func (i *Interpreter) Module(nm *ast.Module, sc *runtime.Scope) runtime.Data {
	if _, ok := i.modules[nm.Name.Value]; ok {
		i.interpreterError(nm, fmt.Sprintf("Module '%s' redeclared", nm.Name.Value))
	} else {
		i.modules[nm.Name.Value] = &runtime.TModule{Name: nm.Name, Body: nm.Body}
	}
	return nil
}

func (i *Interpreter) ModuleAccess(na *ast.ModuleAccess, sc *runtime.Scope) runtime.Data {
	if module, ok := i.modules[na.Object.Value]; ok {
		if results, ok := i.moduleCache[module.Name.Value]; ok {
			if result, ok := results[na.Parameter.Value]; ok {
				return result
			}
		} else {
			newScope := runtime.NewScopeFrom(sc)
			results := map[string]runtime.Data{}
			for _, statement := range module.Body.Statements {
				switch sType := statement.(type) {
				case *ast.ExpressionStatement:
					switch eType := sType.Expression.(type) {
					case *ast.Val:
						result := i.Interpreter(statement, newScope)
						if result == nil {
							return nil
						}
						results[eType.Name.Value] = result
					default:
						i.interpreterError(na, "Only Val statements are accepted as Module members")
						return nil
					}
				default:
					i.interpreterError(na, "Only Val statements are accepted as Module members")
					return nil
				}
			}
			i.moduleCache[module.Name.Value] = results
			if val, ok := i.moduleCache[module.Name.Value][na.Parameter.Value]; ok {
				return val
			} else {
				i.interpreterError(na, fmt.Sprintf("Member '%s' in module '%s' not found", na.Parameter.Value, na.Object.Value))
				return nil
			}
		}
	}
	i.interpreterError(na, fmt.Sprintf("%s.%s not found", na.Object.Value, na.Parameter.Value))
	return nil
}

func (i *Interpreter) Identifier(ni *ast.Identifier, sc *runtime.Scope) runtime.Data {
	if data, ok := sc.Read(ni.Value); ok {
		return data
	}
	i.interpreterError(ni, fmt.Sprintf("Identifier '%s' not found in current memory", ni.Value))
	return nil
}

func (i *Interpreter) BlockStatement(nb *ast.BlockStatement, sc *runtime.Scope) runtime.Data {
	var result runtime.Data
	for _, statement := range nb.Statements {
		result = i.Interpreter(statement, sc)
		if result == nil {
			return nil
		}
		if i.shouldBreakImmediately(result) {
			return result
		}
	}
	return result
}

func (i *Interpreter) Assign(na *ast.Assign, sc *runtime.Scope) runtime.Data {
	var original runtime.Data
	var name string
	var ok bool
	var err error
	switch naType := na.Name.(type) {
	case *ast.Identifier:
		name = naType.Value
	case *ast.Subscript:
		name = naType.Left.(*ast.Identifier).Value
	}
	if original, ok = sc.Read(name); !ok {
		i.interpreterError(na, fmt.Sprintf("Identifier '%s' not found in current memory", name))
		return nil
	}
	if _, ok = i.immutable[name]; ok {
		i.interpreterError(na, fmt.Sprintf("Identifier '%s' is immutable", name))
		return nil
	}
	data := i.Interpreter(na.Right, sc)
	if data == nil {
		return nil
	}
	switch naType := na.Name.(type) {
	case *ast.Subscript:
		data, err = i.AssignSubscript(naType, original, data, sc)
		if err != nil {
			i.interpreterError(na, err.Error())
			return nil
		}
	}
	if data.Type() != original.Type() {
		i.interpreterError(na, fmt.Sprintf("Variable assignment should keep the original data type '%s'", original.Type()))
		return nil
	}
	sc.Update(name, data)
	return data
}

func (i *Interpreter) AssignSubscript(ns *ast.Subscript, original runtime.Data, value runtime.Data, sc *runtime.Scope) (runtime.Data, error) {
	index := i.Interpreter(ns.Index, sc)
	if index == nil {
		return nil, nil
	}
	switch {
	case original.Type() == runtime.TTString && index.Type() == runtime.TTInteger && value.Type() == runtime.TTString:
		str := original.(*runtime.TString)
		idx := index.(*runtime.TInteger).Value
		value := value.(*runtime.TString).Value
		idx, err := i.checkStringBounds(str.Value, idx)
		if err != nil {
			return nil, err
		}
		return &runtime.TString{Value: str.Value[:idx] + value + str.Value[idx+1:]}, nil
	case original.Type() == runtime.TTArray && index.Type() == runtime.TTInteger || index.Type() == runtime.TTPlaceHolder:
		array := original.(*runtime.TArray)
		if index.Type() == runtime.TTInteger {
			idx := index.(*runtime.TInteger).Value
			idx, err := i.checkArrayBounds(array.Elements, idx)
			if err != nil {
				return nil, err
			}
			array.Elements[idx] = value
		} else {
			array.Elements = append(array.Elements, value)
		}
		return array, nil
	case original.Type() == runtime.TTDictionary:
		dictionary := original.(*runtime.TDictionary)
		foundDict := false
		for k := range dictionary.Pairs {
			if k.Check() == index.Check() {
				dictionary.Pairs[k] = value
				foundDict = true
				break
			}
		}
		if !foundDict {
			dictionary.Pairs[index] = value
		}
		return dictionary, nil
	default:
		return nil, rerror.ErrorFmt("Subscript assignment not recognized")
	}
}

func (i *Interpreter) Array(na *ast.Array, sc *runtime.Scope) runtime.Data {
	var result []runtime.Data
	for _, element := range na.List.Elements {
		value := i.Interpreter(element, sc)
		result = append(result, value)
	}
	return &runtime.TArray{Elements: result}
}

func (i *Interpreter) Dictionary(nd *ast.Dictionary, sc *runtime.Scope) runtime.Data {
	result := map[runtime.Data]runtime.Data{}
	for k, v := range nd.Pairs {
		key := i.Interpreter(k, sc)
		if key == nil {
			return nil
		}
		value := i.Interpreter(v, sc)
		result[key] = value
	}
	return &runtime.TDictionary{Pairs: result}
}

func (i *Interpreter) If(ni *ast.If, sc *runtime.Scope) runtime.Data {
	condition := i.Interpreter(ni.Condition, sc)
	if i.isType(condition) {
		return i.Interpreter(ni.Then, runtime.NewScopeFrom(sc))
	} else if ni.Else != nil {
		return i.Interpreter(ni.Else, runtime.NewScopeFrom(sc))
	} else {
		return runtime.Nil
	}
}

func (i *Interpreter) Match(nm *ast.Match, sc *runtime.Scope) runtime.Data {
	var control runtime.Data
	if nm.Control == nil {
		control = runtime.Yes
	} else {
		control = i.Interpreter(nm.Control, sc)
		if control == nil {
			i.interpreterError(nm, "Match control expression couldn't be interpreted")
			return nil
		}
	}
	theWhen, err := i.MatchWhen(nm.Whens, control, sc)
	if err != nil {
		i.interpreterError(nm, err.Error())
		return nil
	}
	if theWhen != nil {
		return i.Interpreter(theWhen.Body, runtime.NewScopeFrom(sc))
	}
	if nm.Else != nil {
		return i.Interpreter(nm.Else, runtime.NewScopeFrom(sc))
	}
	return nil
}

func (i *Interpreter) MatchWhen(whens []*ast.MatchWhen, control runtime.Data, sc *runtime.Scope) (*ast.MatchWhen, error) {
	for _, ws := range whens {
		matches := 0
		for index, element := range ws.Values.Elements {
			parameter := i.Interpreter(element, sc)
			switch {
			case parameter.Type() == control.Type():
				if parameter.Check() == control.Check() {
					return ws, nil
				}
			case control.Type() == runtime.TTArray:
				arrayData := control.(*runtime.TArray).Elements
				if len(ws.Values.Elements) != len(arrayData) {
					break
				}
				if parameter.Type() == arrayData[index].Type() && parameter.Check() == arrayData[index].Check() ||
					parameter.Type() == runtime.TTPlaceHolder {
					matches++
					if matches == len(arrayData) {
						return ws, nil
					}
				}
			case parameter.Type() == runtime.TTSymbol && control.Type() == runtime.TTString:
				if parameter.(*runtime.TSymbol).Value == control.(*runtime.TString).Value {
					return ws, nil
				}
			default:
				return nil, rerror.ErrorFmt("Type '%s' can't be used in a match when with control type '%s'", parameter.Type(), control.Type())
			}
		}
	}
	return nil, nil
}

func (i *Interpreter) Repeat(nr *ast.Repeat, sc *runtime.Scope) runtime.Data {
	if nr.Enumerable == nil {
		return i.RepeatInfinite(nr, sc)
	}
	enumData := i.Interpreter(nr.Enumerable, sc)
	if enumData == nil {
		return nil
	}
	switch enum := enumData.(type) {
	case *runtime.TString:
		return i.ForArray(nr, i.stringToArray(enum), sc)
	case *runtime.TArray:
		return i.ForArray(nr, enum, sc)
	case *runtime.TDictionary:
		return i.ForDictionary(nr, enum, sc)
	case *runtime.TSymbol:
		str := &runtime.TString{Value: enum.Value}
		return i.ForArray(nr, i.stringToArray(str), sc)
	default:
		i.interpreterError(nr, fmt.Sprintf("Type %s is not an enumerable", enumData.Type()))
		return nil
	}
}

func (i *Interpreter) RepeatInfinite(nr *ast.Repeat, sc *runtime.Scope) runtime.Data {
	var out []runtime.Data
	for {
		newScope := runtime.NewScopeFrom(sc)
		result := i.Interpreter(nr.Body, newScope)
		if result == nil {
			return nil
		}
		if result.Type() == runtime.TTBreak {
			break
		} else if result.Type() == runtime.TTContinue {
			continue
		} else if result.Type() == runtime.TTReturn {
			return result
		}
		out = append(out, result)
	}
	return &runtime.TArray{Elements: out}
}

func (i *Interpreter) ForArray(nr *ast.Repeat, array *runtime.TArray, sc *runtime.Scope) runtime.Data {
	var out []runtime.Data
	for index, value := range array.Elements {
		newScope := runtime.NewScopeFrom(sc)
		switch len(nr.Arguments.Elements) {
		case 1:
			sc.Write(nr.Arguments.Elements[0].Value, value)
		case 2:
			sc.Write(nr.Arguments.Elements[0].Value, &runtime.TInteger{Value: int64(index)})
			sc.Write(nr.Arguments.Elements[1].Value, value)
		default:
			i.interpreterError(nr, "A for loop with an array expects at most 2 arguments")
			return nil
		}
		result := i.Interpreter(nr.Body, newScope)
		if result == nil {
			return nil
		}
		if result.Type() == runtime.TTBreak {
			break
		} else if result.Type() == runtime.TTContinue {
			continue
		} else if result.Type() == runtime.TTReturn {
			return result
		}
		out = append(out, result)
	}
	return &runtime.TArray{Elements: out}
}

func (i *Interpreter) ForDictionary(nr *ast.Repeat, dictionary *runtime.TDictionary, sc *runtime.Scope) runtime.Data {
	var out []runtime.Data
	for pair, value := range dictionary.Pairs {
		newScope := runtime.NewScopeFrom(sc)
		switch len(nr.Arguments.Elements) {
		case 1:
			sc.Write(nr.Arguments.Elements[0].Value, value)
		case 2:
			sc.Write(nr.Arguments.Elements[0].Value, pair)
			sc.Write(nr.Arguments.Elements[1].Value, value)
		default:
			i.interpreterError(nr, "A for loop with a dictionary expects at most 2 arguments")
			return nil
		}
		result := i.Interpreter(nr.Body, newScope)
		if result == nil {
			return nil
		}
		if result.Type() == runtime.TTBreak {
			break
		} else if result.Type() == runtime.TTContinue {
			continue
		} else if result.Type() == runtime.TTReturn {
			return result
		}
		out = append(out, result)
	}
	return &runtime.TArray{Elements: out}
}

func (i *Interpreter) Function(nf *ast.FunctionCall, sc *runtime.Scope) runtime.Data {
	switch nfType := nf.Function.(type) {
	case *ast.Identifier:
		if runtimeFn, ok := runtime.FnRuntime[nfType.Value]; ok {
			return i.RuntimeFunction(nf, runtimeFn, sc)
		}
	}
	fn := i.Interpreter(nf.Function, sc)
	if fn == nil {
		return nil
	}
	if fn.Type() != runtime.TTFunction {
		i.interpreterError(nf, "Trying to call a non-function")
		return nil
	}
	function := fn.(*runtime.TFunction)
	fnScope := runtime.NewScopeFrom(function.Scope)
	if !function.Variadic {
		if len(nf.Arguments.Elements) > len(function.Parameters) {
			i.interpreterError(nf, "Too many arguments in function call")
			return nil
		}
	}
	defaultCount := 0
	for _, param := range function.Parameters {
		if param.Default != nil {
			value := i.Interpreter(param.Default, sc)
			if value == nil {
				return nil
			}
			if param.Type != nil {
				if err := i.checkTypeMatch(value.Type(), param.Type.Value); err != nil {
					i.interpreterError(nf, err.Error())
					return nil
				}
			}
			fnScope.Write(param.Name.Value, value)
			defaultCount++
		}
	}
	if len(nf.Arguments.Elements) < len(function.Parameters)-defaultCount {
		i.interpreterError(nf, "Too few arguments in function call")
		return nil
	}
	var arguments []runtime.Data
	countParams := len(function.Parameters) - 1
	for index, element := range nf.Arguments.Elements {
		value := i.Interpreter(element, sc)
		if value == nil {
			return nil
		}
		var paramName *ast.Identifier
		var paramType *ast.Identifier
		if function.Variadic && index >= countParams {
			paramName = function.Parameters[countParams].Name
			paramType = function.Parameters[countParams].Type
		} else {
			paramName = function.Parameters[index].Name
			paramType = function.Parameters[index].Type
		}
		if paramType != nil {
			if err := i.checkTypeMatch(value.Type(), paramType.Value); err != nil {
				i.interpreterError(nf, err.Error())
				return nil
			}
		}
		if function.Variadic && index >= countParams {
			arguments = append(arguments, value)
		} else {
			fnScope.Write(paramName.Value, value)
		}
	}
	if function.Variadic && len(arguments) > 0 {
		fnScope.Write(function.Parameters[len(function.Parameters)-1].Name.Value, &runtime.TArray{Elements: arguments})
	}
	result := i.unwrapReturnValue(i.Interpreter(function.Body, fnScope))
	if result == nil {
		return nil
	}
	if function.ReturnType != nil {
		if err := i.checkTypeMatch(result.Type(), function.ReturnType.Value); err != nil {
			i.interpreterError(nf, err.Error())
			return nil
		}
	}
	return result
}

func (i *Interpreter) RuntimeFunction(nf *ast.FunctionCall, fn runtime.TRuntimeFn, sc *runtime.Scope) runtime.Data {
	var args []runtime.Data
	for _, element := range nf.Arguments.Elements {
		value := i.Interpreter(element, sc)
		if value != nil {
			args = append(args, value)
		}
	}
	data, err := fn(args...)
	if err != nil {
		i.interpreterError(nf, err.Error())
		return nil
	}
	return data
}

func (i *Interpreter) Subscript(ns *ast.Subscript, sc *runtime.Scope) runtime.Data {
	left := i.Interpreter(ns.Left, sc)
	index := i.Interpreter(ns.Index, sc)
	if left == nil || index == nil {
		return nil
	}
	switch {
	case left.Type() == runtime.TTString && index.Type() == runtime.TTInteger:
		result, err := i.StringSubscript(left, index)
		if err != nil {
			i.interpreterError(ns, err.Error())
		}
		return result
	case left.Type() == runtime.TTArray && index.Type() == runtime.TTInteger:
		return i.ArraySubscript(left, index)
	case left.Type() == runtime.TTDictionary:
		return i.DictionarySubscript(left, index)
	default:
		i.interpreterError(ns, fmt.Sprintf("Subscript on '%s' not supported with literal '%s'", left.Type(), index.Type()))
		return nil
	}
}

func (i *Interpreter) ArraySubscript(array, index runtime.Data) runtime.Data {
	arrayData := array.(*runtime.TArray).Elements
	idx := index.(*runtime.TInteger).Value
	idx, err := i.checkArrayBounds(arrayData, idx)
	if err != nil {
		return runtime.Nil
	}
	return arrayData[idx]
}

func (i *Interpreter) DictionarySubscript(dictionary, index runtime.Data) runtime.Data {
	dictionaryData := dictionary.(*runtime.TDictionary).Pairs
	for k, v := range dictionaryData {
		if k.Check() == index.Check() {
			return v
		}
	}
	return runtime.Nil
}

func (i *Interpreter) StringSubscript(str, index runtime.Data) (runtime.Data, error) {
	stringData := str.(*runtime.TString).Value
	idx := index.(*runtime.TInteger).Value
	idx, err := i.checkStringBounds(stringData, idx)
	if err != nil {
		return runtime.Nil, nil
	}
	return &runtime.TString{Value: string(stringData[idx])}, nil
}

func (i *Interpreter) Pipe(np *ast.Pipe, sc *runtime.Scope) runtime.Data {
	switch rightFunction := np.Right.(type) {
	case *ast.FunctionCall:
		rightFunction.Arguments.Elements = append([]ast.Expression{np.Left}, rightFunction.Arguments.Elements...)
		return i.Interpreter(rightFunction, sc)
	default:
		i.interpreterError(np, "Pipe operator expects a function on the right side")
		return nil
	}
}

func (i *Interpreter) Use(nu *ast.Use, sc *runtime.Scope) runtime.Data {
	fileName := i.checkExtFileName(nu.File.Value)
	if cache, ok := i.useCache[fileName]; ok {
		return cache
	}
	source, err := ioutil.ReadFile(fileName)
	if err != nil {
		i.interpreterError(nu, fmt.Sprintf("Couldn't read imported file '%s'", nu.File.Value))
		return nil
	}
	lex := lexer.New(source)
	if rerror.HasErrors() {
		return nil
	}
	parse := parser.New(lex)
	program := parse.Parse()
	if rerror.HasErrors() {
		return nil
	}
	result := i.Interpreter(program, sc)
	i.useCache[fileName] = result
	return result
}

func (i *Interpreter) Is(ni *ast.Is, sc *runtime.Scope) runtime.Data {
	dataIs := i.Interpreter(ni.Left, sc)
	if dataIs == nil {
		return nil
	}
	if !i.checkSupportedType(ni.Right.Value) {
		i.interpreterError(ni, fmt.Sprintf("Unknown type '%s' in is operator", ni.Right.Value))
		return nil
	}
	if dataIs.Type() == ni.Right.Value {
		return runtime.Yes
	}
	return runtime.No
}

func (i *Interpreter) As(na *ast.As, sc *runtime.Scope) runtime.Data {
	if !i.checkSupportedType(na.Right.Value) {
		i.interpreterError(na, fmt.Sprintf("Unknown type '%s' in as operator", na.Right.Value))
		return nil
	}
	nf := &ast.FunctionCall{
		Token:     na.Token,
		Arguments: &ast.ExpressionList{Elements: []ast.Expression{na.Left}},
	}
	switch na.Right.Value {
	case runtime.TTString:
		return i.RuntimeFunction(nf, runtime.FnRuntime[runtime.TTString], sc)
	case runtime.TTInteger:
		return i.RuntimeFunction(nf, runtime.FnRuntime[runtime.TTInteger], sc)
	case runtime.TTFloat:
		return i.RuntimeFunction(nf, runtime.FnRuntime[runtime.TTFloat], sc)
	case runtime.TTArray:
		return i.RuntimeFunction(nf, runtime.FnRuntime[runtime.TTArray], sc)
	default:
		i.interpreterError(na, fmt.Sprintf("Can't convert to type '%s'", na.Right.Value))
		return nil
	}
}

func (i *Interpreter) PrefixExpression(np *ast.PrefixExpression, sc *runtime.Scope) runtime.Data {
	data := i.Interpreter(np.Right, sc)
	if data == nil {
		i.interpreterError(np, fmt.Sprintf("Trying to run operator '%s' with an unknown value", np.Operator))
		return nil
	}
	var out runtime.Data
	var err error
	switch np.Operator {
	case string(token.Not):
		out = i.nativeToBoolean(!i.isType(data))
	case string(token.Minus):
		out, err = i.MinusPrefix(data)
	case string(token.BitwiseNot):
		out, err = i.BitwiseNotPrefix(data)
	default:
		err = rerror.ErrorFmt("Unsupported prefix operator")
	}
	if err != nil {
		i.interpreterError(np, err.Error())
	}
	return out
}

func (i *Interpreter) MinusPrefix(data runtime.Data) (runtime.Data, error) {
	switch data.Type() {
	case runtime.TTInteger:
		return &runtime.TInteger{Value: -data.(*runtime.TInteger).Value}, nil
	case runtime.TTFloat:
		return &runtime.TFloat{Value: -data.(*runtime.TFloat).Value}, nil
	default:
		return nil, rerror.ErrorFmt("Minus prefix can be applied to Integers and Floats only")
	}
}

func (i *Interpreter) BitwiseNotPrefix(data runtime.Data) (runtime.Data, error) {
	switch data.Type() {
	case runtime.TTInteger:
		return &runtime.TInteger{Value: ^data.(*runtime.TInteger).Value}, nil
	default:
		return nil, rerror.ErrorFmt("Bitwise not prefix can be applied to Integers only")
	}
}

func (i *Interpreter) InfixExpression(ni *ast.InfixExpression, sc *runtime.Scope) runtime.Data {
	left := i.Interpreter(ni.Left, sc)
	if ni.Operator == token.LogicalAnd && !i.isType(left) {
		return runtime.No
	}
	if ni.Operator == token.LogicalOr && i.isType(left) {
		return runtime.Yes
	}
	right := i.Interpreter(ni.Right, sc)
	if left == nil || right == nil {
		return nil
	}
	var out runtime.Data
	var err error
	switch {
	case left.Type() == runtime.TTBoolean && right.Type() == runtime.TTString:
		l := &runtime.TString{Value: fmt.Sprintf("%t", left.(*runtime.TBoolean).Value)}
		r := &runtime.TString{Value: fmt.Sprintf("%s", right.(*runtime.TString).Value)}
		out, err = i.StringInfix(ni.Operator, l.Value, r.Value)
	case left.Type() == runtime.TTString && right.Type() == runtime.TTBoolean:
		l := &runtime.TString{Value: fmt.Sprintf("%s", left.(*runtime.TString).Value)}
		r := &runtime.TString{Value: fmt.Sprintf("%t", right.(*runtime.TBoolean).Value)}
		out, err = i.StringInfix(ni.Operator, l.Value, r.Value)
	case left.Type() == runtime.TTBoolean && right.Type() == runtime.TTBoolean:
		out, err = i.BooleanInfix(ni.Operator, left, right)
	case left.Type() == runtime.TTString && right.Type() == runtime.TTString:
		out, err = i.StringInfix(ni.Operator, left.(*runtime.TString).Value, right.(*runtime.TString).Value)
	case left.Type() == runtime.TTString && right.Type() == runtime.TTSymbol:
		out, err = i.StringInfix(ni.Operator, left.(*runtime.TString).Value, right.(*runtime.TSymbol).Value)
	case left.Type() == runtime.TTInteger && right.Type() == runtime.TTString:
		l := &runtime.TString{Value: fmt.Sprintf("%d", left.(*runtime.TInteger).Value)}
		r := &runtime.TString{Value: fmt.Sprintf("%s", right.(*runtime.TString).Value)}
		out, err = i.StringInfix(ni.Operator, l.Value, r.Value)
	case left.Type() == runtime.TTString && right.Type() == runtime.TTInteger:
		l := &runtime.TString{Value: fmt.Sprintf("%s", left.(*runtime.TString).Value)}
		r := &runtime.TString{Value: fmt.Sprintf("%d", right.(*runtime.TInteger).Value)}
		out, err = i.StringInfix(ni.Operator, l.Value, r.Value)
	case left.Type() == runtime.TTInteger && right.Type() == runtime.TTInteger:
		out, err = i.IntegerInfix(ni.Operator, left, right)
	case left.Type() == runtime.TTInteger && right.Type() == runtime.TTFloat:
		out, err = i.FloatInfix(ni.Operator, float64(left.(*runtime.TInteger).Value), right.(*runtime.TFloat).Value)
	case left.Type() == runtime.TTFloat && right.Type() == runtime.TTString:
		l := &runtime.TString{Value: fmt.Sprintf("%f", left.(*runtime.TFloat).Value)}
		r := &runtime.TString{Value: fmt.Sprintf("%s", right.(*runtime.TString).Value)}
		out, err = i.StringInfix(ni.Operator, l.Value, r.Value)
	case left.Type() == runtime.TTString && right.Type() == runtime.TTFloat:
		l := &runtime.TString{Value: fmt.Sprintf("%s", left.(*runtime.TString).Value)}
		r := &runtime.TString{Value: fmt.Sprintf("%f", right.(*runtime.TFloat).Value)}
		out, err = i.StringInfix(ni.Operator, l.Value, r.Value)
	case left.Type() == runtime.TTFloat && right.Type() == runtime.TTFloat:
		out, err = i.FloatInfix(ni.Operator, left.(*runtime.TFloat).Value, right.(*runtime.TFloat).Value)
	case left.Type() == runtime.TTFloat && right.Type() == runtime.TTInteger:
		out, err = i.FloatInfix(ni.Operator, left.(*runtime.TFloat).Value, float64(right.(*runtime.TInteger).Value))
	case left.Type() == runtime.TTArray && right.Type() == runtime.TTArray:
		out, err = i.ArrayInfix(ni.Operator, left, right)
	case left.Type() == runtime.TTDictionary && right.Type() == runtime.TTDictionary:
		out, err = i.DictionaryInfix(ni.Operator, left, right)
	case left.Type() == runtime.TTSymbol && right.Type() == runtime.TTSymbol:
		out, err = i.StringInfix(ni.Operator, left.(*runtime.TSymbol).Value, right.(*runtime.TSymbol).Value)
	case left.Type() == runtime.TTSymbol && right.Type() == runtime.TTString:
		out, err = i.StringInfix(ni.Operator, left.(*runtime.TSymbol).Value, right.(*runtime.TString).Value)
	case left.Type() == runtime.TTNil || right.Type() == runtime.TTNil:
		out, err = i.NilInfix(ni.Operator, left, right)
	case left.Type() != right.Type():
		err = rerror.ErrorFmt("Cannot run expression with types '%s' and '%s'", left.Type(), right.Type())
	default:
		err = rerror.ErrorFmt("Unknown operator %s for types '%s' and '%s'", ni.Operator, left.Type(), right.Type())
	}
	if err != nil {
		i.interpreterError(ni, err.Error())
	}
	return out
}

func (i *Interpreter) IntegerInfix(operator string, left, right runtime.Data) (runtime.Data, error) {
	leftVal := left.(*runtime.TInteger).Value
	rightVal := right.(*runtime.TInteger).Value
	switch operator {
	case string(token.Plus):
		return &runtime.TInteger{Value: leftVal + rightVal}, nil
	case string(token.Minus):
		return &runtime.TInteger{Value: leftVal - rightVal}, nil
	case string(token.Multiply):
		return &runtime.TInteger{Value: leftVal * rightVal}, nil
	case string(token.Divide):
		if rightVal == 0 {
			return nil, rerror.ErrorFmt("Division by 0")
		}
		value := float64(leftVal) / float64(rightVal)
		if math.Trunc(value) == value {
			return &runtime.TInteger{Value: int64(value)}, nil
		}
		return &runtime.TFloat{Value: value}, nil
	case string(token.Modulus):
		return &runtime.TInteger{Value: leftVal % rightVal}, nil
	case token.Exponential:
		return &runtime.TInteger{Value: int64(math.Pow(float64(leftVal), float64(rightVal)))}, nil
	case string(token.Less):
		return i.nativeToBoolean(leftVal < rightVal), nil
	case token.LessEqual:
		return i.nativeToBoolean(leftVal <= rightVal), nil
	case string(token.Greater):
		return i.nativeToBoolean(leftVal > rightVal), nil
	case token.GreaterEqual:
		return i.nativeToBoolean(leftVal >= rightVal), nil
	case token.BitShiftLeft:
		if leftVal < 0 || rightVal < 0 {
			return nil, rerror.ErrorFmt("Bitwise shift requires two unsigned Integers")
		}
		return &runtime.TInteger{Value: int64(uint64(leftVal) << uint64(rightVal))}, nil
	case token.BitShiftRight:
		if leftVal < 0 || rightVal < 0 {
			return nil, rerror.ErrorFmt("Bitwsise shift requires two unsigned Integers")
		}
		return &runtime.TInteger{Value: int64(uint64(leftVal) >> uint64(rightVal))}, nil
	case string(token.BitwiseAnd):
		return &runtime.TInteger{Value: leftVal & rightVal}, nil
	case string(token.BitwiseOr):
		return &runtime.TInteger{Value: leftVal | rightVal}, nil
	case token.Equal:
		return i.nativeToBoolean(leftVal == rightVal), nil
	case token.NotEqual:
		return i.nativeToBoolean(leftVal != rightVal), nil
	case token.Range:
		return i.RangeIntegerInfix(leftVal, rightVal), nil
	default:
		return nil, rerror.ErrorFmt("Unsupported Integer operator '%s'", operator)
	}
}

func (i *Interpreter) FloatInfix(operator string, left, right float64) (runtime.Data, error) {
	switch operator {
	case string(token.Plus):
		return &runtime.TFloat{Value: left + right}, nil
	case string(token.Minus):
		return &runtime.TFloat{Value: left - right}, nil
	case string(token.Multiply):
		return &runtime.TFloat{Value: left * right}, nil
	case string(token.Divide):
		if right == 0 {
			return nil, rerror.ErrorFmt("Division by 0")
		}
		return &runtime.TFloat{Value: left / right}, nil
	case string(token.Modulus):
		return &runtime.TFloat{Value: math.Mod(left, right)}, nil
	case token.Exponential:
		return &runtime.TFloat{Value: math.Pow(left, right)}, nil
	case string(token.Less):
		return i.nativeToBoolean(left < right), nil
	case token.LessEqual:
		return i.nativeToBoolean(left <= right), nil
	case string(token.Greater):
		return i.nativeToBoolean(left > right), nil
	case token.GreaterEqual:
		return i.nativeToBoolean(left >= right), nil
	case token.Equal:
		return i.nativeToBoolean(left == right), nil
	case token.NotEqual:
		return i.nativeToBoolean(left != right), nil
	default:
		return nil, rerror.ErrorFmt("Unsupported Float operator '%s'", operator)
	}
}

func (i *Interpreter) StringInfix(operator string, left, right string) (runtime.Data, error) {
	switch operator {
	case string(token.Plus):
		return &runtime.TString{Value: left + right}, nil
	case string(token.Less):
		return i.nativeToBoolean(len(left) < len(right)), nil
	case token.LessEqual:
		return i.nativeToBoolean(len(left) <= len(right)), nil
	case string(token.Greater):
		return i.nativeToBoolean(len(left) > len(right)), nil
	case token.GreaterEqual:
		return i.nativeToBoolean(len(left) >= len(right)), nil
	case token.Equal:
		return i.nativeToBoolean(left == right), nil
	case token.NotEqual:
		return i.nativeToBoolean(left != right), nil
	case token.Range:
		return i.RangeStringInfix(left, right)
	default:
		return nil, rerror.ErrorFmt("Unsupported String operator '%s'", operator)
	}
}

func (i *Interpreter) BooleanInfix(operator string, left, right runtime.Data) (runtime.Data, error) {
	leftVal := left.(*runtime.TBoolean).Value
	rightVal := right.(*runtime.TBoolean).Value
	switch operator {
	case token.LogicalAnd:
		return i.nativeToBoolean(leftVal && rightVal), nil
	case token.LogicalOr:
		return i.nativeToBoolean(leftVal || rightVal), nil
	case token.Equal:
		return i.nativeToBoolean(leftVal == rightVal), nil
	case token.NotEqual:
		return i.nativeToBoolean(leftVal != rightVal), nil
	default:
		return nil, rerror.ErrorFmt("Unsupported Boolean operator '%s'", operator)
	}
}

func (i *Interpreter) ArrayInfix(operator string, left, right runtime.Data) (runtime.Data, error) {
	leftVal := left.(*runtime.TArray).Elements
	rightVal := right.(*runtime.TArray).Elements
	switch operator {
	case string(token.Plus):
		return &runtime.TArray{Elements: append(leftVal, rightVal...)}, nil
	case token.Equal:
		return i.nativeToBoolean(i.compareArrays(leftVal, rightVal)), nil
	case token.NotEqual:
		return i.nativeToBoolean(!i.compareArrays(leftVal, rightVal)), nil
	case string(token.Less):
		return i.nativeToBoolean(len(leftVal) < len(rightVal)), nil
	case string(token.Greater):
		return i.nativeToBoolean(len(leftVal) > len(rightVal)), nil
	default:
		return nil, rerror.ErrorFmt("Unsupported Array operator '%s'", operator)
	}
}

func (i *Interpreter) DictionaryInfix(operator string, left, right runtime.Data) (runtime.Data, error) {
	leftVal := left.(*runtime.TDictionary).Pairs
	rightVal := right.(*runtime.TDictionary).Pairs
	switch operator {
	case string(token.Plus):
		for k, v := range leftVal {
			rightVal[k] = v
		}
		return &runtime.TDictionary{Pairs: rightVal}, nil
	case token.Equal:
		return i.nativeToBoolean(i.compareDictionaries(leftVal, rightVal)), nil
	case token.NotEqual:
		return i.nativeToBoolean(!i.compareDictionaries(leftVal, rightVal)), nil
	case string(token.Less):
		return i.nativeToBoolean(len(leftVal) < len(rightVal)), nil
	case string(token.Greater):
		return i.nativeToBoolean(len(leftVal) > len(rightVal)), nil
	default:
		return nil, rerror.ErrorFmt("Unsupported Dictionary operator '%s'", operator)
	}
}

func (i *Interpreter) NilInfix(operator string, left, right runtime.Data) (runtime.Data, error) {
	switch operator {
	case token.Equal:
		if left.Type() != runtime.TTNil || right.Type() != runtime.TTNil {
			return runtime.No, nil
		}
		return runtime.Yes, nil
	case token.NotEqual:
		if left.Type() != runtime.TTNil || right.Type() != runtime.TTNil {
			return runtime.Yes, nil
		}
		return runtime.No, nil
	default:
		return nil, rerror.ErrorFmt("Unsupported Nil operator '%s'", operator)
	}
}

func (i *Interpreter) RangeIntegerInfix(left, right int64) runtime.Data {
	var result []runtime.Data
	if left < right {
		for idx := left; idx <= right; idx++ {
			result = append(result, &runtime.TInteger{Value: idx})
		}
	} else {
		for idx := left; idx >= right; idx-- {
			result = append(result, &runtime.TInteger{Value: idx})
		}
	}
	return &runtime.TArray{Elements: result}
}

func (i *Interpreter) RangeStringInfix(left, right string) (runtime.Data, error) {
	if len(left) > 1 || len(right) > 1 {
		return nil, rerror.ErrorFmt("Range operator expects 2 single character strings")
	}
	var result []runtime.Data
	alphabet := "0123456789abcdefghijklmnopqrstuvwxyz"
	leftByte := []int32(strings.ToLower(left))[0]
	rightByte := []int32(strings.ToLower(right))[0]
	if leftByte < rightByte {
		for _, v := range alphabet {
			if v >= leftByte && v <= rightByte {
				result = append(result, &runtime.TString{Value: string(v)})
			}
		}
	} else {
		for i := len(alphabet) - 1; i >= 0; i-- {
			v := int32(alphabet[i])
			if v <= leftByte && v >= rightByte {
				result = append(result, &runtime.TString{Value: string(v)})
			}
		}
	}
	return &runtime.TArray{Elements: result}, nil
}

func (i *Interpreter) checkExtFileName(file string) string {
	ext := filepath.Ext(file)
	if ext == "" {
		file = file + util.FileExtension()
	}
	return file
}

func (i *Interpreter) shouldBreakImmediately(data runtime.Data) bool {
	switch data.Type() {
	case runtime.TTBreak, runtime.TTContinue, runtime.TTReturn:
		return true
	default:
		return false
	}
}

func (i *Interpreter) unwrapReturnValue(data runtime.Data) runtime.Data {
	if returnVal, ok := data.(*runtime.TReturn); ok {
		return returnVal.Value
	}
	return data
}

func (i *Interpreter) compareArrays(left, right []runtime.Data) bool {
	if len(left) != len(right) {
		return false
	}
	for i, v := range left {
		if v.Type() != right[i].Type() || v.Check() != right[i].Check() {
			return false
		}
	}
	return true
}

func (i *Interpreter) compareDictionaries(left, right map[runtime.Data]runtime.Data) bool {
	if len(left) != len(right) {
		return false
	}
	found := 0
	for lk, lv := range left {
		for rk, rv := range right {
			if lk.Check() == rk.Check() && lv.Check() == rv.Check() {
				found += 1
				continue
			}
		}
	}
	return found == len(left)
}

func (i *Interpreter) stringToArray(str *runtime.TString) *runtime.TArray {
	array := &runtime.TArray{}
	array.Elements = []runtime.Data{}
	for _, s := range str.Value {
		array.Elements = append(array.Elements, &runtime.TString{Value: string(s)})
	}
	return array
}

func (i *Interpreter) nativeToBoolean(value bool) runtime.Data {
	if value {
		return runtime.Yes
	}
	return runtime.No
}

func (i *Interpreter) isType(data runtime.Data) bool {
	switch data := data.(type) {
	case *runtime.TBoolean:
		return data.Value
	case *runtime.TString:
		return data.Value != ""
	case *runtime.TInteger:
		return data.Value != 0
	case *runtime.TFloat:
		return data.Value != 0.0
	case *runtime.TArray:
		return len(data.Elements) > 0
	case *runtime.TDictionary:
		return len(data.Pairs) > 0
	case *runtime.TSymbol:
		return true
	case *runtime.TNil:
		return false
	default:
		return false
	}
}

func (i *Interpreter) checkArrayBounds(array []runtime.Data, index int64) (int64, error) {
	originalIdx := index
	if index < 0 {
		index = int64(len(array)) + index
	}
	if index < 0 || index > int64(len(array)-1) {
		return 0, rerror.ErrorFmt("Array index '%d' out of bounds", originalIdx)
	}
	return index, nil
}

func (i *Interpreter) checkStringBounds(str string, index int64) (int64, error) {
	originalIdx := index
	if index < 0 {
		index = int64(len(str)) + index
	}
	if index < 0 || index > int64(len(str)-1) {
		return 0, rerror.ErrorFmt("String index '%d' out of bounds", originalIdx)
	}
	return index, nil
}

func (i *Interpreter) checkSupportedType(t string) bool {
	switch t {
	case runtime.TTBoolean, runtime.TTString, runtime.TTInteger, runtime.TTFloat,
		runtime.TTArray, runtime.TTDictionary, runtime.TTSymbol, runtime.TTFunction:
		return true
	default:
		return false
	}
}

func (i *Interpreter) checkTypeMatch(actual, expected string) error {
	if !i.checkSupportedType(actual) {
		return rerror.ErrorFmt("Unknown type '%s' in function parameter", actual)
	}
	if actual != expected {
		return rerror.ErrorFmt("Function asks for type '%s' but got '%s'", expected, actual)
	}
	return nil
}

func (i *Interpreter) interpreterError(n ast.Node, msg string) {
	rerror.Error(rerror.Runtime, n.TokenPosition(), msg)
}
