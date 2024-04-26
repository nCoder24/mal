package main

import (
	"fmt"
	"slices"

	environ "github.com/nCoder24/mal/impls/golisp/env"
	"github.com/nCoder24/mal/impls/golisp/types"
)

func execDef(env *environ.Env, args []types.MalValue) (types.MalValue, error) {
	val, err := EVAL(args[1], env)
	if err != nil {
		return nil, err
	}

	sym, err := types.SymbolString(args[0])
	if err != nil {
		return nil, err
	}

	env.Set(sym, val)

	return val, nil
}

func execLet(env *environ.Env, args []types.MalValue) (types.MalValue, error) {
	letEnv := environ.New(environ.WithOuterEnv(env))

	bindings, err := types.Sequence(args[0])
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(bindings); i += 2 {
		val, err := EVAL(bindings[i+1], letEnv)
		if err != nil {
			return nil, err
		}

		str, err := types.SymbolString(bindings[i])
		if err != nil {
			return nil, err
		}

		letEnv.Set(str, val)
	}

	return EVAL(args[1], letEnv)
}

func execDo(env *environ.Env, args []types.MalValue) (types.MalValue, error) {
	var (
		res types.MalValue
		err error
	)

	for _, form := range args {
		res, err = EVAL(form, env)

		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func execIf(env *environ.Env, args []types.MalValue) (types.MalValue, error) {
	pred, err := EVAL(args[0], env)
	if err != nil {
		return nil, err
	}

	if isTrue, ok := pred.(types.Bool); pred == types.Nil || ok && !bool(isTrue) {
		if len(args) > 2 {
			return EVAL(args[2], env)
		}

		return types.Nil, nil
	}

	return EVAL(args[1], env)
}

func execFn(env *environ.Env, args []types.MalValue) (types.MalValue, error) {
	return types.Func(func(exprs []types.MalValue) (types.MalValue, error) {
		symbols, err := types.SymbolStrings(args[0])
		if err != nil {
			return nil, err
		}

		if !slices.Contains(symbols, "&") && len(symbols) != len(exprs) {
			return nil, fmt.Errorf("expected %d arguments, got %d", len(symbols), len(exprs))
		}

		return EVAL(args[1], environ.New(environ.WithOuterEnv(env), environ.WithBindings(symbols, exprs)))
	}), nil
}
