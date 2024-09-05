package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var PrintFunc = NewNativeFuncVal(func(scope *Scope, args ...RuntimeVal) RuntimeVal {
	var values []interface{}
	for _, v := range args {
		values = append(values, v.Value())
	}
	fmt.Println(values...)
	return NullVal{}
})
var InputFunc = NewNativeFuncVal(func(scope *Scope, args ...RuntimeVal) RuntimeVal {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	intVal, err := strconv.Atoi(text)
	if err == nil {
		return NewIntVal(intVal)
	}
	return NewStringVal(text)
})

var CountFunc = NewNativeFuncVal(func(scope *Scope, args ...RuntimeVal) RuntimeVal {
	if args[0].Kind() == VaArrayVal {
		return NewIntVal(len(args[0].(ArrayVal).values))
	}
	return NewIntVal(0)
})
