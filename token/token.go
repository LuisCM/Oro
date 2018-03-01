// Copyright 2011 The LuisCM. All rights reserved.
// Use of this source code is license that can be found in the LICENSE file.

// Package token implements types and constants.
package token

type TType string

type Token struct {
	Type     TType
	Literal  string
	Position Position
}

type Position struct {
	Row int
	Col int
}

const (
	// Literals
	Identifier = "IDENTIFIER"
	Boolean    = "BOOLEAN"
	String     = "STRING"
	Integer    = "INTEGER"
	Float      = "FLOAT"
	// Operators
	Assign         = "="
	PlusAssign     = "+="
	MinusAssign    = "-="
	MultiplyAssign = "*="
	DivideAssign   = "/="
	Equal          = "=="
	NotEqual       = "!="
	Greater        = ">"
	GreaterEqual   = ">="
	Less           = "<"
	LessEqual      = "<="
	Plus           = "+"
	Minus          = "-"
	Multiply       = "*"
	Exponential    = "**"
	Modulus        = "%"
	Divide         = "/"
	BitwiseOr      = "|"
	BitwiseAnd     = "&"
	BitwiseNot     = "~"
	BitShiftLeft   = "<<"
	BitShiftRight  = ">>"
	LogicalOr      = "||"
	LogicalAnd     = "&&"
	Not            = "!"
	Pipe           = "|>"
	Arrow          = "->"
	FatArrow       = "=>"
	QuestionMark   = "?"
	// Delimiters
	Space            = " "
	NewLine          = "\n"
	LeftParenthesis  = "("
	LeftBracket      = "["
	LeftBraces       = "{"
	RightParenthesis = ")"
	RightBracket     = "]"
	RightBraces      = "}"
	Colon            = ":"
	SemiColon        = ";"
	Range            = ".."
	Ellipsis         = "..."
	Dot              = "."
	Comma            = ","
	Underscore       = "_"
	// Keywords
	Val      = "val"
	Var      = "var"
	Function = "fn"
	Do       = "do"
	End      = "end"
	Then     = "then"
	If       = "if"
	Else     = "else"
	Repeat   = "repeat"
	In       = "in"
	Is       = "is"
	As       = "as"
	Nil      = "nil"
	Return   = "return"
	Match    = "match"
	With     = "with"
	When     = "when"
	Break    = "break"
	Continue = "continue"
	Module   = "module"
	Use      = "use"
	True     = "true"
	False    = "false"
	// Miscellaneous
	Comment = "COMMENT"
	Eof     = "EOF"
)
