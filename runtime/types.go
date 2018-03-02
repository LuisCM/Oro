// Copyright 2011 The LuisCM. All rights reserved.
// Use of this source code is license that can be found in the LICENSE file.

// Package runtime implements functions to types.
package runtime

import (
	"bytes"
	"fmt"
	"github.com/luiscm/oro/ast"
	"github.com/luiscm/oro/token"
	"strings"
)

const (
	TTBoolean     = "Boolean"
	TTString      = "String"
	TTInteger     = "Integer"
	TTFloat       = "Float"
	TTArray       = "Array"
	TTDictionary  = "Dictionary"
	TTSymbol      = "Symbol"
	TTPlaceHolder = "PlaceHolder"
	TTFunction    = "Function"
	TTModule      = "Module"
	TTBreak       = "Break"
	TTContinue    = "Continue"
	TTReturn      = "Return"
	TTNil         = "Nil"
)

var (
	Nil = &TNil{}
	Yes = &TBoolean{Value: true}
	No  = &TBoolean{Value: false}
)

type Data interface {
	Type() string
	Check() string
}

type TBoolean struct {
	Value bool
}

func (t *TBoolean) Type() string {
	return TTBoolean
}

func (t *TBoolean) Check() string {
	return fmt.Sprintf("%t", t.Value)
}

type TString struct {
	Value string
}

func (t *TString) Type() string {
	return TTString
}

func (t *TString) Check() string {
	return t.Value
}

type TInteger struct {
	Value int64
}

func (t *TInteger) Type() string {
	return TTInteger
}

func (t *TInteger) Check() string {
	return fmt.Sprintf("%d", t.Value)
}

type TFloat struct {
	Value float64
}

func (t *TFloat) Type() string {
	return TTFloat
}

func (t *TFloat) Check() string {
	return fmt.Sprintf("%f", t.Value)
}

type TArray struct {
	Elements []Data
}

func (t *TArray) Type() string {
	return TTArray
}

func (t *TArray) Check() string {
	var out bytes.Buffer
	var elements []string
	for _, e := range t.Elements {
		elements = append(elements, e.Check())
	}
	out.WriteString(token.LeftBracket)
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString(token.RightBracket)
	return out.String()
}

type TDictionary struct {
	Pairs map[Data]Data
}

func (t *TDictionary) Type() string {
	return TTDictionary
}

func (t *TDictionary) Check() string {
	var out bytes.Buffer
	var pairs []string
	for key, value := range t.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s => %s", key.Check(), value.Check()))
	}
	out.WriteString(token.LeftBracket)
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString(token.RightBracket)
	return out.String()
}

type TSymbol struct {
	Value string
}

func (t *TSymbol) Type() string {
	return TTSymbol
}

func (t *TSymbol) Check() string {
	var out bytes.Buffer
	out.WriteString(token.Colon)
	out.WriteString(t.Value)
	return out.String()
}

type TPlaceHolder struct{}

func (t *TPlaceHolder) Type() string {
	return TTPlaceHolder
}

func (t *TPlaceHolder) Check() string {
	return TTPlaceHolder
}

type TNil struct{}

func (t *TNil) Type() string {
	return TTNil
}

func (t *TNil) Check() string {
	return token.Nil
}

type TFunction struct {
	Parameters []*ast.FunctionParameter
	Body       *ast.BlockStatement
	ReturnType *ast.Identifier
	Variadic   bool
	Scope      *Scope
}

func (t *TFunction) Type() string {
	return TTFunction
}

func (t *TFunction) Check() string {
	var out bytes.Buffer
	var parameters []string
	for i, v := range t.Parameters {
		param := v.Check()
		if t.Variadic && i == len(t.Parameters)-1 {
			param = token.Ellipsis + param
		}
		parameters = append(parameters, param)
	}
	out.WriteString(token.Function)
	out.WriteString(token.Space)
	out.WriteString(token.LeftParenthesis)
	out.WriteString(strings.Join(parameters, ", "))
	out.WriteString(token.RightParenthesis)
	out.WriteString(token.Space)
	out.WriteString(token.Arrow)
	out.WriteString(token.Space)
	out.WriteString(t.Body.Check())
	return out.String()
}

type TModule struct {
	Name *ast.Identifier
	Body *ast.BlockStatement
}

func (t *TModule) Type() string {
	return TTModule
}

func (t *TModule) Check() string {
	var out bytes.Buffer
	out.WriteString(token.Module)
	out.WriteString(token.Space)
	out.WriteString(t.Name.Check())
	out.WriteString(token.Space)
	out.WriteString(token.LeftBraces)
	out.WriteString(token.Space)
	out.WriteString(t.Body.Check())
	out.WriteString(token.Space)
	out.WriteString(token.RightBraces)
	out.WriteString(token.Space)
	return out.String()
}

type TBreak struct{}

func (t *TBreak) Type() string {
	return TTBreak
}

func (t *TBreak) Check() string {
	return TTBreak
}

type TContinue struct{}

func (t *TContinue) Type() string {
	return TTContinue
}

func (t *TContinue) Check() string {
	return TTContinue
}

type TReturn struct {
	Value Data
}

func (t *TReturn) Type() string {
	return TTReturn
}

func (t *TReturn) Check() string {
	return t.Value.Check()
}
