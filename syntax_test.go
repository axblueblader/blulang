package main_test

import (
	"blulang"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestMath(t *testing.T) {
	scope := main.NewScope(nil)
	parser := main.NewParser()
	code := `
	((10+4) * 2 - 3) / ((9-7)*(3-2))
	`
	program := parser.CreateAST(code)
	result := main.Eval(program, scope)
	assert.Equal(t, 12, result.Value())
	assert.Equal(t, main.VaIntVal, result.Kind())
}

func TestVariableDeclaration(t *testing.T) {
	scope := main.NewScope(nil)
	parser := main.NewParser()
	code := `
	let a = 10
	a
	`
	program := parser.CreateAST(code)
	result := main.Eval(program, scope)
	assert.Equal(t, 10, result.Value())
	assert.Equal(t, main.VaIntVal, result.Kind())
}

func TestVariableAssignment(t *testing.T) {
	scope := main.NewScope(nil)
	parser := main.NewParser()
	code := `
	let a = 10
	let b1 = 20
	a = a*b1
	a
	`
	program := parser.CreateAST(code)
	result := main.Eval(program, scope)
	assert.Equal(t, 200, result.Value())
	assert.Equal(t, main.VaIntVal, result.Kind())
}

func TestComparisonExpression(t *testing.T) {
	scope := main.NewGlobalScope()
	parser := main.NewParser()
	sources := []string{
		"4*3 == 2*6",
		"1+1 == abs(2-4)*1",
		"3 < 4",
		"3 <= 4",
		"4 <= 4",
		"4 > 3",
		"4 >= 3",
		"3 >= 3",
		"1 != 2",
	}
	for _, code := range sources {
		program := parser.CreateAST(code)
		result := main.Eval(program, scope)
		assert.Equal(t, true, result.Value())
		assert.Equal(t, main.VaBoolVal, result.Kind())
	}
}

func TestLogicalExpression(t *testing.T) {
	scope := main.NewGlobalScope()
	parser := main.NewParser()
	sources := []string{
		"(4 < 3) || (3 < 4)",
		"(1 == 1) && (1 != 2)",
		"true",
		"true == true",
		"!false",
		"!(1 != 1)",
	}
	for _, code := range sources {
		program := parser.CreateAST(code)
		result := main.Eval(program, scope)
		assert.Equalf(t, true, result.Value(), "%v", program)
		assert.Equal(t, main.VaBoolVal, result.Kind())
	}
	falseSources := []string{
		"(4 < 3) || (5 < 4)",
		"(1 == 1) && (2 != 2)",
		"false",
		"false != false",
		"!true",
		"!(1 == 1)",
	}
	for _, code := range falseSources {
		program := parser.CreateAST(code)
		result := main.Eval(program, scope)
		assert.Equalf(t, false, result.Value(), "%v", program)
		assert.Equal(t, main.VaBoolVal, result.Kind())
	}
}

func TestFunctionDeclareAndCall(t *testing.T) {
	scope := main.NewScope(nil)
	parser := main.NewParser()
	code := `
	fn sum (a,b) {
		a = a*a
		a+b
	}
	let a = 10
	let b = 1
	sum(a,b)
	`
	program := parser.CreateAST(code)
	result := main.Eval(program, scope)
	assert.Equal(t, 101, result.Value())
	assert.Equal(t, main.VaIntVal, result.Kind())
}

func TestAnonFunctionDeclareAndCall(t *testing.T) {
	scope := main.NewScope(nil)
	parser := main.NewParser()
	code := `
	let sum = fn (a,b) {
		a = a*a
		a+b
	}
	let a = 10
	let b = 1
	sum(a,b)
	`
	program := parser.CreateAST(code)
	result := main.Eval(program, scope)
	assert.Equal(t, 101, result.Value())
	assert.Equal(t, main.VaIntVal, result.Kind())
}

func TestFunctionCallReturn(t *testing.T) {
	scope := main.NewScope(nil)
	parser := main.NewParser()
	code := `
	fn increase(a) {
		while a < 100 {
			a = a + 1
			if a == 50 {
				a
				return
			}
		}
	}
	let a = 1
	increase(a)
	`
	program := parser.CreateAST(code)
	result := main.Eval(program, scope)
	assert.Equal(t, 50, result.Value())
	assert.Equal(t, main.VaIntVal, result.Kind())
}

func TestConditionalStatement(t *testing.T) {
	scope := main.NewGlobalScope()
	parser := main.NewParser()
	code := `
	let a = 1000
	if 0 == 1 {
		a = 1
	} else if 1 == 0 {
		a = 0
	} else {
		a = 2
	}

	let b = if 1-1 == 2-2 {
		b = 10
	} else {
		b = 20
	}
	
	fn getC() {
		if true {
			100
		} else {
			200
		}
	}
	let c = getC()
	c + b + a
	`
	program := parser.CreateAST(code)
	result := main.Eval(program, scope)
	assert.Equal(t, 112, result.Value())
	assert.Equal(t, main.VaIntVal, result.Kind())
}

func TestWhileLoop(t *testing.T) {
	scope := main.NewScope(nil)
	parser := main.NewParser()
	code := `
	let a = 1
	let b = while a != 100 {
		a = a + 1
	}
	a + b
	`
	program := parser.CreateAST(code)
	result := main.Eval(program, scope)
	assert.Equal(t, 200, result.Value())
	assert.Equal(t, main.VaIntVal, result.Kind())
}

func TestWhileLoopBreak(t *testing.T) {
	scope := main.NewScope(nil)
	parser := main.NewParser()
	code := `
	let a = 1
	while a != 100 {
		a = a + 1
		if a == 50 { break }
	}
	a
	`
	program := parser.CreateAST(code)
	result := main.Eval(program, scope)
	assert.Equal(t, 50, result.Value())
	assert.Equal(t, main.VaIntVal, result.Kind())
}

func TestArray(t *testing.T) {
	scope := main.NewGlobalScope()
	parser := main.NewParser()
	code, err := os.ReadFile("./sample/array.blu")
	assert.NoError(t, err)
	program := parser.CreateAST(string(code))
	result := main.Eval(program, scope)
	assert.Equal(t, 27, result.Value())
	assert.Equal(t, main.VaIntVal, result.Kind())
}

func TestVietnamese(t *testing.T) {
	scope := main.NewGlobalScope()
	parser := main.NewParser()
	code, err := os.ReadFile("./sample/chao.blu")
	assert.NoError(t, err)
	program := parser.CreateAST(string(code))
	result := main.Eval(program, scope)
	assert.Equal(t, []main.RuntimeVal{main.NewIntVal(12), main.NewIntVal(10), main.NewIntVal(10), main.NewIntVal(32)}, result.Value())
	assert.Equal(t, main.VaArrayVal, result.Kind())
}

func TestEnglish(t *testing.T) {
	scope := main.NewGlobalScope()
	parser := main.NewParser()
	code, err := os.ReadFile("./sample/hello.blu")
	assert.NoError(t, err)
	program := parser.CreateAST(string(code))
	result := main.Eval(program, scope)
	assert.Equal(t, []main.RuntimeVal{main.NewIntVal(12), main.NewIntVal(10), main.NewIntVal(10), main.NewIntVal(32)}, result.Value())
	assert.Equal(t, main.VaArrayVal, result.Kind())
}

func TestFibonacciRecursion(t *testing.T) {
	scope := main.NewScope(nil)
	parser := main.NewParser()
	code, err := os.ReadFile("./sample/fibonacci.blu")
	assert.NoError(t, err)
	program := parser.CreateAST(string(code))
	result := main.Eval(program, scope)
	assert.Equal(t, 21, result.Value())
	assert.Equal(t, main.VaIntVal, result.Kind())
}
