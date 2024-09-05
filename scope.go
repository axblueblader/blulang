package main

import (
	"errors"
	"fmt"
	"log"
	"os"
)

type Scope struct {
	parent    *Scope
	variables map[string]RuntimeVal
}

func NewScope(parent *Scope) *Scope {
	return &Scope{
		parent:    parent,
		variables: make(map[string]RuntimeVal),
	}
}

func NewGlobalScope() *Scope {
	globalScope := NewScope(nil)
	globalScope.DeclareVar("true", NewBoolVal(true))
	globalScope.DeclareVar("false", NewBoolVal(false))
	globalScope.DeclareVar("đúng", NewBoolVal(true))
	globalScope.DeclareVar("sai", NewBoolVal(false))
	globalScope.DeclareVar("print", PrintFunc)
	globalScope.DeclareVar("in", PrintFunc)
	globalScope.DeclareVar("count", CountFunc)
	globalScope.DeclareVar("đếm", CountFunc)
	globalScope.DeclareVar("input", InputFunc)
	globalScope.DeclareVar("nhập", InputFunc)
	globalScope.DeclareVar("abs", NewNativeFuncVal(func(scope *Scope, args ...RuntimeVal) RuntimeVal {
		if args[0].Value().(int) < 0 {
			return NewIntVal(-args[0].Value().(int))
		}
		return args[0]
	}))
	globalScope.DeclareVar("exit", NewNativeFuncVal(func(scope *Scope, args ...RuntimeVal) RuntimeVal {
		fmt.Println("Good bye")
		os.Exit(0)
		return NullVal{}
	}))
	return globalScope
}

func (s *Scope) DeclareVar(name string, value RuntimeVal) RuntimeVal {
	if s.variables[name] != nil {
		log.Panic("variable already defined")
	}

	s.variables[name] = value
	return value
}

func (s *Scope) AssignVar(name string, value RuntimeVal) RuntimeVal {
	scope, err := s.resolve(name)
	if err != nil {
		return NullVal{}
	}
	scope.variables[name] = value
	return value
}

func (s *Scope) GetVarVal(name string) RuntimeVal {
	scope, err := s.resolve(name)
	if err != nil {
		return NullVal{}
	}
	return scope.variables[name]
}

func (s *Scope) resolve(name string) (*Scope, error) {
	if s.variables[name] != nil {
		return s, nil
	}

	if s.parent == nil {
		return nil, errors.New("variable not found")
	}

	return s.parent.resolve(name)
}
