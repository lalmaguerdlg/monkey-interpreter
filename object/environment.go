package object

import (
	"fmt"
)

type Environment struct {
	outer *Environment
	store map[string]Object
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{
		store: s,
		outer: nil,
	}
}

func NewEnclousedEnvironment(outer *Environment) *Environment {
	s := make(map[string]Object)
	return &Environment{
		store: s,
		outer: outer,
	}
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		return e.outer.Get(name)
	}
	return obj, ok
}

// Assign recursively searches for a variable where to store it
func (e *Environment) Assign(name string, value Object) Object {
	_, found := e.store[name]
	if found {
		e.store[name] = value
		return value
	}

	if e.outer == nil {
		return &Error{Message: fmt.Sprintf("assign to an undefined identifier %s", name)}
	}

	return e.outer.Assign(name, value)
}

// Set only assigns a value in the current scope, used for let statements
// let a = 1;
func (e *Environment) Set(name string, value Object) Object {
	e.store[name] = value
	return value
}
