package main

import (
	"log"
)

func Eval(statement Statement, scope *Scope) RuntimeVal {
	switch statement.Kind() {
	case StmtProgram:
		return EvalProgram(statement.(Program), scope)
	case StmtBinaryExpr:
		return EvalBinaryExpression(statement.(BinaryExpression), scope)
	case StmtIntLiteral:
		return NewIntVal(statement.(IntLiteral).value)
	case StmtStringLiteral:
		return NewStringVal(statement.(StringLiteral).value)
	case StmtArrayLiteral:
		return EvalArrayLiteral(statement.(ArrayLiteral), scope)
	case StmtVarDeclareExpr:
		return EvalVarDeclareExpression(statement.(VarDeclareExpression), scope)
	case StmtFuncDeclareExpr:
		return EvalFuncDeclareExpression(statement.(FuncDeclareExpression), scope)
	case StmtFuncCallExpr:
		return EvalFuncCallExpression(statement.(FuncCallExpression), scope)
	case StmtArrayAccessExpr:
		return EvalArrayAccessExpression(statement.(ArrayAccessExpr), scope)
	case StmtIdentifier:
		return EvalIdentifier(statement.(Identifier), scope)
	case StmtConditionalExpr:
		return EvalConditionalExpression(statement.(ConditionalExpression), scope)
	case StmtWhileLoopExpr:
		return EvalWhileLoopExpression(statement.(WhileLoopExpression), scope)
	case StmtNullLiteral:
		return NullVal{}
	}
	log.Panicf("Invalid statement: %v", statement)
	return nil
}

func EvalArrayAccessExpression(expr ArrayAccessExpr, scope *Scope) RuntimeVal {
	index := Eval(expr.index, scope).(IntVal)
	arrayVal := EvalIdentifier(NewIdentifier(expr.name), scope).(ArrayVal)
	return arrayVal.values[index.value]
}

func EvalArrayLiteral(statement ArrayLiteral, scope *Scope) RuntimeVal {
	var runTimeValues []RuntimeVal
	for _, expr := range statement.values {
		runTimeValues = append(runTimeValues, Eval(expr, scope))
	}
	return NewArrayVal(runTimeValues)
}

func EvalProgram(program Program, scope *Scope) RuntimeVal {
	var lastValue RuntimeVal = NullVal{}
	for _, statement := range program.body {
		lastValue = Eval(statement, scope)
	}
	return lastValue
}

func EvalConditionalBody(body []Statement, scope *Scope) RuntimeVal {
	var lastValue RuntimeVal = NullVal{}
	for _, statement := range body {
		if statement.Kind() == StmtBreak {
			return NewBreakVal(lastValue)
		}
		if statement.Kind() == StmtReturn {
			return NewReturnVal(lastValue)
		}
		lastValue = Eval(statement, scope)
	}
	return lastValue
}
func EvalConditionalExpression(conditionStatement ConditionalExpression, scope *Scope) RuntimeVal {
	conditionResult := Eval(conditionStatement.condition, scope)
	bodyScope := NewScope(scope)
	if conditionResult.Value() == true {
		return EvalConditionalBody(conditionStatement.trueBody, bodyScope)
	} else {
		return EvalConditionalBody(conditionStatement.falseBody, bodyScope)
	}
}

func EvalWhileLoopExpression(expression WhileLoopExpression, scope *Scope) RuntimeVal {
	conditionResult := Eval(expression.condition, scope)
	bodyScope := NewScope(scope)
	var lastValue RuntimeVal = NullVal{}
	for conditionResult.Value() == true {
		for _, statement := range expression.body {
			lastValue = Eval(statement, bodyScope)
			if lastValue.Kind() == VaBreakVal {
				// return the actual value
				return lastValue.(BreakVal).lastValue
			}
			if lastValue.Kind() == VaReturnVal {
				// return the flow control val
				return lastValue
			}
		}
		conditionResult = Eval(expression.condition, scope)
	}
	return lastValue
}

func EvalFuncCallExpression(call FuncCallExpression, scope *Scope) RuntimeVal {
	funcVal := scope.GetVarVal(call.name)
	switch funcVal.(type) {
	case FunctionVal:
		return EvalUserFuncCallExpression(funcVal.(FunctionVal), call.arguments, scope)
	case NativeFuncVal:
		return EvalNativeFuncCallExpression(funcVal.(NativeFuncVal), call.arguments, scope)
	}
	return NullVal{}
}

func EvalNativeFuncCallExpression(funcVal NativeFuncVal, argExpressions []Expression, scope *Scope) RuntimeVal {
	var argsVal []RuntimeVal
	for _, expr := range argExpressions {
		argsVal = append(argsVal, Eval(expr, scope))
	}
	return funcVal.Invoke(scope, argsVal...)
}

func EvalUserFuncCallExpression(functionVal FunctionVal, argExpressions []Expression, scope *Scope) RuntimeVal {
	funcScope := NewScope(scope)
	var lastValue RuntimeVal = NullVal{}
	for i, identifier := range functionVal.arguments {
		funcScope.DeclareVar(identifier.name, Eval(argExpressions[i], funcScope))
	}
	for _, statement := range functionVal.body {
		lastValue = Eval(statement, funcScope)
		if lastValue.Kind() == VaReturnVal {
			// return the actual value
			return lastValue.(ReturnVal).lastValue
		}
	}
	return lastValue
}

func EvalBinaryExpression(binaryExp BinaryExpression, scope *Scope) RuntimeVal {
	lhs := Eval(binaryExp.left, scope)
	rhs := Eval(binaryExp.right, scope)
	operator := binaryExp.operator
	// assignment
	if operator == "=" {
		return EvalAssignmentExpression(binaryExp.left, rhs, scope)
	}
	if operator == "||" || operator == "&&" {
		return EvalLogicalBinaryExpression(lhs, rhs, operator)
	}
	// comparison operator
	if operator == "==" || operator == "!=" || operator == "<" || operator == ">" || operator == "<=" || operator == ">=" {
		return EvalComparisonBinaryExpression(lhs, rhs, operator)
	}
	// normal math operator
	if rhs.Kind() == VaIntVal {
		if lhs.Kind() == VaIntVal {
			return EvalIntBinaryExpression(lhs.(IntVal), rhs.(IntVal), operator)
		} else if operator == "-" {
			return EvalIntBinaryExpression(NewIntVal(0), rhs.(IntVal), operator)
		}
	}

	// array math operator
	if lhs.Kind() == rhs.Kind() && rhs.Kind() == VaArrayVal {
		return EvalArrayBinaryExpression(lhs.(ArrayVal), rhs.(ArrayVal), operator)
	}
	return NullVal{}
}

func EvalLogicalBinaryExpression(lhs RuntimeVal, rhs RuntimeVal, operator string) RuntimeVal {
	switch operator {
	case "||":
		if lhs.Value() == true || rhs.Value() == true {
			return TrueVal
		} else {
			return FalseVal
		}
	case "&&":
		if lhs.Value() == true && rhs.Value() == true {
			return TrueVal
		} else {
			return FalseVal
		}
	}
	return NullVal{}
}

func EvalArrayBinaryExpression(lhs ArrayVal, rhs ArrayVal, operator string) RuntimeVal {
	switch operator {
	case "+":
		return NewArrayVal(append(lhs.values, rhs.values...))
	}
	log.Panic("Unsupported operator for array: ", operator)
	return NullVal{}
}

func EvalComparisonBinaryExpression(lhs RuntimeVal, rhs RuntimeVal, operator string) RuntimeVal {
	if operator == "==" {
		if lhs.Kind() == rhs.Kind() {
			return NewBoolVal(lhs.Value() == rhs.Value())
		} else {
			return FalseVal
		}
	}
	if operator == "!=" {
		if lhs.Kind() == rhs.Kind() {
			return NewBoolVal(lhs.Value() != rhs.Value())
		} else {
			return TrueVal
		}
	}
	if lhs.Kind() == rhs.Kind() && lhs.Kind() == VaIntVal {
		return EvalIntComparisonExpression(lhs.(IntVal), rhs.(IntVal), operator)
	}
	return NullVal{}
}

func EvalIntComparisonExpression(lhs IntVal, rhs IntVal, operator string) RuntimeVal {
	lhsVal := lhs.value
	rhsVal := rhs.value
	switch operator {
	case "<=":
		return NewBoolVal(lhsVal <= rhsVal)
	case ">=":
		return NewBoolVal(lhsVal >= rhsVal)
	case "<":
		return NewBoolVal(lhsVal < rhsVal)
	case ">":
		return NewBoolVal(lhsVal > rhsVal)
	}
	log.Panic("Unsupported operator: ", operator)
	return NullVal{}
}

func EvalIdentifier(identifier Identifier, scope *Scope) RuntimeVal {
	return scope.GetVarVal(identifier.name)
}

func EvalVarDeclareExpression(varDeclareExpr VarDeclareExpression, scope *Scope) RuntimeVal {
	varName := varDeclareExpr.name
	varValue := Eval(varDeclareExpr.valueExpr, scope)
	// create variable in scope
	scope.DeclareVar(varName, varValue)
	return varValue
}

func EvalFuncDeclareExpression(funcDeclareExpr FuncDeclareExpression, scope *Scope) RuntimeVal {
	funcName := funcDeclareExpr.name
	funcVal := NewFuncVal(funcName, funcDeclareExpr.arguments, funcDeclareExpr.body)
	scope.DeclareVar(funcName, funcVal)
	return funcVal
}

func EvalAssignmentExpression(expr Expression, varValue RuntimeVal, scope *Scope) RuntimeVal {
	switch expr.Kind() {
	case StmtIdentifier:
		// assign variable in scope
		scope.AssignVar(expr.(Identifier).name, varValue)
		return varValue
	case StmtArrayAccessExpr:
		arrayAccessExpr := expr.(ArrayAccessExpr)
		index := Eval(arrayAccessExpr.index, scope).(IntVal)
		arrayVal := EvalIdentifier(NewIdentifier(arrayAccessExpr.name), scope).(ArrayVal)
		arrayVal.values[index.value] = varValue
		return varValue
	}

	return varValue
}

func EvalIntBinaryExpression(val IntVal, val2 IntVal, operator string) RuntimeVal {
	switch operator {
	case "+":
		return IntVal{
			value: val.value + val2.value,
		}
	case "-":
		return IntVal{
			value: val.value - val2.value,
		}
	case "*":
		return IntVal{
			value: val.value * val2.value,
		}
	case "/":
		return IntVal{
			value: val.value / val2.value,
		}
	}
	return NullVal{}
}
