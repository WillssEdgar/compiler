package environment

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	return &Environment{
		store: make(map[string]Object),
		outer: nil,
	}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (e *Environment) Get(name string) (Object, bool) {
	if val, ok := e.store[name]; ok {
		return val, true
	}
	if e.outer != nil {
		return e.outer.Get(name)
	}
	return nil, false
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
