// Copyright 2011 The LuisCM. All rights reserved.
// Use of this source code is license that can be found in the LICENSE file.

package interpreter

import (
	"github.com/luiscm/oro/lexer"
	"github.com/luiscm/oro/parser"
	"github.com/luiscm/oro/rerror"
	"github.com/luiscm/oro/runtime"
	"testing"
)

func TestInterpreterString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"hello"`, "hello"},
		{`"hello"+"world"`, "helloworld"},
		{`"hello"+" "+"world"`, "hello world"},
	}
	for _, test := range tests {
		lex := lexer.New([]byte(test.input))
		parse := parser.New(lex)
		program := parse.Parse()
		runner := New()
		actual := runner.Interpreter(program, runtime.NewScope())
		checkInterpreterErrors(t)
		testString(t, actual, test.expected)
	}
}

func TestInterpreterInteger(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{`10`, 10},
		{`1234567`, 1234567},
		{`1 + 1`, 2},
		{`-10`, -10},
		{`-10 + 10`, 0},
		{`5 * 2`, 10},
		{`5 * (2 + 2)`, 20},
		{`2 ** 8`, 256},
		{`5 % 2`, 1},
	}
	for _, test := range tests {
		lex := lexer.New([]byte(test.input))
		parse := parser.New(lex)
		program := parse.Parse()
		runner := New()
		actual := runner.Interpreter(program, runtime.NewScope())
		checkInterpreterErrors(t)
		testInteger(t, actual, test.expected)
	}
}

func TestInterpreterFloat(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{`10.0`, 10.0},
		{`10.0 + 1.2`, 11.2},
		{`1 - 0.5`, 0.5},
		{`4.5 * 2`, 9.0},
		{`-5.2`, -5.2},
		{`9.0 / 3`, 3.0},
	}
	for _, test := range tests {
		lex := lexer.New([]byte(test.input))
		parse := parser.New(lex)
		program := parse.Parse()
		runner := New()
		actual := runner.Interpreter(program, runtime.NewScope())
		checkInterpreterErrors(t)
		testFloat(t, actual, test.expected)
	}
}

func TestInterpreterBoolean(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`true`, true},
		{`false`, false},
		{`!false`, true},
		{`1 == 1`, true},
		{`1 == 2`, false},
		{`1 != 2`, true},
		{`1 != 1`, false},
		{`5 > 1`, true},
		{`5 >= 5`, true},
		{`10 > 100`, false},
		{`(1 < 2) == (2 > 1)`, true},
		{`5.3 > 5.2`, true},
		{`"four" > "one"`, true},
		{`"hello" == "world"`, false},
		{`[1, 2] == [3, 4]`, false},
		{`[1, 2] == [1, 2]`, true},
		{`true == !false`, true},
		{`true && true`, true},
		{`true && false`, false},
		{`false || false`, false},
		{`false || true`, true},
	}
	for _, test := range tests {
		lex := lexer.New([]byte(test.input))
		parse := parser.New(lex)
		program := parse.Parse()
		runner := New()
		actual := runner.Interpreter(program, runtime.NewScope())
		checkInterpreterErrors(t)
		testBoolean(t, actual, test.expected)
	}
}

func TestInterpreterVal(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`val x = 5`, 5},
		{`val x = 2`, 2},
		//{`val x = 2 x = 3`, 2},
	}
	for _, test := range tests {
		lex := lexer.New([]byte(test.input))
		parse := parser.New(lex)
		program := parse.Parse()
		runner := New()
		actual := runner.Interpreter(program, runtime.NewScope())
		checkInterpreterErrors(t)
		result, ok := test.expected.(int)
		if !ok {
			t.Errorf("Expected Integer but got %T", test.expected)
		}
		testInteger(t, actual, int64(result))
	}
}

func TestInterpreterIf(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`if 5 > 2 then 10 end`, 10},
		{`if 5 < 2 then 10 else 15 end`, 15},
		{`if true then 10 end`, 10},
	}
	for _, test := range tests {
		lex := lexer.New([]byte(test.input))
		parse := parser.New(lex)
		program := parse.Parse()
		runner := New()
		actual := runner.Interpreter(program, runtime.NewScope())
		checkInterpreterErrors(t)
		result, ok := test.expected.(int)
		if !ok {
			t.Errorf("Expected Integer but got %T", test.expected)
		}
		testInteger(t, actual, int64(result))
	}
}

func TestInterpreterMatch(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`match 1 with when 1 then 10 when 2 then 20 end`, 10},
		{`match 2 with when 1 then 10 when 2 then 20 end`, 20},
		{`match 3 with when 1 then 10 else then 20 end`, 20},
		{`var a = 1 match a with when 1 then a + 1 when 2 then a + 2 else then a + 3 end`, 2},
		{`match true with when true then 100 end`, 100},
		{`val a = 5 match a with when 2, 3 then 2 + 3 when 5 then 5 else then 0 end`, 5},
		{`match ["game", "of", "thrones"] with when "game", "thrones" then 1 when "game", "of", "thrones" then 2 end`, 2},
		{`match ["Luis", "Carlos", 2] with when "Luis", _, _ then 10 when _, _ 2 then 20 else then -1 end`, 10},
	}
	for _, test := range tests {
		lex := lexer.New([]byte(test.input))
		parse := parser.New(lex)
		program := parse.Parse()
		runner := New()
		actual := runner.Interpreter(program, runtime.NewScope())
		checkInterpreterErrors(t)
		result, ok := test.expected.(int)
		if !ok {
			t.Errorf("Expected Integer but got %T", test.expected)
		}
		testInteger(t, actual, int64(result))
	}
}

func testString(t *testing.T, tp runtime.Data, expected string) bool {
	result, ok := tp.(*runtime.TString)
	if !ok {
		t.Errorf("Expected String but got %t", tp)
		return false
	}
	if result.Value != expected {
		t.Errorf("Expected %s but got %s", expected, result.Value)
		return false
	}
	return true
}

func testInteger(t *testing.T, tp runtime.Data, expected int64) bool {
	result, ok := tp.(*runtime.TInteger)
	if !ok {
		t.Errorf("Expected Integer but got %t", tp)
		return false
	}
	if result.Value != expected {
		t.Errorf("Expected %d but got %d", expected, result.Value)
		return false
	}
	return true
}

func testFloat(t *testing.T, tp runtime.Data, expected float64) bool {
	result, ok := tp.(*runtime.TFloat)
	if !ok {
		t.Errorf("Expected Float but got %t", tp)
		return false
	}
	if result.Value != expected {
		t.Errorf("Expected %f but got %f", expected, result.Value)
		return false
	}
	return true
}

func testBoolean(t *testing.T, tp runtime.Data, expected bool) bool {
	result, ok := tp.(*runtime.TBoolean)
	if !ok {
		t.Errorf("Expected Boolean but got %t", tp)
		return false
	}
	if result.Value != expected {
		t.Errorf("Expected %t but got %t", expected, result.Value)
		return false
	}
	return true
}

func checkInterpreterErrors(t *testing.T) {
	if rerror.HasErrors() {
		t.Errorf(rerror.ParseErrors)
		for _, e := range rerror.GetErrors() {
			t.Errorf(e)
		}
	}
}
