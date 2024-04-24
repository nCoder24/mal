package env

import "github.com/nCoder24/mal/impls/golisp/types"

type Option func(env *Env)

func WithOuterEnv(outer *Env) Option {
	return func(env *Env) {
		env.outer = outer
	}
}

func WithData(data map[string]types.MalValue) Option {
	return func(env *Env) {
		env.data = data
	}
}
