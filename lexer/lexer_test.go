// Copyright 2011 The LuisCM. All rights reserved.
// Use of this source code is license that can be found in the LICENSE file.

package lexer

import (
	"github.com/luiscm/oro/token"
	"testing"
)

func TestOperators(t *testing.T) {
	input := `val a = 1 + 2 * 3 % 1 / (5 + 2) ** 2 + 1..5
val b = true && false || 0 >= 1 < 5 && !true
val c = 10 & 5 >> 1 | 0 & ~1`
	tests := []struct {
		Type    token.TType
		Literal string
	}{
		{token.Val, "val"},
		{token.Identifier, "a"},
		{token.Assign, "="},
		{token.Integer, "1"},
		{token.Plus, "+"},
		{token.Integer, "2"},
		{token.Multiply, "*"},
		{token.Integer, "3"},
		{token.Modulus, "%"},
		{token.Integer, "1"},
		{token.Divide, "/"},
		{token.LeftParenthesis, "("},
		{token.Integer, "5"},
		{token.Plus, "+"},
		{token.Integer, "2"},
		{token.RightParenthesis, ")"},
		{token.Exponential, "**"},
		{token.Integer, "2"},
		{token.Plus, "+"},
		{token.Integer, "1"},
		{token.Range, ".."},
		{token.Integer, "5"},
		{token.NewLine, "\n"},
		{token.Val, "val"},
		{token.Identifier, "b"},
		{token.Assign, "="},
		{token.Boolean, "true"},
		{token.LogicalAnd, "&&"},
		{token.Boolean, "false"},
		{token.LogicalOr, "||"},
		{token.Integer, "0"},
		{token.GreaterEqual, ">="},
		{token.Integer, "1"},
		{token.Less, "<"},
		{token.Integer, "5"},
		{token.LogicalAnd, "&&"},
		{token.Not, "!"},
		{token.Boolean, "true"},
		{token.NewLine, "\n"},
		{token.Val, "val"},
		{token.Identifier, "c"},
		{token.Assign, "="},
		{token.Integer, "10"},
		{token.BitwiseAnd, "&"},
		{token.Integer, "5"},
		{token.BitShiftRight, ">>"},
		{token.Integer, "1"},
		{token.BitwiseOr, "|"},
		{token.Integer, "0"},
		{token.BitwiseAnd, "&"},
		{token.BitwiseNot, "~"},
		{token.Integer, "1"},
	}
	lex := New([]byte(input))
	for i, v := range tests {
		tok := lex.NextToken()
		if tok.Type != v.Type || tok.Literal != v.Literal {
			t.Errorf("Expected [%s %s] but got [%s %s] in line %d", string(v.Type), v.Literal, string(tok.Type), tok.Literal, i)
		}
	}
}

func TestDataTypes(t *testing.T) {
	input := `1 5 true 5.20 3.4789 false "yes"`
	tests := []struct {
		Type    token.TType
		Literal string
	}{
		{token.Integer, "1"},
		{token.Integer, "5"},
		{token.Boolean, "true"},
		{token.Float, "5.20"},
		{token.Float, "3.4789"},
		{token.Boolean, "false"},
		{token.String, "yes"},
	}
	lex := New([]byte(input))
	for i, v := range tests {
		tok := lex.NextToken()
		if tok.Type != v.Type || tok.Literal != v.Literal {
			t.Errorf("Expected [%s %s] but got [%s %s] in line %d", string(v.Type), v.Literal, string(tok.Type), tok.Literal, i)
		}
	}
}

func TestDelimiters(t *testing.T) {
	input := `(1, 2, a) ["yes", 5.1, b] [a: b, c: d] a.b a..b`
	tests := []struct {
		Type    token.TType
		Literal string
	}{
		{token.LeftParenthesis, "("},
		{token.Integer, "1"},
		{token.Comma, ","},
		{token.Integer, "2"},
		{token.Comma, ","},
		{token.Identifier, "a"},
		{token.RightParenthesis, ")"},
		{token.LeftBracket, "["},
		{token.String, "yes"},
		{token.Comma, ","},
		{token.Float, "5.1"},
		{token.Comma, ","},
		{token.Identifier, "b"},
		{token.RightBracket, "]"},
		{token.LeftBracket, "["},
		{token.Identifier, "a"},
		{token.Colon, ":"},
		{token.Identifier, "b"},
		{token.Comma, ","},
		{token.Identifier, "c"},
		{token.Colon, ":"},
		{token.Identifier, "d"},
		{token.RightBracket, "]"},
		{token.Identifier, "a"},
		{token.Dot, "."},
		{token.Identifier, "b"},
		{token.Identifier, "a"},
		{token.Range, ".."},
		{token.Identifier, "b"},
	}
	lex := New([]byte(input))
	for i, v := range tests {
		tok := lex.NextToken()
		if tok.Type != v.Type || tok.Literal != v.Literal {
			t.Errorf("Expected [%s %s] but got [%s %s] in line %d", string(v.Type), v.Literal, string(tok.Type), tok.Literal, i)
		}
	}
}

func TestKeywords(t *testing.T) {
	input := `val var fn do end not if else right repeat in left then return middle match not when module yes`
	tests := []struct {
		Type    token.TType
		Literal string
	}{
		{token.Val, "val"},
		{token.Var, "var"},
		{token.Function, "fn"},
		{token.Do, "do"},
		{token.End, "end"},
		{token.Identifier, "not"},
		{token.If, "if"},
		{token.Else, "else"},
		{token.Identifier, "right"},
		{token.Repeat, "repeat"},
		{token.In, "in"},
		{token.Identifier, "left"},
		{token.Then, "then"},
		{token.Return, "return"},
		{token.Identifier, "middle"},
		{token.Match, "match"},
		{token.Identifier, "not"},
		{token.When, "when"},
		{token.Module, "module"},
		{token.Identifier, "yes"},
	}
	lex := New([]byte(input))
	for i, v := range tests {
		tok := lex.NextToken()
		if tok.Type != v.Type || tok.Literal != v.Literal {
			t.Errorf("Expected [%s %s] but got [%s %s] in line %d", string(v.Type), v.Literal, string(tok.Type), tok.Literal, i)
		}
	}
}

func TestMiniProgram(t *testing.T) {
	input := `val a = 10
val b = 20.2
if b > a then
  repeat i in 5..10
    i + 2
  end
else
  "exiting..."
end
val c = fn x, y, z
  "hi" + x + y + z
end`
	tests := []struct {
		Type    token.TType
		Literal string
	}{
		{token.Val, "val"},
		{token.Identifier, "a"},
		{token.Assign, "="},
		{token.Integer, "10"},
		{token.NewLine, "\n"},
		{token.Val, "val"},
		{token.Identifier, "b"},
		{token.Assign, "="},
		{token.Float, "20.2"},
		{token.NewLine, "\n"},
		{token.If, "if"},
		{token.Identifier, "b"},
		{token.Greater, ">"},
		{token.Identifier, "a"},
		{token.Then, "then"},
		{token.NewLine, "\n"},
		{token.Repeat, "repeat"},
		{token.Identifier, "i"},
		{token.In, "in"},
		{token.Integer, "5"},
		{token.Range, ".."},
		{token.Integer, "10"},
		{token.NewLine, "\n"},
		{token.Identifier, "i"},
		{token.Plus, "+"},
		{token.Integer, "2"},
		{token.NewLine, "\n"},
		{token.End, "end"},
		{token.NewLine, "\n"},
		{token.Else, "else"},
		{token.NewLine, "\n"},
		{token.String, "exiting..."},
		{token.NewLine, "\n"},
		{token.End, "end"},
		{token.NewLine, "\n"},
		{token.Val, "val"},
		{token.Identifier, "c"},
		{token.Assign, "="},
		{token.Function, "fn"},
		{token.Identifier, "x"},
		{token.Comma, ","},
		{token.Identifier, "y"},
		{token.Comma, ","},
		{token.Identifier, "z"},
		{token.NewLine, "\n"},
		{token.String, "hi"},
		{token.Plus, "+"},
		{token.Identifier, "x"},
		{token.Plus, "+"},
		{token.Identifier, "y"},
		{token.Plus, "+"},
		{token.Identifier, "z"},
		{token.NewLine, "\n"},
		{token.End, "end"},
	}
	lex := New([]byte(input))
	for i, v := range tests {
		tok := lex.NextToken()
		if tok.Type != v.Type || tok.Literal != v.Literal {
			t.Errorf("Expected [%s %s] but got [%s %s] in line %d", string(v.Type), v.Literal, string(tok.Type), tok.Literal, i)
		}
	}
}
