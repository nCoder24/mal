package env

import (
	"fmt"

	"github.com/nCoder24/mal/impls/golisp/types"
)

type Env struct {
	outer *Env
	data  map[string]types.MalValue
}

func New(options ...Option) *Env {
	env := &Env{
		data: make(map[string]types.MalValue),
	}

	for _, option := range options {
		option(env)
	}

	return env
}

func (e *Env) find(key string) *Env {
	if _, ok := e.data[key]; ok {
		return e
	}

	if e.outer != nil {
		return e.outer.find(key)
	}

	return nil
}

func (e *Env) Get(key string) (types.MalValue, error) {
	if env := e.find(key); env != nil {
		return env.data[key], nil
	}

	return nil, fmt.Errorf("symbol '%s' not found", key)
}

func (e *Env) Set(key string, value types.MalValue) {
	e.data[key] = value
}
