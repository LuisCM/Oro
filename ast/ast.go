// Copyright 2011 The LuisCM. All rights reserved.
// Use of this source code is license that can be found in the LICENSE file.

// Package ast implements functions to abstract syntax tree.
package ast

import (
	"bytes"
	"fmt"
	"github.com/luiscm/oro/token"
	"strings"
)

type Node interface {
	TokenLiteral() string
	TokenPosition() token.Position
	Check() string
}

type Statement interface {
	Node
	Statement()
}

type Expression interface {
	Node
	Expression()
}

type Program struct {
	Statements []Statement
}

type Boolean struct {
	Token token.Token
	Value bool
}

type String struct {
	Token        token.Token
	Value        string
	Interpolated map[string]Expression
}

type Integer struct {
	Token token.Token
	Value int64
}

type Float struct {
	Token token.Token
	Value float64
}

type Array struct {
	Token token.Token
	List  *ExpressionList
}

type Dictionary struct {
	Token token.Token
	Pairs map[Expression]Expression
}

type Symbol struct {
	Token token.Token
	Value string
}

type Nil struct {
	Token token.Token
}

type Val struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

type Var struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

type Is struct {
	Token token.Token
	Left  Expression
	Right *Identifier
}

type As struct {
	Token token.Token
	Left  Expression
	Right *Identifier
}

type If struct {
	Token     token.Token
	Condition Expression
	Then      *BlockStatement
	Else      *BlockStatement
}

type Repeat struct {
	Token      token.Token
	Arguments  *IdentifierList
	Enumerable Expression
	Body       *BlockStatement
}

type Match struct {
	Token   token.Token
	Control Expression
	Whens   []*MatchWhen
	Else    *BlockStatement
}

type MatchWhen struct {
	Token  token.Token
	Values *ExpressionList
	Body   *BlockStatement
}

type Break struct {
	Token token.Token
}

type Continue struct {
	Token token.Token
}

type Return struct {
	Token token.Token
	Value Expression
}

type Pipe struct {
	Token token.Token
	Left  Expression
	Right Expression
}

type PlaceHolder struct {
	Token token.Token
}

type Function struct {
	Token      token.Token
	Parameters []*FunctionParameter
	Body       *BlockStatement
	ReturnType *Identifier
	Variadic   bool
}

type FunctionParameter struct {
	Token   token.Token
	Name    *Identifier
	Type    *Identifier
	Default Expression
}

type FunctionCall struct {
	Token     token.Token
	Function  Expression
	Arguments *ExpressionList
}

type Module struct {
	Token token.Token
	Name  *Identifier
	Body  *BlockStatement
}

type ModuleAccess struct {
	Token     token.Token
	Object    *Identifier
	Parameter *Identifier
}

type Subscript struct {
	Token token.Token
	Left  Expression
	Index Expression
}

type Use struct {
	Token token.Token
	File  *String
}

type Assign struct {
	Token    token.Token
	Operator string
	Name     Expression
	Right    Expression
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

type ExpressionList struct {
	Token    token.Token
	Elements []Expression
}

type Identifier struct {
	Token token.Token
	Value string
}

type IdentifierList struct {
	Token    token.Token
	Elements []*Identifier
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (a *Program) TokenLiteral() string {
	if len(a.Statements) > 0 {
		return a.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (a *Program) TokenPosition() token.Position {
	return token.Position{}
}

func (a *Program) Check() string {
	var out bytes.Buffer
	for _, s := range a.Statements {
		out.WriteString(s.Check())
	}
	return out.String()
}

func (a *Boolean) Expression() {
}

func (a *Boolean) TokenLiteral() string {
	return a.Token.Literal
}

func (a *Boolean) TokenPosition() token.Position {
	return a.Token.Position
}

func (a *Boolean) Check() string {
	return a.Token.Literal
}

func (a *String) Expression() {
}

func (a *String) TokenLiteral() string {
	return a.Token.Literal
}

func (a *String) TokenPosition() token.Position {
	return a.Token.Position
}

func (a *String) Check() string {
	return a.Token.Literal
}

func (a *Integer) Expression() {
}

func (a *Integer) TokenLiteral() string {
	return a.Token.Literal
}

func (a *Integer) TokenPosition() token.Position {
	return a.Token.Position
}

func (a *Integer) Check() string {
	return a.Token.Literal
}

func (a *Float) Expression() {
}

func (a *Float) TokenLiteral() string {
	return a.Token.Literal
}

func (a *Float) TokenPosition() token.Position {
	return a.Token.Position
}

func (a *Float) Check() string {
	return a.Token.Literal
}

func (a *Array) Expression() {
}

func (a *Array) TokenLiteral() string {
	return a.Token.Literal
}

func (a *Array) TokenPosition() token.Position {
	return a.Token.Position
}

func (a *Array) Check() string {
	var out bytes.Buffer
	out.WriteString("Array")
	out.WriteString(token.LeftParenthesis)
	out.WriteString(a.List.Check())
	out.WriteString(token.RightParenthesis)
	return out.String()
}

func (a *Dictionary) Expression() {
}

func (a *Dictionary) TokenLiteral() string {
	return a.Token.Literal
}

func (a *Dictionary) TokenPosition() token.Position {
	return a.Token.Position
}

func (a *Dictionary) Check() string {
	var out bytes.Buffer
	var pairs []string
	for key, value := range a.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s => %s", key.Check(), value.Check()))
	}
	out.WriteString(token.LeftBracket)
	out.WriteString(strings.Join(pairs,", "))
	out.WriteString(token.RightBracket)
	return out.String()
}

func (a *Symbol) Expression() {
}

func (a *Symbol) TokenLiteral() string {
	return a.Token.Literal
}

func (a *Symbol) TokenPosition() token.Position {
	return a.Token.Position
}
func (a *Symbol) Check() string {
	var out bytes.Buffer
	out.WriteString(token.Colon)
	out.WriteString(a.Token.Literal)
	return out.String()
}

func (a *Nil) Expression() {
}

func (a *Nil) TokenLiteral() string {
	return a.Token.Literal
}

func (a *Nil) TokenPosition() token.Position {
	return a.Token.Position
}

func (a *Nil) Check() string {
	return a.Token.Literal
}

func (a *Val) Expression() {
}

func (a *Val) TokenLiteral() string {
	return a.Token.Literal
}

func (a *Val) TokenPosition() token.Position {
	return a.Token.Position
}

func (a *Val) Check() string {
	var out bytes.Buffer
	out.WriteString(token.Val)
	out.WriteString(token.Space)
	out.WriteString(a.Name.Check())
	out.WriteString(token.Space)
	out.WriteString(token.Assign)
	out.WriteString(token.Space)
	if a.Value != nil {
		out.WriteString(a.Value.Check())
	}
	return out.String()
}

func (a *Var) Expression() {
}

func (a *Var) TokenLiteral() string {
	return a.Token.Literal
}

func (a *Var) TokenPosition() token.Position {
	return a.Token.Position
}

func (a *Var) Check() string {
	var out bytes.Buffer
	out.WriteString(token.Var)
	out.WriteString(token.Space)
	out.WriteString(a.Name.Check())
	out.WriteString(token.Space)
	out.WriteString(token.Assign)
	out.WriteString(token.Space)
	if a.Value != nil {
		out.WriteString(a.Value.Check())
	}
	return out.String()
}

func (a *Is) Expression() {
}

func (a *Is) TokenLiteral() string {
	return a.Token.Literal
}

func (a *Is) TokenPosition() token.Position {
	return a.Token.Position
}

func (a *Is) Check() string {
	var out bytes.Buffer
	out.WriteString(token.LeftParenthesis)
	out.WriteString(a.Left.Check())
	out.WriteString(token.Space)
	out.WriteString(token.Is)
	out.WriteString(token.Space)
	out.WriteString(a.Right.Check())
	out.WriteString(token.RightParenthesis)
	return out.String()
}

func (a *As) Expression() {
}

func (a *As) TokenLiteral() string {
	return a.Token.Literal
}

func (a *As) TokenPosition() token.Position {
	return a.Token.Position
}

func (a *As) Check() string {
	var out bytes.Buffer
	out.WriteString(token.LeftParenthesis)
	out.WriteString(a.Left.Check())
	out.WriteString(token.Space)
	out.WriteString(token.As)
	out.WriteString(token.Space)
	out.WriteString(a.Right.Check())
	out.WriteString(token.RightParenthesis)
	return out.String()
}

func (a *If) Expression() {
}

func (a *If) TokenLiteral() string {
	return a.Token.Literal
}

func (a *If) TokenPosition() token.Position {
	return a.Token.Position
}

func (a *If) Check() string {
	var out bytes.Buffer
	out.WriteString(token.If)
	out.WriteString(token.Space)
	out.WriteString(a.Condition.Check())
	out.WriteString(token.Space)
	out.WriteString(token.Then)
	out.WriteString(token.Space)
	out.WriteString(a.Then.Check())
	if a.Else != nil {
		out.WriteString(token.Space)
		out.WriteString(token.Else)
		out.WriteString(token.Space)
		out.WriteString(a.Else.Check())
	}
	return out.String()
}

func (a *Repeat) Expression() {
}

func (a *Repeat) TokenLiteral() string {
	return a.Token.Literal
}

func (a *Repeat) TokenPosition() token.Position {
	return a.Token.Position
}

func (a *Repeat) Check() string {
	var out bytes.Buffer
	out.WriteString(a.Token.Literal)
	if a.Enumerable != nil {
		out.WriteString(token.Space)
		out.WriteString(token.LeftParenthesis)
		if a.Arguments != nil {
			out.WriteString(a.Arguments.Check())
			out.WriteString(token.Space)
			out.WriteString(token.In)
			out.WriteString(token.Space)
		}
		out.WriteString(a.Enumerable.Check())
		out.WriteString(token.RightParenthesis)
	}
	out.WriteString(token.Space)
	out.WriteString(token.Arrow)
	out.WriteString(token.Space)
	out.WriteString(a.Body.Check())
	return out.String()
}

func (a *Match) Expression() {
}

func (a *Match) TokenLiteral() string {
	return a.Token.Literal
}

func (a *Match) TokenPosition() token.Position {
	return a.Token.Position
}

func (a *Match) Check() string {
	var out bytes.Buffer
	var whens []string
	for _, w := range a.Whens {
		whens = append(whens, w.Check())
	}
	out.WriteString(token.Match)
	out.WriteString(token.Space)
	if a.Control != nil {
		out.WriteString(a.Control.Check())
	}
	out.WriteString(token.Space)
	out.WriteString(token.Arrow)
	out.WriteString(token.Space)
	out.WriteString(strings.Join(whens, "; "))
	if a.Else != nil {
		out.WriteString(token.SemiColon)
		out.WriteString(token.Space)
		out.WriteString(token.Else)
		out.WriteString(token.Space)
		out.WriteString(a.Else.Check())
	}
	return out.String()
}

func (a *MatchWhen) Expression() {
}

func (a *MatchWhen) TokenLiteral() string {
	return a.Token.Literal
}

func (a *MatchWhen) TokenPosition() token.Position {
	return a.Token.Position
}

func (a *MatchWhen) Check() string {
	var out bytes.Buffer
	out.WriteString(token.When)
	out.WriteString(token.Space)
	out.WriteString(a.Values.Check())
	out.WriteString(token.Space)
	out.WriteString(token.Then)
	out.WriteString(token.Space)
	out.WriteString(a.Body.Check())
	return out.String()
}

func (a *Break) Statement() {
}

func (a *Break) TokenLiteral() string {
	return a.Token.Literal
}

func (a *Break) TokenPosition() token.Position {
	return a.Token.Position
}

func (a *Break) Check() string {
	return a.Token.Literal
}

func (a *Continue) Statement() {
}

func (a *Continue) TokenLiteral() string {
	return a.Token.Literal
}

func (a *Continue) TokenPosition() token.Position {
	return a.Token.Position
}

func (a *Continue) Check() string {
	return a.Token.Literal
}

func (a *Return) Statement() {
}

func (a *Return) TokenLiteral() string {
	return a.Token.Literal
}

func (a *Return) TokenPosition() token.Position {
	return a.Token.Position
}

func (a *Return) Check() string {
	var out bytes.Buffer
	out.WriteString(a.Token.Literal)
	out.WriteString(token.Space)
	if a.Value != nil {
		out.WriteString(a.Value.Check())
	}
	return out.String()
}

func (a *Pipe) Expression() {
}

func (a *Pipe) TokenLiteral() string {
	return a.Token.Literal
}

func (a *Pipe) TokenPosition() token.Position {
	return a.Token.Position
}

func (a *Pipe) Check() string {
	var out bytes.Buffer
	out.WriteString(a.Left.Check())
	out.WriteString(token.Space)
	out.WriteString(token.Pipe)
	out.WriteString(token.Space)
	out.WriteString(a.Right.Check())
	return out.String()
}

func (a *PlaceHolder) Expression() {
}

func (a *PlaceHolder) TokenLiteral() string {
	return a.Token.Literal
}

func (a *PlaceHolder) TokenPosition() token.Position {
	return a.Token.Position
}

func (a *PlaceHolder) Check() string {
	return ""
}

func (a *Function) Expression() {
}

func (a *Function) TokenLiteral() string {
	return a.Token.Literal
}

func (a *Function) TokenPosition() token.Position {
	return a.Token.Position
}

func (a *Function) Check() string {
	var out bytes.Buffer
	var parameters []string
	for i, v := range a.Parameters {
		param := v.Check()
		if a.Variadic && i == len(a.Parameters)-1 {
			param = token.Ellipsis + param
		}
		parameters = append(parameters, param)
	}
	out.WriteString(a.Token.Literal)
	out.WriteString(token.Space)
	out.WriteString(token.LeftParenthesis)
	out.WriteString(strings.Join(parameters,", "))
	out.WriteString(token.RightParenthesis)
	out.WriteString(token.Space)
	if a.ReturnType != nil {
		out.WriteString(token.Space)
		out.WriteString(token.Arrow)
		out.WriteString(token.Space)
		out.WriteString(a.ReturnType.Value)
		out.WriteString(token.NewLine)
	}
	out.WriteString(a.Body.Check())
	return out.String()
}

func (a *FunctionParameter) Expression() {
}

func (a *FunctionParameter) TokenLiteral() string {
	return a.Token.Literal
}

func (a *FunctionParameter) TokenPosition() token.Position {
	return a.Token.Position
}

func (a *FunctionParameter) Check() string {
	var out bytes.Buffer
	out.WriteString(a.Name.Value)
	if a.Type != nil {
		out.WriteString(token.Colon)
		out.WriteString(a.Type.Value)
	}
	if a.Default != nil {
		out.WriteString(token.Space)
		out.WriteString(token.Assign)
		out.WriteString(token.Space)
		out.WriteString(a.Default.Check())
	}
	return out.String()
}

func (a *FunctionCall) Expression() {
}

func (a *FunctionCall) TokenLiteral() string {
	return a.Token.Literal
}

func (a *FunctionCall) TokenPosition() token.Position {
	return a.Token.Position
}

func (a *FunctionCall) Check() string {
	var out bytes.Buffer
	out.WriteString(a.Function.Check())
	out.WriteString(token.LeftParenthesis)
	out.WriteString(a.Arguments.Check())
	out.WriteString(token.RightParenthesis)
	return out.String()
}

func (a *Module) Expression() {
}

func (a *Module) TokenLiteral() string {
	return a.Token.Literal
}

func (a *Module) TokenPosition() token.Position {
	return a.Token.Position
}

func (a *Module) Check() string {
	var out bytes.Buffer
	out.WriteString(token.Module)
	out.WriteString(token.Space)
	out.WriteString(a.Name.Check())
	out.WriteString(token.Space)
	out.WriteString(token.LeftBraces)
	out.WriteString(token.Space)
	out.WriteString(a.Body.Check())
	out.WriteString(token.Space)
	out.WriteString(token.RightBraces)
	out.WriteString(token.Space)
	return out.String()
}

func (a *ModuleAccess) Expression() {
}

func (a *ModuleAccess) TokenLiteral() string {
	return a.Token.Literal
}

func (a *ModuleAccess) TokenPosition() token.Position {
	return a.Token.Position
}

func (a *ModuleAccess) Check() string {
	var out bytes.Buffer
	out.WriteString(a.Object.Check())
	out.WriteString(token.Arrow)
	out.WriteString(a.Parameter.Check())
	return out.String()
}

func (a *Subscript) Expression() {
}

func (a *Subscript) TokenLiteral() string {
	return a.Token.Literal
}

func (a *Subscript) TokenPosition() token.Position {
	return a.Token.Position
}

func (a *Subscript) Check() string {
	var out bytes.Buffer
	out.WriteString(a.Left.Check())
	out.WriteString(token.LeftBracket)
	out.WriteString(a.Index.Check())
	out.WriteString(token.RightBracket)
	return out.String()
}

func (a *Use) Expression() {
}

func (a *Use) TokenLiteral() string {
	return a.Token.Literal
}

func (a *Use) TokenPosition() token.Position {
	return a.Token.Position
}

func (a *Use) Check() string {
	var out *bytes.Buffer
	out.WriteString(token.Use)
	out.WriteString(token.Space)
	out.WriteString(a.File.Value)
	return out.String()
}

func (a *Assign) Expression() {
}

func (a *Assign) TokenLiteral() string {
	return a.Token.Literal
}

func (a *Assign) TokenPosition() token.Position {
	return a.Token.Position
}

func (a *Assign) Check() string {
	var out bytes.Buffer
	out.WriteString(a.Name.Check())
	out.WriteString(token.Space)
	out.WriteString(token.Assign)
	out.WriteString(token.Space)
	out.WriteString(a.Right.Check())
	return out.String()
}

func (a *BlockStatement) Statement() {
}

func (a *BlockStatement) TokenLiteral() string {
	return a.Token.Literal
}

func (a *BlockStatement) TokenPosition() token.Position {
	return a.Token.Position
}

func (a *BlockStatement) Check() string {
	var out bytes.Buffer
	for _, s := range a.Statements {
		out.WriteString(s.Check())
	}
	return out.String()
}

func (a *ExpressionStatement) Statement() {
}

func (a *ExpressionStatement) TokenLiteral() string {
	return a.Token.Literal
}

func (a *ExpressionStatement) TokenPosition() token.Position {
	return a.Token.Position
}

func (a *ExpressionStatement) Check() string {
	if a.Expression != nil {
		return a.Expression.Check()
	}
	return ""
}

func (a *ExpressionList) Expression() {
}

func (a *ExpressionList) TokenLiteral() string {
	return a.Token.Literal
}

func (a *ExpressionList) TokenPosition() token.Position {
	return a.Token.Position
}

func (a *ExpressionList) Check() string {
	var out bytes.Buffer
	var elements []string
	for _, el := range a.Elements {
		elements = append(elements, el.Check())
	}
	out.WriteString(strings.Join(elements,", "))
	return out.String()
}

func (a *Identifier) Expression() {
}

func (a *Identifier) TokenLiteral() string {
	return a.Token.Literal
}

func (a *Identifier) TokenPosition() token.Position {
	return a.Token.Position
}

func (a *Identifier) Check() string {
	return a.Value
}

func (a *IdentifierList) Statement() {
}

func (a *IdentifierList) TokenLiteral() string {
	return a.Token.Literal
}

func (a *IdentifierList) TokenPosition() token.Position {
	return a.Token.Position
}

func (a *IdentifierList) Check() string {
	var out bytes.Buffer
	var elements []string
	for _, el := range a.Elements {
		elements = append(elements, el.Check())
	}
	out.WriteString(strings.Join(elements,", "))
	return out.String()
}

func (a *PrefixExpression) Expression() {
}

func (a *PrefixExpression) TokenLiteral() string {
	return a.Token.Literal
}

func (a *PrefixExpression) TokenPosition() token.Position {
	return a.Token.Position
}

func (a *PrefixExpression) Check() string {
	var out bytes.Buffer
	out.WriteString(token.LeftParenthesis)
	out.WriteString(a.Operator)
	out.WriteString(a.Right.Check())
	out.WriteString(token.RightParenthesis)
	return out.String()
}

func (a *InfixExpression) Expression() {
}

func (a *InfixExpression) TokenLiteral() string {
	return a.Token.Literal
}

func (a *InfixExpression) TokenPosition() token.Position {
	return a.Token.Position
}

func (a *InfixExpression) Check() string {
	var out bytes.Buffer
	out.WriteString(token.LeftParenthesis)
	out.WriteString(a.Left.Check())
	out.WriteString(token.Space)
	out.WriteString(a.Operator)
	out.WriteString(token.Space)
	out.WriteString(a.Right.Check())
	out.WriteString(token.RightParenthesis)
	return out.String()
}
