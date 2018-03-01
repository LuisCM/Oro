// Copyright 2011 The LuisCM. All rights reserved.
// Use of this source code is license that can be found in the LICENSE file.

// Package parser implements functions to parser.
package parser

import (
	"fmt"
	"github.com/luiscm/oro/ast"
	"github.com/luiscm/oro/lexer"
	"github.com/luiscm/oro/rerror"
	"github.com/luiscm/oro/token"
	"reflect"
	"strconv"
	"strings"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	lexer          *lexer.Lexer
	token          token.Token
	peekToken      token.Token
	prefixFunction map[token.TType]prefixParseFn
	infixFunction  map[token.TType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	parser := &Parser{lexer: l}
	parser.prefixFunction = make(map[token.TType]prefixParseFn)
	parser.infixFunction = make(map[token.TType]infixParseFn)
	parser.prefix(token.Val, parser.parseVal)
	parser.prefix(token.Var, parser.parseVar)
	parser.prefix(token.Module, parser.parseModule)
	parser.prefix(token.If, parser.parseIf)
	parser.prefix(token.Match, parser.parseMatch)
	parser.prefix(token.Repeat, parser.parseRepeat)
	parser.prefix(token.Function, parser.parseFunction)
	parser.prefix(token.Use, parser.parseUse)
	parser.prefix(token.LeftBracket, parser.parseArrayOrDictionary)
	parser.prefix(token.Identifier, parser.parseIdentifier)
	parser.prefix(token.Integer, parser.parseInteger)
	parser.prefix(token.Float, parser.parseFloat)
	parser.prefix(token.String, parser.parseString)
	parser.prefix(token.Boolean, parser.parseBoolean)
	parser.prefix(token.Nil, parser.parseNil)
	parser.prefix(token.Underscore, parser.parsePlaceHolder)
	parser.prefix(token.Colon, parser.parseSymbol)
	parser.prefix(token.Not, parser.parsePrefix)
	parser.prefix(token.BitwiseNot, parser.parsePrefix)
	parser.prefix(token.Minus, parser.parsePrefix)
	parser.prefix(token.LeftParenthesis, parser.parseGroup)
	parser.infix(token.Assign, parser.parseAssign)
	parser.infix(token.PlusAssign, parser.parseAssign)
	parser.infix(token.MinusAssign, parser.parseAssign)
	parser.infix(token.MultiplyAssign, parser.parseAssign)
	parser.infix(token.DivideAssign, parser.parseAssign)
	parser.infix(token.Dot, parser.parseModuleAccess)
	parser.infix(token.LeftParenthesis, parser.parseFunctionCall)
	parser.infix(token.LeftBracket, parser.parseSubscript)
	parser.infix(token.Pipe, parser.parsePipe)
	parser.infix(token.Arrow, parser.parseArrowFunction)
	parser.infix(token.QuestionMark, parser.parseTernary)
	parser.infix(token.Is, parser.parseIs)
	parser.infix(token.As, parser.parseAs)
	parser.infix(token.Range, parser.parseInfix)
	parser.infix(token.Plus, parser.parseInfix)
	parser.infix(token.Minus, parser.parseInfix)
	parser.infix(token.Divide, parser.parseInfix)
	parser.infix(token.Multiply, parser.parseInfix)
	parser.infix(token.Modulus, parser.parseInfix)
	parser.infix(token.Exponential, parser.parseInfix)
	parser.infix(token.Equal, parser.parseInfix)
	parser.infix(token.NotEqual, parser.parseInfix)
	parser.infix(token.Less, parser.parseInfix)
	parser.infix(token.LessEqual, parser.parseInfix)
	parser.infix(token.Greater, parser.parseInfix)
	parser.infix(token.GreaterEqual, parser.parseInfix)
	parser.infix(token.LogicalOr, parser.parseInfixRight)
	parser.infix(token.LogicalAnd, parser.parseInfixRight)
	parser.infix(token.BitwiseAnd, parser.parseInfix)
	parser.infix(token.BitwiseOr, parser.parseInfix)
	parser.infix(token.BitShiftLeft, parser.parseInfix)
	parser.infix(token.BitShiftRight, parser.parseInfix)
	parser.nextToken()
	parser.nextToken()
	return parser
}

func (p *Parser) prefix(tokenType token.TType, fn prefixParseFn) {
	p.prefixFunction[tokenType] = fn
}

func (p *Parser) infix(tokenType token.TType, fn infixParseFn) {
	p.infixFunction[tokenType] = fn
}

func (p *Parser) precedence() int {
	if parser, ok := precedences[p.token.Type]; ok {
		return parser
	}
	return Lowest
}

func (p *Parser) peekPrecedence() int {
	if parser, ok := precedences[p.peekToken.Type]; ok {
		return parser
	}
	return Lowest
}

func (p *Parser) nextToken() {
	p.token = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) matchToken(tokenType ...token.TType) bool {
	for _, t := range tokenType {
		if p.token.Type == t {
			return true
		}
	}
	return false
}

func (p *Parser) peekTokenMatch(tokenType ...token.TType) bool {
	for _, t := range tokenType {
		if p.peekToken.Type == t {
			return true
		}
	}
	return false
}

func (p *Parser) Parse() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for !p.matchToken(token.Eof) {
		statement := p.parseStatement()
		if reflect.ValueOf(statement).IsValid() {
			program.Statements = append(program.Statements, statement)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.token.Type {
	case token.Comment, token.NewLine:
		return nil
	case token.Break:
		return p.parseBreak()
	case token.Continue:
		return p.parseContinue()
	case token.Return:
		return p.parseReturn()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseVal() ast.Expression {
	expression := &ast.Val{Token: p.token}
	if !p.peekTokenMatch(token.Identifier) {
		p.parserError("VAL expects an identifier")
		return nil
	}
	p.nextToken()
	expression.Name = &ast.Identifier{Token: p.token, Value: p.token.Literal}
	if !p.peekTokenMatch(token.Assign) {
		p.parserError("Missing assignment in VAL")
		return nil
	}
	p.nextToken()
	p.nextToken()
	expression.Value = p.parseExpression(Lowest)
	return expression
}

func (p *Parser) parseVar() ast.Expression {
	expression := &ast.Var{Token: p.token}
	if !p.peekTokenMatch(token.Identifier) {
		p.parserError("VAR expects an identifier")
		return nil
	}
	p.nextToken()
	expression.Name = &ast.Identifier{Token: p.token, Value: p.token.Literal}
	if !p.peekTokenMatch(token.Assign) {
		p.parserError("Missing assignment in VAR")
		return nil
	}
	p.nextToken()
	p.nextToken()
	expression.Value = p.parseExpression(Lowest)
	return expression
}

func (p *Parser) parseModule() ast.Expression {
	expression := &ast.Module{Token: p.token}
	if !p.peekTokenMatch(token.Identifier) {
		p.parserError("Expecting an identifier as MODULE name")
		return nil
	}
	p.nextToken()
	expression.Name = &ast.Identifier{Token: p.token, Value: p.token.Literal}
	if p.matchToken(token.Do) {
		p.nextToken()
	}
	expression.Body = p.parseBlockBody()
	if !p.matchToken(token.End) {
		p.parserError("Missing END closing statement in MODULE")
		return nil
	}
	return expression
}

func (p *Parser) parseModuleAccess(left ast.Expression) ast.Expression {
	switch object := left.(type) {
	case *ast.Identifier:
		expression := &ast.ModuleAccess{Token: p.token, Object: object}
		if p.peekTokenMatch(token.Identifier) {
			p.nextToken()
			expression.Parameter = &ast.Identifier{Token: p.token, Value: p.token.Literal}
			return expression
		}
	default:
		p.parserError(fmt.Sprintf("Cannot use '%s' as MODULE caller", object.TokenLiteral()))
	}
	return nil
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.token, Value: p.token.Literal}
}

func (p *Parser) parseInteger() ast.Expression {
	literalInteger := &ast.Integer{Token: p.token}
	literal := p.token.Literal
	var value int64
	var err error
	if strings.HasPrefix(literal, "0b") {
		value, err = strconv.ParseInt(literal[2:], 2, 64)
	} else if strings.HasPrefix(literal, "0x") {
		value, err = strconv.ParseInt(literal, 0, 64)
	} else if strings.HasPrefix(literal, "0o") {
		value, err = strconv.ParseInt(literal[:1]+literal[2:], 0, 64)
	} else {
		value, err = strconv.ParseInt(literal, 0, 64)
	}
	if err != nil {
		p.parserError(fmt.Sprintf("Couldn't parse %s as Integer", literal))
		return nil
	}
	literalInteger.Value = value
	return literalInteger
}

func (p *Parser) parseFloat() ast.Expression {
	literalFloat := &ast.Float{Token: p.token}
	literal := p.token.Literal
	value, err := strconv.ParseFloat(literal, 64)
	if err != nil {
		p.parserError(fmt.Sprintf("Couldn't parse %s as Float", literal))
		return nil
	}
	literalFloat.Value = value
	return literalFloat
}

func (p *Parser) parseString() ast.Expression {
	return &ast.String{Token: p.token, Value: p.token.Literal}
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.token, Value: p.token.Literal == token.True}
}

func (p *Parser) parseNil() ast.Expression {
	return &ast.Nil{Token: p.token}
}

func (p *Parser) parseIf() ast.Expression {
	expression := &ast.If{Token: p.token}
	p.nextToken()
	expression.Condition = p.parseExpression(Lowest)
	if expression.Condition == nil {
		p.parserError("Missing condition expression in IF")
		return nil
	}
	p.nextToken()
	if p.matchToken(token.Then, token.Do) {
		p.nextToken()
	}
	block := &ast.BlockStatement{Token: p.token}
	block.Statements = []ast.Statement{}
	for !p.matchToken(token.End, token.Else, token.Eof) {
		statement := p.parseStatement()
		if statement != nil {
			block.Statements = append(block.Statements, statement)
		}
		p.nextToken()
	}
	if len(block.Statements) == 0 {
		p.parserError("Empty body in IF")
		return nil
	}
	expression.Then = block
	if p.matchToken(token.Else) {
		elseBody := p.parseBlockBody()
		if len(elseBody.Statements) == 0 {
			p.parserError("Empty ELSE body in IF")
			return nil
		}
		expression.Else = elseBody
	}
	if !p.matchToken(token.End) {
		p.parserError("Missing END closing statement in IF")
		return nil
	}
	return expression
}

func (p *Parser) parseMatch() ast.Expression {
	expression := &ast.Match{Token: p.token}
	p.nextToken()
	expression.Control = p.parseExpression(Lowest)
	if expression.Control != nil {
		p.nextToken()
	}
	if !p.matchToken(token.With, token.NewLine) {
		p.parserError("Missing WITH statement in inline MATCH")
	}
	p.nextToken()
	var whens []*ast.MatchWhen
	for !p.matchToken(token.End, token.Eof) {
		switch p.token.Type {
		case token.When:
			matchWhen := &ast.MatchWhen{Token: p.token}
			list := &ast.ExpressionList{Token: p.token}
			p.nextToken()
			list.Elements = p.parseDelimited(token.Comma, token.NewLine, token.Then)
			if len(list.Elements) == 0 {
				p.parserError("Missing expression in MATCH WHEN")
				break
			}
			matchWhen.Values = list
			matchWhen.Body = p.parseMatchWhen()
			whens = append(whens, matchWhen)
		case token.Else:
			if !p.peekTokenMatch(token.Then, token.NewLine) {
				p.parserError("ELSE when in MATCH can't have parameters")
				return nil
			}
			p.nextToken()
			expression.Else = p.parseMatchWhen()
			if len(expression.Else.Statements) == 0 {
				p.parserError("Missing ELSE when body in MATCH")
				return nil
			}
		}
		p.nextToken()
	}
	expression.Whens = whens
	if !p.matchToken(token.End) {
		p.parserError("Missing END closing statement in MATCH")
		return nil
	}
	return expression
}

func (p *Parser) parseMatchWhen() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.token}
	block.Statements = []ast.Statement{}
	for !p.peekTokenMatch(token.When, token.Else, token.End, token.Eof) {
		p.nextToken()
		statement := p.parseStatement()
		if statement != nil {
			block.Statements = append(block.Statements, statement)
		}
	}
	return block
}

func (p *Parser) parseRepeat() ast.Expression {
	expression := &ast.Repeat{Token: p.token}
	expression.Arguments = &ast.IdentifierList{}
	var arguments []*ast.Identifier
	p.nextToken()
loop:
	for !p.matchToken(token.Do, token.NewLine, token.Eof) {
		switch p.token.Type {
		case token.Comma:
		case token.In:
			p.nextToken()
			expression.Arguments.Elements = arguments
			expression.Enumerable = p.parseExpression(Lowest)
			if expression.Enumerable == nil {
				p.parserError("Missing enumerable in REPEAT loop")
				return nil
			}
			break loop
		default:
			arguments = append(arguments, &ast.Identifier{Token: p.token, Value: p.token.Literal})
		}
		p.nextToken()
	}
	if p.peekTokenMatch(token.Do) {
		p.nextToken()
	}
	expression.Body = p.parseBlockBody()
	if len(expression.Body.Statements) == 0 {
		p.parserError("Empty body in REPEAT loop")
		return nil
	}
	if !p.matchToken(token.End) {
		p.parserError("Missing END closing statement in REPEAT loop")
		return nil
	}
	return expression
}

func (p *Parser) parseFunction() ast.Expression {
	expression := &ast.Function{Token: p.token, Variadic: false}
	expression.Parameters = []*ast.FunctionParameter{}
	p.nextToken()
	for !p.matchToken(token.Do, token.NewLine) {
		switch p.token.Type {
		case token.LeftParenthesis, token.RightParenthesis:
		case token.Comma:
			if expression.Variadic {
				p.parserError("Variadic argument in function should be the last parameter")
				return nil
			}
		case token.Ellipsis:
			if expression.Variadic {
				p.parserError("Function expects only 1 variadic argument")
				return nil
			}
			if !p.peekTokenMatch(token.Identifier) {
				p.parserError("Variadic argument in function expects an identifier")
				return nil
			}
			expression.Variadic = true
		case token.Eof:
			p.parserError("Missing body in function")
			return nil
		case token.Arrow:
			if p.peekTokenMatch(token.Identifier) {
				p.nextToken()
				expression.ReturnType = &ast.Identifier{Token: p.token, Value: p.token.Literal}
			} else {
				p.parserError("Function expecting a return types")
			}
		case token.Identifier:
			var paramType *ast.Identifier
			var defaultValue ast.Expression
			paramName := &ast.Identifier{Token: p.token, Value: p.token.Literal}
			if p.peekTokenMatch(token.Colon) {
				p.nextToken()
				if p.peekTokenMatch(token.Identifier) {
					p.nextToken()
					paramType = &ast.Identifier{Token: p.token, Value: p.token.Literal}
				} else {
					p.parserError(fmt.Sprintf("Function parameter '%s' expecting a types", paramName.Value))
					return nil
				}
			}
			if p.peekTokenMatch(token.Assign) {
				p.nextToken()
				p.nextToken()
				defaultValue = p.parseExpression(Lowest)
			}
			expression.Parameters = append(expression.Parameters, &ast.FunctionParameter{
				Token:   p.token,
				Name:    paramName,
				Type:    paramType,
				Default: defaultValue,
			})
		default:
			p.parserError(fmt.Sprintf("Unexpected token '%s' as function parameter", p.token.Type))
			return nil
		}
		p.nextToken()
	}
	expression.Body = p.parseBlockBody()
	if len(expression.Body.Statements) == 0 {
		p.parserError("Empty body in function")
		return nil
	}
	if !p.matchToken(token.End) {
		p.parserError("Missing END statement in function")
		return nil
	}
	return expression
}

func (p *Parser) parseFunctionCall(function ast.Expression) ast.Expression {
	expression := &ast.FunctionCall{Token: p.token, Function: function}
	list := &ast.ExpressionList{Token: p.token}
	p.nextToken()
	list.Elements = p.parseDelimited(token.Comma, token.RightParenthesis)
	expression.Arguments = list
	return expression
}

func (p *Parser) parseUse() ast.Expression {
	expression := &ast.Use{Token: p.token}
	p.nextToken()
	switch {
	case p.matchToken(token.String, token.Identifier):
		expression.File = &ast.String{Token: p.token, Value: p.token.Literal}
	default:
		p.parserError("USE expects a string or identifier as filename")
		return nil
	}
	return expression
}

func (p *Parser) parseArrayOrDictionary() ast.Expression {
	p.nextToken()
	var list []ast.Expression
	isDictionary := false
	for !p.matchToken(token.RightBracket) {
		switch {
		case p.matchToken(token.NewLine, token.Eof):
			p.parserError("Missing closing ']' in enumerable")
			return nil
		case p.matchToken(token.FatArrow):
			isDictionary = true
		case p.matchToken(token.Comma):
		default:
			expression := p.parseExpression(Lowest)
			if expression == nil {
				return nil
			}
			list = append(list, expression)
		}
		p.nextToken()
	}
	if !isDictionary {
		expression := &ast.Array{Token: p.token}
		expression.List = &ast.ExpressionList{Elements: list}
		return expression
	}
	if len(list)%2 == 1 {
		p.parserError("Dictionary expects elements as Key:Value")
		return nil
	}
	expression := &ast.Dictionary{Token: p.token}
	expression.Pairs = map[ast.Expression]ast.Expression{}
	for i, v := range list {
		if i%2 == 0 {
			if v == nil || list[i+1] == nil {
				return nil
			}
			expression.Pairs[v] = list[i+1]
		}
	}
	return expression
}

func (p *Parser) parseReturn() *ast.Return {
	statement := &ast.Return{Token: p.token}
	p.nextToken()
	statement.Value = p.parseExpression(Lowest)
	return statement
}

func (p *Parser) parseBreak() ast.Statement {
	return &ast.Break{Token: p.token}
}

func (p *Parser) parseContinue() ast.Statement {
	return &ast.Continue{Token: p.token}
}

func (p *Parser) parseSubscript(left ast.Expression) ast.Expression {
	expression := &ast.Subscript{Token: p.token, Left: left}
	p.nextToken()
	if p.matchToken(token.RightBracket) {
		expression.Index = &ast.PlaceHolder{Token: p.token}
		return expression
	}
	if p.matchToken(token.Underscore) {
		p.nextToken()
		expression.Index = &ast.PlaceHolder{Token: p.token}
		return expression
	}
	expression.Index = p.parseExpression(Lowest)
	if !p.peekTokenMatch(token.RightBracket) {
		p.parserError("Missing closing ] in subscript expression")
		return nil
	}
	p.nextToken()
	return expression
}

func (p *Parser) parsePipe(left ast.Expression) ast.Expression {
	expression := &ast.Pipe{
		Token: p.token,
		Left:  left,
	}
	p.nextToken()
	expression.Right = p.parseExpression(Pipe)
	return expression
}

func (p *Parser) parseArrowFunction(left ast.Expression) ast.Expression {
	expression := &ast.Function{Token: p.token}
	expression.Parameters = []*ast.FunctionParameter{}
	switch expressionType := left.(type) {
	case *ast.Identifier:
		expression.Parameters = append(expression.Parameters, &ast.FunctionParameter{
			Token: p.token,
			Name:  expressionType,
		})
	case *ast.ExpressionList:
		for _, v := range expressionType.Elements {
			switch parameter := v.(type) {
			case *ast.Identifier:
				expression.Parameters = append(expression.Parameters, &ast.FunctionParameter{
					Token: p.token,
					Name:  parameter,
				})
			default:
				p.parserError("Arrow function expects a list of identifiers as arguments")
				return nil
			}
		}
	default:
		p.parserError("Arrow function expects identifiers as arguments")
		return nil
	}
	p.nextToken()
	expression.Body = &ast.BlockStatement{
		Statements: []ast.Statement{
			p.parseExpressionStatement(),
		},
	}
	return expression
}

func (p *Parser) parseAssign(left ast.Expression) ast.Expression {
	expression := &ast.Assign{
		Token:    p.token,
		Operator: p.token.Literal,
	}
	switch identifier := left.(type) {
	case *ast.Identifier:
		expression.Name = identifier
	case *ast.Subscript:
		switch identifier.Left.(type) {
		case *ast.Identifier:
			expression.Name = identifier
		default:
			p.parserError("Assignment operator expects an identifier")
			return nil
		}
	default:
		p.parserError("Assignment operator expects an identifier")
		return nil
	}
	p.nextToken()
	expression.Right = p.parseExpression(Lowest)
	if expression.Right == nil {
		return nil
	}
	switch expression.Operator {
	case token.PlusAssign, token.MinusAssign, token.MultiplyAssign, token.DivideAssign:
		expression.Right = &ast.InfixExpression{
			Token:    p.token,
			Left:     left,
			Right:    expression.Right,
			Operator: string(expression.Operator[0]),
		}
	}
	return expression
}

func (p *Parser) parseTernary(left ast.Expression) ast.Expression {
	expression := &ast.If{Token: p.token}
	expression.Condition = left
	p.nextToken()
	then := p.parseExpression(Lowest)
	if then == nil {
		p.parserError("Missing THEN condition in ternary operator")
		return nil
	}
	expression.Then = &ast.BlockStatement{
		Statements: []ast.Statement{
			&ast.ExpressionStatement{Expression: then},
		},
	}
	p.nextToken()
	if !p.matchToken(token.Colon) {
		p.parserError("Ternary operator expects an else (:) expression")
		return nil
	}
	p.nextToken()
	elseExpr := p.parseExpression(Lowest)
	if elseExpr == nil {
		p.parserError("Missing ELSE condition in ternary operator")
		return nil
	}
	expression.Else = &ast.BlockStatement{
		Statements: []ast.Statement{
			&ast.ExpressionStatement{Expression: elseExpr},
		},
	}
	return expression
}

func (p *Parser) parseIs(left ast.Expression) ast.Expression {
	expression := &ast.Is{
		Token: p.token,
		Left:  left,
	}
	p.nextToken()
	if !p.matchToken(token.Identifier) {
		p.parserError("IS operator expects a types")
		return nil
	}
	expression.Right = &ast.Identifier{Token: p.token, Value: p.token.Literal}
	return expression
}

func (p *Parser) parseAs(left ast.Expression) ast.Expression {
	expression := &ast.As{
		Token: p.token,
		Left:  left,
	}
	p.nextToken()
	if !p.matchToken(token.Identifier) {
		p.parserError("AS operator expects a types")
		return nil
	}
	expression.Right = &ast.Identifier{Token: p.token, Value: p.token.Literal}
	return expression
}

func (p *Parser) parsePlaceHolder() ast.Expression {
	return &ast.PlaceHolder{Token: p.token}
}

func (p *Parser) parseSymbol() ast.Expression {
	p.nextToken()
	expression := &ast.Symbol{Token: p.token}
	if !p.matchToken(token.Identifier) {
		p.parserError("Symbol expects an identifier")
		return nil
	}
	expression.Value = p.token.Literal
	return expression
}

func (p *Parser) parseDelimited(delimiter token.TType, end ...token.TType) []ast.Expression {
	var list []ast.Expression
	for !p.matchToken(end...) {
		switch p.token.Type {
		case delimiter:
		case token.NewLine, token.Eof:
			p.parserError(fmt.Sprintf("Missing closing '%s' in parameter list", end))
			return list
		default:
			elem := p.parseExpression(Lowest)
			if elem == nil {
				p.parserError(fmt.Sprintf("Unexpected '%s' in expression list", p.token.Literal))
				return list
			}
			list = append(list, elem)
		}
		p.nextToken()
	}
	return list
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	statement := &ast.ExpressionStatement{Token: p.token}
	statement.Expression = p.parseExpression(Lowest)
	return statement
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixFunction[p.token.Type]
	if prefix == nil {
		p.parserError(fmt.Sprintf("Unexpected expression '%s'", p.token.Literal))
		return nil
	}
	left := prefix()
	for precedence < p.peekPrecedence() {
		infix := p.infixFunction[p.peekToken.Type]
		if infix == nil {
			return left
		}
		p.nextToken()
		left = infix(left)
	}
	return left
}

func (p *Parser) parseGroup() ast.Expression {
	p.nextToken()
	expression := p.parseExpression(Lowest)
	if p.peekTokenMatch(token.Comma) {
		p.nextToken()
		list := &ast.ExpressionList{}
		list.Elements = []ast.Expression{expression}
		rest := p.parseDelimited(token.Comma, token.RightParenthesis)
		if rest != nil {
			list.Elements = append(list.Elements, rest...)
		}
		return list
	}
	if !p.peekTokenMatch(token.RightParenthesis) {
		p.parserError("Missing closing ')' for grouped expression")
		return nil
	}
	p.nextToken()
	return expression
}

func (p *Parser) parsePrefix() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.token,
		Operator: p.token.Literal,
	}
	p.nextToken()
	expression.Right = p.parseExpression(Prefix)
	return expression
}

func (p *Parser) parseInfix(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.token,
		Operator: p.token.Literal,
		Left:     left,
	}
	precedence := p.precedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}

func (p *Parser) parseInfixRight(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.token,
		Operator: p.token.Literal,
		Left:     left,
	}
	precedence := p.precedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence - 1)
	return expression
}

func (p *Parser) parseBlockBody() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.token}
	block.Statements = []ast.Statement{}
	p.nextToken()
	for !p.matchToken(token.End, token.Eof) {
		if statement := p.parseStatement(); statement != nil {
			block.Statements = append(block.Statements, statement)
		}
		p.nextToken()
	}
	return block
}

func (p *Parser) synchronize() {
	for !p.matchToken(token.Eof) {
		switch p.token.Type {
		case token.Val, token.If, token.Repeat, token.Match, token.When,
			token.Else, token.Return, token.Function, token.Module:
			return
		}
		p.nextToken()
	}
}

func (p *Parser) parserError(msg string) {
	rerror.Error(rerror.Parse, p.token.Position, msg)
	p.synchronize()
}
