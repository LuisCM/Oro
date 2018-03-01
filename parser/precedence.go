// Copyright 2011 The LuisCM. All rights reserved.
// Use of this source code is license that can be found in the LICENSE file.

// Package parser implements functions to precedence.
package parser

import (
	"github.com/luiscm/oro/token"
)

const (
	_ = iota
	Lowest
	Assign
	Pipe
	Arrow
	Ternary
	Boolean
	Bitwise
	Comparison
	Range
	BitShift
	Sum
	Product
	Exponential
	Prefix
	Call
	Index
	As
)

var precedences = map[token.TType]int{
	token.Assign:          Assign,
	token.PlusAssign:      Assign,
	token.MinusAssign:     Assign,
	token.MultiplyAssign:  Assign,
	token.DivideAssign:    Assign,
	token.Plus:            Sum,
	token.Minus:           Sum,
	token.Multiply:        Product,
	token.Divide:          Product,
	token.Modulus:         Product,
	token.Exponential:     Exponential,
	token.Equal:           Comparison,
	token.NotEqual:        Comparison,
	token.Less:            Comparison,
	token.LessEqual:       Comparison,
	token.GreaterEqual:    Comparison,
	token.Greater:         Comparison,
	token.LogicalOr:       Boolean,
	token.LogicalAnd:      Boolean,
	token.Dot:             Call,
	token.LeftParenthesis: Call,
	token.LeftBracket:     Index,
	token.BitwiseOr:       Bitwise,
	token.BitwiseAnd:      Bitwise,
	token.BitwiseNot:      Bitwise,
	token.BitShiftLeft:    BitShift,
	token.BitShiftRight:   BitShift,
	token.Range:           Range,
	token.Pipe:            Pipe,
	token.Arrow:           Arrow,
	token.QuestionMark:    Ternary,
	token.Is:              Assign,
	token.As:              As,
}
