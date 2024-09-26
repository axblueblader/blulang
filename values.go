package main

type ValueType string

const (
	VaIntVal        ValueType = "IntVal"
	VaBoolVal       ValueType = "BoolVal"
	VaNullVal       ValueType = "NullVal"
	VaStringVal     ValueType = "StringVal"
	VaArrayVal      ValueType = "ArrayVal"
	VaFuncVal       ValueType = "FuncVal"
	VaNativeFuncVal ValueType = "NativeFuncVal"
	VaBreakVal      ValueType = "BreakVal"
	VaReturnVal     ValueType = "ReturnVal"
	VaObjectVal     ValueType = "ObjectVal"
)

type RuntimeVal interface {
	Kind() ValueType
	Value() any
}
type IntVal struct {
	value int
}

func (v IntVal) Kind() ValueType {
	return VaIntVal
}
func (v IntVal) Value() any {
	return v.value
}

func NewIntVal(value int) IntVal {
	return IntVal{value: value}
}

type BoolVal struct {
	value bool
}

func (v BoolVal) Kind() ValueType {
	return VaBoolVal
}
func (v BoolVal) Value() any {
	return v.value
}

func NewBoolVal(value bool) BoolVal {
	return BoolVal{value: value}
}

var TrueVal = NewBoolVal(true)
var FalseVal = NewBoolVal(false)

type NullVal struct {
}

func (v NullVal) Kind() ValueType {
	return VaNullVal
}
func (v NullVal) Value() any {
	return nil
}

type StringVal struct {
	value string
}

func (v StringVal) Kind() ValueType {
	return VaStringVal
}
func (v StringVal) Value() any {
	return v.value
}
func NewStringVal(value string) StringVal {
	return StringVal{value: value}
}

type ArrayVal struct {
	values []RuntimeVal
}

func (v ArrayVal) Kind() ValueType {
	return VaArrayVal
}
func (v ArrayVal) Value() any {
	return v.values
}
func NewArrayVal(values []RuntimeVal) ArrayVal {
	return ArrayVal{values: values}
}

type FunctionVal struct {
	name      string
	arguments []Identifier
	body      []Statement
}

func (v FunctionVal) Kind() ValueType {
	return VaFuncVal
}
func (v FunctionVal) Value() any {
	return v
}

func NewFuncVal(name string, args []Identifier, body []Statement) FunctionVal {
	return FunctionVal{
		name:      name,
		arguments: args,
		body:      body,
	}
}

type NativeFuncVal struct {
	call NativeFunc
}

type NativeFunc func(scope *Scope, args ...RuntimeVal) RuntimeVal

func (v NativeFuncVal) Kind() ValueType {
	return VaNativeFuncVal
}

func (v NativeFuncVal) Value() any {
	return v
}

func (v NativeFuncVal) Invoke(scope *Scope, args ...RuntimeVal) RuntimeVal {
	return v.call(scope, args...)
}

func NewNativeFuncVal(call NativeFunc) NativeFuncVal {
	return NativeFuncVal{call: call}
}

type BreakVal struct {
	lastValue RuntimeVal
}

func (v BreakVal) Kind() ValueType {
	return VaBreakVal
}

func (v BreakVal) Value() any {
	return v.lastValue
}
func NewBreakVal(lastValue RuntimeVal) BreakVal {
	return BreakVal{lastValue: lastValue}
}

type ReturnVal struct {
	lastValue RuntimeVal
}

func (v ReturnVal) Kind() ValueType {
	return VaReturnVal
}

func (v ReturnVal) Value() any {
	return v.lastValue
}
func NewReturnVal(lastValue RuntimeVal) ReturnVal {
	return ReturnVal{lastValue: lastValue}
}

type ObjectVal struct {
	properties *Scope
}

func (v ObjectVal) Kind() ValueType {
	return VaObjectVal
}

func (v ObjectVal) Value() any {
	return v.properties.variables
}

func NewObjectVal(props *Scope) ObjectVal {
	return ObjectVal{properties: props}
}
