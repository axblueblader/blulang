package main

import "fmt"

var PrintFunc = NewNativeFuncVal(func(scope *Scope, args ...RuntimeVal) RuntimeVal {
	var values []interface{}
	for _, v := range args {
		values = append(values, v.Value())
	}
	fmt.Println(values...)
	return NullVal{}
})

var CountFunc = NewNativeFuncVal(func(scope *Scope, args ...RuntimeVal) RuntimeVal {
	if args[0].Kind() == VaArrayVal {
		return NewIntVal(len(args[0].(ArrayVal).values))
	}
	return NewIntVal(0)
})
