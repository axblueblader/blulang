package main

type StmtType string

const (
	StmtBinaryExpr      StmtType = "BinaryExpr"
	StmtProgram         StmtType = "Program"
	StmtIntLiteral      StmtType = "IntLiteral"
	StmtStringLiteral   StmtType = "StringLiteral"
	StmtNullLiteral     StmtType = "NullLiteral"
	StmtVarDeclareExpr  StmtType = "VarDeclareExpr"
	StmtFuncDeclareExpr StmtType = "FuncDeclareExpr"
	StmtFuncCallExpr    StmtType = "FuncCallExpr"
	StmtIdentifier      StmtType = "Identifier"
	StmtConditionalExpr StmtType = "ConditionalExpr"
	StmtWhileLoopExpr   StmtType = "WhileLoopExpr"
	StmtBreak           StmtType = "BreakStmt"
	StmtReturn          StmtType = "ReturnStmt"
	StmtArrayLiteral    StmtType = "ArrayLiteral"
	StmtArrayAccessExpr StmtType = "ArrayAccessExpr"
	StmtObjDeclareExpr  StmtType = "ObjectDeclareExpr"
	StmtObjAccessExpr   StmtType = "ObjectAccessExpr"
)

type Statement interface {
	Kind() StmtType
}

type Expression interface {
	Statement
}

type Program struct {
	body []Statement
}

func NewProgram() Program {
	return Program{
		body: make([]Statement, 0),
	}
}

func (p Program) Kind() StmtType {
	return StmtProgram
}

type WhileLoopExpression struct {
	condition Expression
	body      []Statement
}

func (e WhileLoopExpression) Kind() StmtType {
	return StmtWhileLoopExpr
}

func NewWhileLoopExpression(condition Expression, body []Statement) WhileLoopExpression {
	return WhileLoopExpression{
		condition: condition,
		body:      body,
	}
}

type BreakStatement struct{}

func (s BreakStatement) Kind() StmtType { return StmtBreak }

func NewBreakStatement() BreakStatement { return BreakStatement{} }

type ReturnStatement struct{}

func (s ReturnStatement) Kind() StmtType { return StmtReturn }

func NewReturnStatement() ReturnStatement { return ReturnStatement{} }

type ConditionalExpression struct {
	condition Expression
	trueBody  []Statement
	falseBody []Statement
}

func (e ConditionalExpression) Kind() StmtType {
	return StmtConditionalExpr
}

func NewConditionalExpression(condition Expression, trueBody []Statement, falseBody []Statement) ConditionalExpression {
	return ConditionalExpression{
		condition: condition,
		trueBody:  trueBody,
		falseBody: falseBody,
	}
}

type BinaryExpression struct {
	left     Expression
	right    Expression
	operator string
}

func NewBinaryExpression(left Expression, right Expression, operator string) BinaryExpression {
	return BinaryExpression{
		left:     left,
		right:    right,
		operator: operator,
	}
}

func (b BinaryExpression) Kind() StmtType {
	return StmtBinaryExpr
}

type VarDeclareExpression struct {
	name      string
	valueExpr Expression
}

func (v VarDeclareExpression) Kind() StmtType {
	return StmtVarDeclareExpr
}

func NewVarDeclareExpression(name string, value Expression) VarDeclareExpression {
	return VarDeclareExpression{
		name:      name,
		valueExpr: value,
	}
}

type FuncDeclareExpression struct {
	name      string
	arguments []Identifier
	body      []Statement
}

func (v FuncDeclareExpression) Kind() StmtType {
	return StmtFuncDeclareExpr
}

func NewFuncDeclareExpression(name string, arguments []Identifier, body []Statement) FuncDeclareExpression {
	return FuncDeclareExpression{
		name:      name,
		arguments: arguments,
		body:      body,
	}
}

type FuncCallExpression struct {
	name      string
	arguments []Expression
}

func (v FuncCallExpression) Kind() StmtType {
	return StmtFuncCallExpr
}

func NewFuncCallExpression(name string, arguments []Expression) FuncCallExpression {
	return FuncCallExpression{
		name:      name,
		arguments: arguments,
	}
}

type IntLiteral struct {
	value int
}

func (l IntLiteral) Kind() StmtType {
	return StmtIntLiteral
}

func NewIntLiteral(value int) IntLiteral {
	return IntLiteral{
		value: value,
	}
}

type StringLiteral struct {
	value string
}

func (l StringLiteral) Kind() StmtType {
	return StmtStringLiteral
}

func NewStringLiteral(value string) StringLiteral {
	return StringLiteral{value: value}
}

type NullLiteral struct{}

func (n NullLiteral) Kind() StmtType {
	return StmtNullLiteral
}

type Identifier struct {
	name string
}

func (i Identifier) Kind() StmtType {
	return StmtIdentifier
}

func NewIdentifier(name string) Identifier {
	return Identifier{name: name}
}

type ArrayLiteral struct {
	values []Expression
}

func (a ArrayLiteral) Kind() StmtType {
	return StmtArrayLiteral
}

func NewArrayLiteral(values []Expression) ArrayLiteral {
	return ArrayLiteral{values: values}
}

type ArrayAccessExpr struct {
	name  string
	index Expression
}

func (a ArrayAccessExpr) Kind() StmtType {
	return StmtArrayAccessExpr
}

func NewArrayAccessExpr(name string, index Expression) ArrayAccessExpr {
	return ArrayAccessExpr{name: name, index: index}
}

type ObjectDeclareExpr struct {
	props map[string]Expression
}

func (e ObjectDeclareExpr) Kind() StmtType {
	return StmtObjDeclareExpr
}

func NewObjectDeclareExpr(props map[string]Expression) ObjectDeclareExpr {
	return ObjectDeclareExpr{props: props}
}

type ObjectAccessExpr struct {
	owner    Identifier
	property Expression
}

func (e ObjectAccessExpr) Kind() StmtType { return StmtObjAccessExpr }

func NewObjectAccessExpr(owner Identifier, property Expression) ObjectAccessExpr {
	return ObjectAccessExpr{owner: owner, property: property}
}
