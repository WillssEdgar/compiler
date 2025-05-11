package environment

import "compiler/object"

type Environment struct {
	store map[string]object.Object
}

func NewEnvironment() *Environment {
	return &Environment{store: make(map[string]object.Object)}
}

func (e *Environment) Get(name string) (object.Object, bool) {
	val, ok := e.store[name]
	return val, ok
}

func (e *Environment) Set(name string, val object.Object) object.Object {
	e.store[name] = val
	return val
}
