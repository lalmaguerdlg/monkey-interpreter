package object

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

func (e *Environment) Set(name string, value Object) Object {
	_, ok := e.store[name]
	if !ok && e.outer != nil {
		return e.outer.Set(name, value)
	}
	e.store[name] = value
	return value
}

func (e *Environment) Shadow(name string, value Object) Object {
	e.store[name] = value
	return value
}
