package env

import (
	"github.com/nCoder24/mal/impls/golisp/types"
)

type Option func(env *Env)

func WithOuterEnv(outer *Env) Option {
	return func(env *Env) {
		env.outer = outer
	}
}

func WithLookup(lookup map[string]types.MalValue) Option {
	return func(env *Env) {
		env.lookup = lookup
	}
}

func WithBindings(symbols []string, exprs []types.MalValue) Option {
	return func(env *Env) {
		for i, symbol := range symbols {
			if symbol == "&" {
				env.lookup[symbols[i+1]] = types.List(exprs[i:])
				return
			}

			env.lookup[symbol] = exprs[i]
		}
	}
}
