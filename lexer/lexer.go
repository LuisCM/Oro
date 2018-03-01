// Copyright 2011 The LuisCM. All rights reserved.
// Use of this source code is license that can be found in the LICENSE file.

// Package lexer implements functions to lexer.
package lexer

import (
	"bytes"
	"fmt"
	"github.com/luiscm/oro/rerror"
	"github.com/luiscm/oro/runtime/cmd"
	"github.com/luiscm/oro/token"
)

type Lexer struct {
	buff    []byte
	offset  int
	curr    int
	nextr   int
	row     int
	col     int
	chr     rune
	token   token.Token
	command *cmd.Command
}

func New(buffer []byte) *Lexer {
	l := &Lexer{
		buff:    buffer,
		row:     1,
		col:     1,
		command: &cmd.Command{},
	}
	l.command.InsertAll()
	l.next()
	return l
}

func (l *Lexer) NextToken() token.Token {
	l.skipWhiteSpace()
	switch {
	case l.chr == 0:
		l.assignToken(token.Eof, "")
	case l.chr == '#':
		l.next()
		l.skipComment()
	case l.chr == '=':
		switch l.peek() {
		case '=':
			l.next()
			l.assignToken(token.Equal, token.Equal)
		case '>':
			l.next()
			l.assignToken(token.FatArrow, token.FatArrow)
		default:
			l.assignToken(token.Assign, string(l.chr))
		}
	case l.chr == '>':
		switch l.peek() {
		case '=':
			l.next()
			l.assignToken(token.GreaterEqual, token.GreaterEqual)
		case '>':
			l.next()
			l.assignToken(token.BitShiftRight, token.BitShiftRight)
		default:
			l.assignToken(token.Greater, string(l.chr))
		}
	case l.chr == '<':
		switch l.peek() {
		case '=':
			l.next()
			l.assignToken(token.LessEqual, token.LessEqual)
		case '<':
			l.next()
			l.assignToken(token.BitShiftLeft, token.BitShiftLeft)
		default:
			l.assignToken(token.Less, string(l.chr))
		}
	case l.chr == '+':
		switch l.peek() {
		case '=':
			l.next()
			l.assignToken(token.PlusAssign, token.PlusAssign)
		default:
			l.assignToken(token.Plus, string(l.chr))
		}
	case l.chr == '-':
		switch l.peek() {
		case '>':
			l.next()
			l.assignToken(token.Arrow, token.Arrow)
		case '=':
			l.next()
			l.assignToken(token.MinusAssign, token.MinusAssign)
		default:
			l.assignToken(token.Minus, string(l.chr))
		}
	case l.chr == '*':
		switch l.peek() {
		case '*':
			l.next()
			l.assignToken(token.Exponential, token.Exponential)
		case '=':
			l.next()
			l.assignToken(token.MultiplyAssign, token.MultiplyAssign)
		default:
			l.assignToken(token.Multiply, string(l.chr))
		}
	case l.chr == '/':
		switch l.peek() {
		case '/':
			l.next()
			l.skipComment()
		case '*':
			l.next()
			l.skipMultiLineComment()
		case '=':
			l.next()
			l.assignToken(token.DivideAssign, token.DivideAssign)
		default:
			l.assignToken(token.Divide, string(l.chr))
		}
	case l.chr == '%':
		l.assignToken(token.Modulus, string(l.chr))
	case l.chr == ',':
		l.assignToken(token.Comma, string(l.chr))
	case l.chr == '.':
		switch l.peek() {
		case '.':
			l.next()
			switch l.peek() {
			case '.':
				l.next()
				l.assignToken(token.Ellipsis, token.Ellipsis)
			default:
				l.assignToken(token.Range, token.Range)
			}
		default:
			l.assignToken(token.Dot, string(l.chr))
		}
	case l.chr == '|':
		switch l.peek() {
		case '|':
			l.next()
			l.assignToken(token.LogicalOr, token.LogicalOr)
		case '>':
			l.next()
			l.assignToken(token.Pipe, token.Pipe)
		default:
			l.assignToken(token.BitwiseOr, string(l.chr))
		}
	case l.chr == '&':
		switch l.peek() {
		case '&':
			l.next()
			l.assignToken(token.LogicalAnd, token.LogicalAnd)
		default:
			l.assignToken(token.BitwiseAnd, string(l.chr))
		}
	case l.chr == '~':
		l.assignToken(token.BitwiseNot, string(l.chr))
	case l.chr == '!':
		switch l.peek() {
		case '=':
			l.next()
			l.assignToken(token.NotEqual, token.NotEqual)
		default:
			l.assignToken(token.Not, string(l.chr))
		}
	case l.chr == '(':
		l.assignToken(token.LeftParenthesis, string(token.LeftParenthesis))
	case l.chr == ')':
		l.assignToken(token.RightParenthesis, string(token.RightParenthesis))
	case l.chr == '[':
		l.assignToken(token.LeftBracket, string(token.LeftBracket))
	case l.chr == ']':
		l.assignToken(token.RightBracket, string(token.RightBracket))
	case l.chr == '?':
		l.assignToken(token.QuestionMark, string(token.QuestionMark))
	case l.chr == ':':
		l.assignToken(token.Colon, string(token.Colon))
	case l.chr == '_':
		l.assignToken(token.Underscore, string(token.Underscore))
	case l.chr == '\n':
		l.assignToken(token.NewLine, string(token.NewLine))
	case l.chr == '"':
		l.skipString()
	case l.chr == '0' && l.peek() == 'x':
		l.skipSpecialInteger(l.isHex)
	case l.chr == '0' && l.peek() == 'o':
		l.skipSpecialInteger(l.isOctal)
	case l.chr == '0' && l.peek() == 'b':
		l.skipSpecialInteger(l.isBinary)
	case l.isNumber(l.chr):
		l.skipNumeric()
	default:
		if l.isAlpha(l.chr) {
			l.skipIdentifier()
		} else {
			l.reportError(fmt.Sprintf("Unidentified character '%s'", string(l.chr)))
		}
	}
	l.next()
	return l.token
}

func (l *Lexer) next() rune {
	if l.nextr >= len(l.buff) {
		l.chr = 0
	} else {
		l.chr = rune(l.buff[l.nextr])
	}
	l.curr = l.nextr
	l.offset = l.curr
	l.nextr++
	l.col++
	if l.chr == '\n' {
		l.row++
		l.col = 0
	}
	return l.chr
}

func (l *Lexer) peek() rune {
	if l.nextr >= len(l.buff) {
		return 0
	}
	return rune(l.buff[l.nextr])
}

func (l *Lexer) rewind() {
	if l.nextr >= len(l.buff) {
		l.chr = 0
	} else {
		l.chr = rune(l.buff[l.curr])
	}
	l.nextr = l.curr
	l.offset = l.nextr
	l.col--
}

func (l *Lexer) assignToken(tokenType token.TType, value string) {
	l.token = token.Token{
		Type:     tokenType,
		Literal:  value,
		Position: token.Position{Row: l.row, Col: l.col},
	}
}

func (l *Lexer) isAlpha(char rune) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') ||
		(char >= '0' && char <= '9') || char == '_' || char == '!' || char == '?'
}

func (l *Lexer) isNumber(char rune) bool {
	return char >= '0' && char <= '9'
}

func (l *Lexer) isHex(char rune) bool {
	return l.isNumber(char) || (char >= 'a' && char <= 'f' || char >= 'A' && char <= 'F')
}

func (l *Lexer) isOctal(char rune) bool {
	return char >= '0' && char <= '7'
}

func (l *Lexer) isBinary(char rune) bool {
	return char == '0' || char == '1'
}

func (l *Lexer) readName() string {
	var out bytes.Buffer
	out.WriteRune(l.chr)
	for l.isAlpha(l.peek()) {
		l.next()
		out.WriteRune(l.chr)
	}
	return out.String()
}

func (l *Lexer) skipWhiteSpace() {
	for l.chr == ' ' || l.chr == '\t' || l.chr == '\r' {
		l.next()
	}
}

func (l *Lexer) skipIdentifier() {
	identifier := l.readName()
	if tokType, found := l.command.Lookup(identifier); found {
		l.assignToken(tokType, identifier)
	} else {
		l.assignToken(token.Identifier, identifier)
	}
}

func (l *Lexer) skipString() {
	var out bytes.Buffer
	l.next()
loop:
	for {
		switch l.chr {
		case '\\':
			l.next()
			switch l.chr {
			case '"':
				out.WriteRune('\\')
				out.WriteRune('"')
			case '\\':
				out.WriteRune('\\')
			case 'n', 't', 'r', 'a', 'b', 'f', 'v':
				out.WriteRune('\\')
				out.WriteRune(l.chr)
			default:
				l.reportError(fmt.Sprintf("Invalid escape character '%s", string(l.chr)))
			}
		case 0:
			l.reportError("Unterminated string")
			break loop
		case '"':
			break loop
		default:
			out.WriteRune(l.chr)
		}
		l.next()
	}
	l.assignToken(token.String, out.String())
}

func (l *Lexer) skipNumeric() {
	var out bytes.Buffer
	out.WriteRune(l.chr)
	floatFound := false
	scientificFound := false
loop:
	for {
		l.next()
		switch {
		case l.isNumber(l.chr):
			out.WriteRune(l.chr)
		case l.chr == '_':
		case l.chr == '.' && l.isNumber(l.peek()):
			floatFound = true
			out.WriteRune('.')
		case l.chr == 'e' && (l.isNumber(l.peek()) || l.peek() == '-'):
			floatFound = true
			scientificFound = true
			out.WriteRune('e')
		case l.chr == '-' && scientificFound:
			out.WriteRune('-')
		case l.chr == '.' && l.peek() == '.':
			l.rewind()
			break loop
		case l.chr == 0:
			break loop
		default:
			l.rewind()
			break loop
		}
	}
	if floatFound {
		l.assignToken(token.Float, out.String())
	} else {
		l.assignToken(token.Integer, out.String())
	}
}

func (l *Lexer) skipSpecialInteger(fn func(rune) bool) {
	var out bytes.Buffer
	out.WriteRune(l.chr)
	out.WriteRune(l.peek())
	l.next()
	for fn(l.peek()) {
		out.WriteRune(l.peek())
		l.next()
	}
	ret := out.String()
	if len(ret) == 2 {
		l.reportError(fmt.Sprintf("Literal sequence '%s' started but not continued", ret))
	}
	l.assignToken(token.Integer, ret)
}

func (l *Lexer) skipComment() {
	var out bytes.Buffer
	l.next()
loop:
	for {
		switch l.chr {
		case '\n', 0:
			break loop
		case '\r':
			l.next()
			switch l.chr {
			case '\n', 0:
				break loop
			default:
				l.reportError("Unexpected comment line ending")
				break loop
			}
		default:
			out.WriteRune(l.chr)
		}
		l.next()
	}
	l.assignToken(token.Comment, out.String())
}

func (l *Lexer) skipMultiLineComment() {
	var out bytes.Buffer
loop:
	for {
		l.next()
		switch l.chr {
		case '*':
			switch l.peek() {
			case '/':
				l.next()
				break loop
			}
		case 0:
			l.reportError("Unterminated multi line comment")
			break loop
		default:
			out.WriteRune(l.chr)
		}
	}
	l.assignToken(token.Comment, out.String())
}

func (l *Lexer) reportError(msg string) {
	rerror.Error(rerror.Parse, token.Position{l.row, l.col}, msg)
}
