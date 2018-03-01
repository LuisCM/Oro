// Copyright 2011 The LuisCM. All rights reserved.
// Use of this source code is license that can be found in the LICENSE file.

package parser

import (
	"github.com/luiscm/oro/ast"
	"github.com/luiscm/oro/lexer"
	"github.com/luiscm/oro/rerror"
	"testing"
)

func TestString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"first test"`, "first test"},
		{`"second test"`, "second test"},
	}
	for _, test := range tests {
		lex := lexer.New([]byte(test.input))
		parse := New(lex)
		program := parse.Parse()
		checkParserErrors(t)
		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("Expected an ast.ExpressionStatement but got %T", program.Statements[0])
		}
		literal, ok := statement.Expression.(*ast.String)
		if !ok {
			t.Errorf("Expected an ast.String but got %T", statement.Expression)
		}
		if literal.Value != test.expected {
			t.Errorf("Expected a String %s but got %s", test.expected, literal.Value)
		}
	}
}

func TestInteger(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{`10`, 10},
		{`12345`, 12345},
	}
	for _, test := range tests {
		lex := lexer.New([]byte(test.input))
		parse := New(lex)
		program := parse.Parse()
		checkParserErrors(t)
		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("Expected an ast.ExpressionStatement but got %T", program.Statements[0])
		}
		literal, ok := statement.Expression.(*ast.Integer)
		if !ok {
			t.Errorf("Expected an ast.Integer but got %T", statement.Expression)
		}
		if literal.Value != test.expected {
			t.Errorf("Expected an Integer %d but got %d", test.expected, literal.Value)
		}
	}
}

func TestFloat(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{`10.2`, 10.2},
		{`1050.23488`, 1050.23488},
	}
	for _, test := range tests {
		lex := lexer.New([]byte(test.input))
		parse := New(lex)
		program := parse.Parse()
		checkParserErrors(t)
		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("Expected an ast.ExpressionStatement but got %T", program.Statements[0])
		}
		literal, ok := statement.Expression.(*ast.Float)
		if !ok {
			t.Errorf("Expected an ast.Float but got %T", statement.Expression)
		}
		if literal.Value != test.expected {
			t.Errorf("Expected a Float %f but got %f", test.expected, literal.Value)
		}
	}
}

func TestBoolean(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`true`, true},
		{`false`, false},
	}
	for _, test := range tests {
		lex := lexer.New([]byte(test.input))
		parse := New(lex)
		program := parse.Parse()
		checkParserErrors(t)
		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("Expected an ast.ExpressionStatement but got %T", program.Statements[0])
		}
		literal, ok := statement.Expression.(*ast.Boolean)
		if !ok {
			t.Errorf("Expected an ast.Boolean but got %T", statement.Expression)
		}
		if literal.Value != test.expected {
			t.Errorf("Expected a Boolean %t but got %t", test.expected, literal.Value)
		}
	}
}

func TestArray(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{`[1, 2, 3]`, 3},
		{`["a", "b"]`, 2},
		{`[]`, 0},
	}
	for _, test := range tests {
		lex := lexer.New([]byte(test.input))
		parse := New(lex)
		program := parse.Parse()
		checkParserErrors(t)
		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("Expected an ast.ExpressionStatement but got %T", program.Statements[0])
		}
		literal, ok := statement.Expression.(*ast.Array)
		if !ok {
			t.Errorf("Expected an ast.Array but got %T", statement.Expression)
		}
		if len(literal.List.Elements) != test.expected {
			t.Errorf("Expected an Array with %d elements but got %d", test.expected, len(literal.List.Elements))
		}
	}
}

func TestDictionary(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{`["a" => "b", "c" => 2]`, 2},
		{`["a" => "b", "c" => 2, "d" => 10]`, 3},
	}
	for _, test := range tests {
		lex := lexer.New([]byte(test.input))
		parse := New(lex)
		program := parse.Parse()
		checkParserErrors(t)
		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("Expected an ast.ExpressionStatement but got %T", program.Statements[0])
		}
		literal, ok := statement.Expression.(*ast.Dictionary)
		if !ok {
			t.Errorf("Expected an ast.Dictionary but got %T", statement.Expression)
		}
		if len(literal.Pairs) != test.expected {
			t.Errorf("Expected a Dictionary with %d elements but got %d", test.expected, len(literal.Pairs))
		}
	}
}

func TestVar(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{`var one = 1`, 1},
		{`var two = 2`, 2},
	}
	for _, test := range tests {
		lex := lexer.New([]byte(test.input))
		parse := New(lex)
		program := parse.Parse()
		checkParserErrors(t)
		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("Expected an ast.ExpressionStatement but got %T", program.Statements[0])
		}
		literal, ok := statement.Expression.(*ast.Var)
		if !ok {
			t.Errorf("Expected an ast.Var but got %T", statement.Expression)
		}
		if literal.Value == nil {
			t.Errorf("Expected a condition but got nil")
		}
	}
}

func TestVal(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`val first = "string"`, "first"},
		{`val second = 10 + 2`, "second"},
	}
	for _, test := range tests {
		lex := lexer.New([]byte(test.input))
		parse := New(lex)
		program := parse.Parse()
		checkParserErrors(t)
		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("Expected an ast.ExpressionStatement but got %T", program.Statements[0])
		}
		literal, ok := statement.Expression.(*ast.Val)
		if !ok {
			t.Errorf("Expected an ast.Val but got %T", statement.Expression)
		}
		if literal.Value == nil {
			t.Errorf("Expected a condition but got nil")
		}
	}
}

func TestIf(t *testing.T) {
	input := `if a == 1
  a + 2
else
  a + 3
end`
	lex := lexer.New([]byte(input))
	parse := New(lex)
	program := parse.Parse()
	checkParserErrors(t)
	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("Expected an ast.ExpressionStatement but got %T", program.Statements[0])
	}
	literal, ok := statement.Expression.(*ast.If)
	if !ok {
		t.Errorf("Expected an ast.If but got %T", statement.Expression)
	}
	if literal.Condition == nil {
		t.Errorf("Expected a condition but got nil")
	}
	if len(literal.Then.Statements) != 1 {
		t.Errorf("Expected %d statement in THEN block but got %d", 1, len(literal.Then.Statements))
	}
	if len(literal.Else.Statements) != 1 {
		t.Errorf("Expected %d statement in ELSE block but got %d", 1, len(literal.Then.Statements))
	}
}

func TestRepeat(t *testing.T) {
	input := `repeat a, b in arr
  a + 1
end`
	lex := lexer.New([]byte(input))
	parse := New(lex)
	program := parse.Parse()
	checkParserErrors(t)
	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("Expected an ast.ExpressionStatement but got %T", program.Statements[0])
	}
	literal, ok := statement.Expression.(*ast.Repeat)
	if !ok {
		t.Errorf("Expected an ast.Repeat but got %T", statement.Expression)
	}
	if len(literal.Arguments.Elements) != 2 {
		t.Errorf("Expected %d arguments but got %d", 2, len(literal.Arguments.Elements))
	}
	if literal.Enumerable == nil {
		t.Errorf("Expected an enumerable but got nil")
	}
	if len(literal.Body.Statements) != 1 {
		t.Errorf("Expected %d statement in body but got %d", 1, len(literal.Body.Statements))
	}
}

func TestMatch(t *testing.T) {
	input := `match a
when 1
  a + 1
when 2
  a + 2
else
  a + 3
end`
	lex := lexer.New([]byte(input))
	parse := New(lex)
	program := parse.Parse()
	checkParserErrors(t)
	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("Expected an ast.ExpressionStatement but got %T", program.Statements[0])
	}
	literal, ok := statement.Expression.(*ast.Match)
	if !ok {
		t.Errorf("Expected an ast.Match but got %T", statement.Expression)
	}
	if literal.Control == nil {
		t.Errorf("Expected a control expression but got nil")
	}
	if len(literal.Whens) != 2 {
		t.Errorf("Expected %d whens but got %d",2, len(literal.Whens))
	}
	for _, ws := range literal.Whens {
		if len(ws.Body.Statements) != 1 {
			t.Errorf("Expected %d statement in when but got %d",1, len(ws.Body.Statements))
		}
	}
	if len(literal.Else.Statements) != 1 {
		t.Errorf("Expected %d statement in else when but got %d",1, len(literal.Else.Statements))
	}
}

func TestModule(t *testing.T) {
	input := `module math
  val pi = 3.14
  val add = fn (x, y)
    x + y
  end
end`
	lex := lexer.New([]byte(input))
	parse := New(lex)
	program := parse.Parse()
	checkParserErrors(t)
	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("Expected an ast.ExpressionStatement but got %T", program.Statements[0])
	}
	literal, ok := statement.Expression.(*ast.Module)
	if !ok {
		t.Errorf("Expected an ast.Module but got %T", statement.Expression)
	}
	if literal.Name.Value != "math" {
		t.Errorf("Expected module name %s but got %s", "math", literal.Name.Value)
	}
	if len(literal.Body.Statements) != 2 {
		t.Errorf("Expected %d statements in body but got %d", 2, len(literal.Body.Statements))
	}
}

func TestFunction(t *testing.T) {
	input := `fn x, y, z
  x + y
  z + x
end`
	lex := lexer.New([]byte(input))
	parse := New(lex)
	program := parse.Parse()
	checkParserErrors(t)
	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("Expected an ast.ExpressionStatement but got %T", program.Statements[0])
	}
	literal, ok := statement.Expression.(*ast.Function)
	if !ok {
		t.Errorf("Expected an ast.Function but got %T", statement.Expression)
	}
	if len(literal.Parameters) != 3 {
		t.Errorf("Expected %d parameters but got %d", 3, len(literal.Parameters))
	}
	if len(literal.Body.Statements) != 2 {
		t.Errorf("Expected %d statements in body but got %d", 2, len(literal.Body.Statements))
	}
}

func TestModuleAccess(t *testing.T) {
	input := `Math.pi`
	lex := lexer.New([]byte(input))
	parse := New(lex)
	program := parse.Parse()
	checkParserErrors(t)
	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("Expected an ast.ExpressionStatement but got %T", program.Statements[0])
	}
	literal, ok := statement.Expression.(*ast.ModuleAccess)
	if !ok {
		t.Errorf("Expected an ast.ModuleAccess but got %T", statement.Expression)
	}
	if literal.Object.Value != "Math" {
		t.Errorf("Expected %s as object but got %s", "Math", literal.Object.Value)
	}
	if literal.Parameter.Value != "pi" {
		t.Errorf("Expected %s as parameter but got %s", "Math", literal.Parameter.Value)
	}
}

func TestFunctionCall(t *testing.T) {
	input := `myfn (1, 2)`
	lex := lexer.New([]byte(input))
	parse := New(lex)
	program := parse.Parse()
	checkParserErrors(t)
	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("Expected an ast.ExpressionStatement but got %T", program.Statements[0])
	}
	literal, ok := statement.Expression.(*ast.FunctionCall)
	if !ok {
		t.Errorf("Expected an ast.FunctionCall but got %T", statement.Expression)
	}
	function, ok := literal.Function.(*ast.Identifier)
	if !ok {
		t.Errorf("Expected an ast.Identifier as function name but got %T", statement.Expression)
	}
	if function.Value != "myfn" {
		t.Errorf("Expected %s as function name but got %s", "myfn", function.Value)
	}
	if len(literal.Arguments.Elements) != 2 {
		t.Errorf("Expected %d arguments but got %d", 2, len(literal.Arguments.Elements))
	}
}

func TestSubscript(t *testing.T) {
	input := `arr[1]`
	lex := lexer.New([]byte(input))
	parse := New(lex)
	program := parse.Parse()
	checkParserErrors(t)
	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("Expected an ast.ExpressionStatement but got %T", program.Statements[0])
	}
	literal, ok := statement.Expression.(*ast.Subscript)
	if !ok {
		t.Errorf("Expected an ast.Subscript but got %T", statement.Expression)
	}
	left, ok := literal.Left.(*ast.Identifier)
	if !ok {
		t.Errorf("Expected an ast.Identifier as subscript caller but got %T", statement.Expression)
	}
	index, ok := literal.Index.(*ast.Integer)
	if !ok {
		t.Errorf("Expected an ast.Integer as subscript but got %T", statement.Expression)
	}
	if left.Value != "arr" {
		t.Errorf("Expected %s as subscript caller but got %s", "arr", left.Value)
	}
	if index.Value != 1 {
		t.Errorf("Expected %d as subscript but got %d", 1, index.Value)
	}
}

func TestOperatorPrecedence(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a - b * c",
			"(a - (b * c))",
		},
		{
			"a / b * c + d",
			"(((a / b) * c) + d)",
		},
		{
			"(a + b) * c",
			"((a + b) * c)",
		},
		{
			"a % b * c",
			"((a % b) * c)",
		},
		{
			"a * b ** c",
			"(a * (b ** c))",
		},
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"b > a && c <= d || e > f",
			"((b > a) && ((c <= d) || (e > f)))",
		},
		{
			"c > b || d <= e && f > g",
			"((c > b) || ((d <= e) && (f > g)))",
		},
		{
			"a + b == c + d",
			"((a + b) == (c + d))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"!true != true",
			"((!true) != true)",
		},
		{
			"a >> b + c",
			"(a >> (b + c))",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"(5 + 5) * 2 * (5 + 5)",
			"(((5 + 5) * 2) * (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
	}
	for _, test := range tests {
		lex := lexer.New([]byte(test.input))
		parse := New(lex)
		program := parse.Parse()
		checkParserErrors(t)
		actual := program.Check()
		if actual != test.expected {
			t.Errorf("Expected %s but got %s", test.expected, actual)
		}
	}
}

func checkParserErrors(t *testing.T) {
	if rerror.HasErrors() {
		t.Errorf(rerror.ParseErrors)
		for _, e := range rerror.GetErrors() {
			t.Errorf(e)
		}
	}
}
