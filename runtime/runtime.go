// Copyright 2011 The LuisCM. All rights reserved.
// Use of this source code is license that can be found in the LICENSE file.

// Package runtime implements functions to runtime.
package runtime

import (
	"bufio"
	"fmt"
	"github.com/luiscm/oro/rerror"
	"github.com/luiscm/oro/util"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type TRuntimeFn func(args ...Data) (Data, error)

var FnRuntime = map[string]TRuntimeFn{

	"echo": func(args ...Data) (Data, error) {
		for _, arg := range args {
			fmt.Println(arg.Check())
		}
		return &TString{Value: ""}, nil
	},

	"put": func(args ...Data) (Data, error) {
		for _, arg := range args {
			fmt.Print(arg.Check())
		}
		return &TString{Value: ""}, nil
	},

	"puts": func(args ...Data) (Data, error) {
		for _, arg := range args {
			fmt.Println(arg.Check())
		}
		return &TString{Value: ""}, nil
	},

	"write": func(args ...Data) (Data, error) {
		for _, arg := range args {
			fmt.Print(arg.Check())
		}
		return &TString{Value: ""}, nil
	},

	"writeln": func(args ...Data) (Data, error) {
		for _, arg := range args {
			fmt.Println(arg.Check())
		}
		return &TString{Value: ""}, nil
	},

	"print": func(args ...Data) (Data, error) {
		for _, arg := range args {
			fmt.Print(arg.Check())
		}
		return &TString{Value: ""}, nil
	},

	"println": func(args ...Data) (Data, error) {
		for _, arg := range args {
			fmt.Println(arg.Check())
		}
		return &TString{Value: ""}, nil
	},

	"prompt": func(args ...Data) (Data, error) {
		reader := bufio.NewReader(os.Stdin)
		if len(args) > 0 {
			fmt.Print(strings.Trim(args[0].Check(), "\""))
		}
		out, _ := reader.ReadString('\n')
		return &TString{Value: strings.Trim(out, "\r\n")}, nil
	},

	"quit": func(args ...Data) (Data, error) {
		os.Exit(0)
		return nil, nil
	},

	"panic": func(args ...Data) (Data, error) {
		var message string
		if len(args) > 0 {
			message = args[0].Check()
		}
		return nil, rerror.ErrorFmt(message)
	},

	"typeof": func(args ...Data) (Data, error) {
		if len(args) != 1 {
			return nil, rerror.ErrorFmt("typeof() expects exactly 1 argument")
		}
		return &TString{Value: args[0].Type()}, nil
	},

	"len": func(args ...Data) (Data, error) {
		if len(args) != 1 {
			return nil, rerror.ErrorFmt("len() expects exactly 1 argument")
		}
		switch object := args[0].(type) {
		case *TArray:
			return &TInteger{Value: int64(len(object.Elements))}, nil
		case *TString:
			return &TInteger{Value: int64(len(object.Value))}, nil
		default:
			return nil, rerror.ErrorFmt("argument to `len` not supported, got %s", object.Type())
		}
	},

	"first": func(args ...Data) (Data, error) {
		if len(args) != 1 {
			return nil, rerror.ErrorFmt("first() expects exactly 1 argument")
		}
		if args[0].Type() != TTArray {
			return nil, rerror.ErrorFmt("argument to `first` must be ARRAY, got %s", args[0].Type())
		}
		arr := args[0].(*TArray)
		if len(arr.Elements) > 0 {
			return arr.Elements[0], nil
		}
		return nil, nil
	},

	"last": func(args ...Data) (Data, error) {
		if len(args) != 1 {
			return nil, rerror.ErrorFmt("last() expects exactly 1 argument")
		}
		if args[0].Type() != TTArray {
			return nil, rerror.ErrorFmt("argument to `last` must be ARRAY, got %s", args[0].Type())
		}
		arr := args[0].(*TArray)
		length := len(arr.Elements)
		if length > 0 {
			return arr.Elements[length-1], nil
		}
		return nil, nil
	},

	"rest": func(args ...Data) (Data, error) {
		if len(args) != 1 {
			return nil, rerror.ErrorFmt("rest() expects exactly 1 argument")
		}
		if args[0].Type() != TTArray {
			return nil, rerror.ErrorFmt("argument to `rest` must be ARRAY, got %s", args[0].Type())
		}
		arr := args[0].(*TArray)
		length := len(arr.Elements)
		if length > 0 {
			newElements := make([]Data, length-1, length-1)
			copy(newElements, arr.Elements[1:length])
			return &TArray{Elements: newElements}, nil
		}
		return nil, nil
	},

	"push": func(args ...Data) (Data, error) {
		if len(args) != 2 {
			return nil, rerror.ErrorFmt("push() expects exactly 2 argument")
		}
		if args[0].Type() != TTArray {
			return nil, rerror.ErrorFmt("argument to `push` must be ARRAY, got %s", args[0].Type())
		}
		arr := args[0].(*TArray)
		length := len(arr.Elements)
		newElements := make([]Data, length+1, length+1)
		copy(newElements, arr.Elements)
		newElements[length] = args[1]
		return &TArray{Elements: newElements}, nil
	},

	"String": func(args ...Data) (Data, error) {
		if len(args) != 1 {
			return nil, rerror.ErrorFmt("String() expects exactly 1 argument")
		}
		switch object := args[0].(type) {
		case *TInteger:
			return &TString{Value: fmt.Sprintf("%d", object.Value)}, nil
		case *TFloat:
			return &TString{Value: fmt.Sprintf("%f", object.Value)}, nil
		case *TBoolean:
			return &TString{Value: fmt.Sprintf("%t", object.Value)}, nil
		case *TString:
			return object, nil
		default:
			return nil, rerror.ErrorFmt("String() can't convert '%s' to String", object.Type())
		}
	},

	"Int": func(args ...Data) (Data, error) {
		if len(args) != 1 {
			return nil, rerror.ErrorFmt("Int() expects exactly 1 argument")
		}
		switch object := args[0].(type) {
		case *TString:
			i, err := strconv.Atoi(object.Value)
			if err != nil {
				return nil, rerror.ErrorFmt("Int() can't convert '%s' to Integer", object.Value)
			}
			return &TInteger{Value: int64(i)}, nil
		case *TFloat:
			return &TInteger{Value: int64(object.Value)}, nil
		case *TBoolean:
			result := 0
			if object.Value {
				result = 1
			}
			return &TInteger{Value: int64(result)}, nil
		case *TInteger:
			return object, nil
		default:
			return nil, rerror.ErrorFmt("Int() can't convert '%s' to Integer", object.Type())
		}
	},

	"Float": func(args ...Data) (Data, error) {
		if len(args) != 1 {
			return nil, rerror.ErrorFmt("Float() expects exactly 1 argument")
		}
		switch object := args[0].(type) {
		case *TString:
			i, err := strconv.ParseFloat(object.Value, 64)
			if err != nil {
				return nil, rerror.ErrorFmt("Float() can't convert '%s' to Integer", object.Value)
			}
			return &TFloat{Value: i}, nil
		case *TInteger:
			return &TFloat{Value: float64(object.Value)}, nil
		case *TBoolean:
			result := 0
			if object.Value {
				result = 1
			}
			return &TFloat{Value: float64(result)}, nil
		case *TFloat:
			return &TFloat{Value: object.Value}, nil
		default:
			return nil, rerror.ErrorFmt("Float() can't convert '%s' to Integer", object.Type())
		}
	},

	"Array": func(args ...Data) (Data, error) {
		if len(args) != 1 {
			return nil, rerror.ErrorFmt("Array() expects exactly 1 argument")
		}
		switch object := args[0].(type) {
		case *TArray:
			return object, nil
		default:
			return &TArray{Elements: []Data{object}}, nil
		}
	},

	"runtime_rand": func(args ...Data) (Data, error) {
		if len(args) != 2 {
			return nil, rerror.ErrorFmt("runtime_rand() expects exactly 2 arguments")
		}
		if args[0].Type() != TTInteger || args[1].Type() != TTInteger {
			return nil, rerror.ErrorFmt("runtime_rand() expects min and max as Integers")
		}
		min := int(args[0].(*TInteger).Value)
		max := int(args[1].(*TInteger).Value)
		if max < min {
			return nil, rerror.ErrorFmt("runtime_rand() expects max higher than min")
		}
		rand.Seed(time.Now().UnixNano())
		random := rand.Intn(max-min) + min
		return &TInteger{Value: int64(random)}, nil
	},

	"runtime_tolower": func(args ...Data) (Data, error) {
		if len(args) != 1 {
			return nil, rerror.ErrorFmt("runtime_tolower() expects exactly 1 argument")
		}
		if args[0].Type() != TTString {
			return nil, rerror.ErrorFmt("runtime_tolower() expects a String")
		}
		str := args[0].(*TString).Value
		return &TString{Value: strings.ToLower(str)}, nil
	},

	"runtime_toupper": func(args ...Data) (Data, error) {
		if len(args) != 1 {
			return nil, rerror.ErrorFmt("runtime_toupper() expects exactly 1 argument")
		}
		if args[0].Type() != TTString {
			return nil, rerror.ErrorFmt("runtime_toupper() expects a String")
		}
		str := args[0].(*TString).Value
		return &TString{Value: strings.ToUpper(str)}, nil
	},

	"runtime_regex_match": func(args ...Data) (Data, error) {
		if len(args) != 2 {
			return nil, rerror.ErrorFmt("runtime_regex_match() expects exactly 2 arguments")
		}
		if args[0].Type() != TTString {
			return nil, rerror.ErrorFmt("runtime_regex_match() expects a String")
		}
		if args[1].Type() != TTString {
			return nil, rerror.ErrorFmt("runtime_regex_match() expects a String regex")
		}
		object := args[0].(*TString).Value
		match := args[1].(*TString).Value
		regx, err := regexp.Compile(match)
		if err != nil {
			return nil, rerror.ErrorFmt("runtime_regex_match() couldn't compile the regular expression")
		}
		return &TBoolean{Value: regx.Find([]byte(object)) != nil}, nil
	},

	"Environment": func(args ...Data) (Data, error) {
		return &TString{Value: util.Environment()}, nil
	},

	"NameVersionEnvironment": func(args ...Data) (Data, error) {
		return &TString{Value: util.NameVersionEnvironment()}, nil
	},

	"Name": func(args ...Data) (Data, error) {
		return &TString{Value: util.Name()}, nil
	},

	"Version": func(args ...Data) (Data, error) {
		return &TString{Value: util.Version()}, nil
	},

	"NameVersion": func(args ...Data) (Data, error) {
		return &TString{Value: util.NameVersion()}, nil
	},

	"AuthorName": func(args ...Data) (Data, error) {
		return &TString{Value: util.AuthorName()}, nil
	},

	"AuthorEmail": func(args ...Data) (Data, error) {
		return &TString{Value: util.AuthorEmail()}, nil
	},

	"Copyright": func(args ...Data) (Data, error) {
		return &TString{Value: util.CopyrightDescription()}, nil
	},
}
