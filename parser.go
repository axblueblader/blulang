package main

import (
	"log"
	"strconv"
)

type Parser struct {
	tokens []Token
}

func NewParser() Parser {
	return Parser{
		tokens: make([]Token, 0),
	}
}
func (p *Parser) CreateAST(source string) Program {
	tokens := Tokenize(source)
	p.tokens = tokens
	program := NewProgram()

	for len(p.tokens) > 0 {
		program.body = append(program.body, p.parseStatement())
	}
	return program
}

func (p *Parser) peek() Token {
	if len(p.tokens) > 0 {
		return p.tokens[0]
	}
	return Token{}
}
func (p *Parser) peekNext() Token {
	if len(p.tokens) > 1 {
		return p.tokens[1]
	}
	return Token{}
}

func (p *Parser) pop() {
	p.tokens = p.tokens[1:]
}

// Order of precedence;
// Variable declaration
// Conditional
// Loop
// Assignment
// Logical
// Comparison
// Additive
// Multiplicative
// Variables and function call
// Literals

func (p *Parser) parseStatement() Statement {
	return p.parseExpression()
}

func (p *Parser) parseExpression() Expression {
	if p.peek().name == TkDeclareVar {
		return p.parseVariableDeclarationExpression()
	}
	if p.peek().name == TkDeclareFunc {
		return p.parseFunctionDeclarationExpression()
	}
	if p.peek().name == TkIf {
		return p.parseConditionalExpression()
	}
	if p.peek().name == TkWhile {
		return p.parseWhileLoopExpression()
	}
	return p.parseAssignmentExpression()
}

func (p *Parser) parseWhileLoopExpression() Expression {
	p.pop() // pop 'if'
	conditionExpr := p.parseLogicalExpression()

	var statements []Statement
	statements = p.parseCodeBlock(statements)
	return NewWhileLoopExpression(conditionExpr, statements)
}

func (p *Parser) parseConditionalExpression() Expression {
	p.pop() // pop 'if'
	conditionExpr := p.parseLogicalExpression()

	var trueBodyStatements []Statement
	trueBodyStatements = p.parseCodeBlock(trueBodyStatements)

	var falseBodyStatements []Statement
	if p.peek().name == TkElse {
		p.pop() // pop 'else'
		if p.peek().name == TkIf {
			falseBodyStatements = append(falseBodyStatements, p.parseConditionalExpression())
		} else {
			falseBodyStatements = p.parseCodeBlock(falseBodyStatements)
		}
	}

	return NewConditionalExpression(conditionExpr, trueBodyStatements, falseBodyStatements)
}

func (p *Parser) parseCodeBlock(statements []Statement) []Statement {
	// pop open curly bracket
	p.pop()
	for p.peek().name != TkCloseCurly {
		statements = append(statements, p.parseStatement())
	}
	// pop close curly bracket
	p.pop()
	return statements
}

func (p *Parser) parseAssignmentExpression() Expression {
	expr := p.parseLogicalExpression()
	if p.peek().value == "=" {
		p.pop() // pop equal sign
		return NewBinaryExpression(expr, p.parseExpression(), "=")
	}
	return expr
}

func (p *Parser) parseVariableDeclarationExpression() Expression {
	// pop the declaration keyword
	p.pop()

	variableName := p.peek().value
	p.pop()

	// pop the equal sign
	p.pop()

	value := p.parseExpression()
	return NewVarDeclareExpression(variableName, value)
}

func (p *Parser) parseFunctionDeclarationExpression() Expression {
	// pop the declaration keyword
	p.pop()

	functionName := ""
	if p.peek().name == TkIdentifier {
		functionName = p.peek().value
		p.pop()
	}

	// pop open round bracket
	p.pop()
	var arguments []Identifier
	for p.peek().name != TkCloseRound {
		arguments = append(arguments, NewIdentifier(p.peek().value))
		p.pop()
		if p.peek().name == TkComma {
			p.pop()
		}
	}
	// pop close round bracket
	p.pop()

	var statements []Statement
	statements = p.parseCodeBlock(statements)

	return NewFuncDeclareExpression(functionName, arguments, statements)
}

func (p *Parser) parseIdentifierOrFunctionCallExpression() Expression {
	// parse function call
	if p.peekNext().name == TkOpenRound {
		funcName := p.peek().value
		// pop name
		p.pop()
		// pop open round bracket
		p.pop()
		var args []Expression
		for p.peek().name != TkCloseRound {
			args = append(args, p.parseExpression())
			if p.peek().name == TkComma {
				p.pop()
			}
		}
		// pop close round bracket
		p.pop()
		return NewFuncCallExpression(funcName, args)
	}
	// parse array index access
	if p.peekNext().name == TKOpenSquare {
		identifierName := p.peek().value
		p.pop() // pop name
		p.pop() // pop [
		indexExpr := p.parseExpression()
		p.pop() // pop ]
		return NewArrayAccessExpr(identifierName, indexExpr)
	}
	// parse property access
	if p.peekNext().name == TkDot {
	}
	identifierName := p.peek().value
	p.pop()
	return NewIdentifier(identifierName)
}

func (p *Parser) parseLogicalExpression() Expression {
	leftExp := p.parseComparisonExpression()
	operator := p.peek().value
	for operator == "&&" || operator == "||" {
		p.pop()
		rightExp := p.parseExpression()
		leftExp = NewBinaryExpression(leftExp, rightExp, operator)
		operator = p.peek().value
	}
	return leftExp
}

func (p *Parser) parseComparisonExpression() Expression {
	leftExp := p.parseAdditiveExpression()
	operator := p.peek().value
	for operator == "==" || operator == "!=" || operator == "<" || operator == ">" || operator == "<=" || operator == ">=" {
		p.pop()
		rightExp := p.parseExpression()
		leftExp = NewBinaryExpression(leftExp, rightExp, operator)
		operator = p.peek().value
	}
	return leftExp
}

func (p *Parser) parseAdditiveExpression() Expression {
	leftExp := p.parseMultiplicativeExpression()
	operator := p.peek().value
	for operator == "+" || operator == "-" {
		p.pop()
		rightExp := p.parseMultiplicativeExpression()
		leftExp = NewBinaryExpression(leftExp, rightExp, operator)
		operator = p.peek().value
	}

	return leftExp
}

func (p *Parser) parseMultiplicativeExpression() Expression {
	leftExp := p.parsePrimaryExpression()
	operator := p.peek().value
	for operator == "*" || operator == "/" {
		p.pop()
		rightExp := p.parseMultiplicativeExpression()
		leftExp = NewBinaryExpression(leftExp, rightExp, operator)
		operator = p.peek().value
	}

	return leftExp
}

func (p *Parser) parseArrayExpression() Expression {
	var values []Expression
	p.pop() // pop [
	for p.peek().name != TkCloseSquare {
		values = append(values, p.parseExpression())
		if p.peek().name == TkComma {
			p.pop() // pop comma after an element
		}
	}
	p.pop() // pop ]
	return NewArrayLiteral(values)
}

func (p *Parser) parseGroupedExpression() Expression {
	p.pop()
	expr := p.parseLogicalExpression()
	p.pop()
	return expr
}

func (p *Parser) parseNotExpression() Expression {
	p.pop() // pop !
	expression := p.parseExpression()
	return NewBinaryExpression(expression, NewIdentifier("true"), "!=")
}

func (p *Parser) parsePrimaryExpression() Expression {
	token := p.peek()
	switch token.name {
	case TkNumber:
		p.pop()
		intVal, _ := strconv.Atoi(token.value)
		return NewIntLiteral(intVal)
	case TkString:
		p.pop()
		return NewStringLiteral(token.value)
	case TkBreak:
		p.pop()
		return NewBreakStatement()
	case TkReturn:
		p.pop()
		return NewReturnStatement()
	case TkNot:
		return p.parseNotExpression()
	case TkIdentifier:
		return p.parseIdentifierOrFunctionCallExpression()
	case TKOpenSquare:
		return p.parseArrayExpression()
	case TkOpenRound:
		return p.parseGroupedExpression()
	default:
		log.Panicf("Unknown token %v", token)
	}
	return NullLiteral{}
}
