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

	env.Set(string(args[0].(types.Symbol)), val)

	return val, nil
}

func execLet(env *environ.Env, args []types.MalValue) (types.MalValue, error) {
	letEnv := environ.New(environ.WithOuterEnv(env))

	bindings, err := types.Seq(args[0])
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(bindings); i += 2 {
		val, err := EVAL(bindings[i+1], letEnv)
		if err != nil {
			return nil, err
		}

		letEnv.Set(string(bindings[i].(types.Symbol)), val)
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
	return types.Func(func(argList []types.MalValue) (types.MalValue, error) {
		fnEnv := environ.New(environ.WithOuterEnv(env))
		bindings, err := types.Seq(args[0])
		if err != nil {
			return nil, err
		}

		bindTill, hasVariadicParam := slices.Index(bindings, types.MalValue(types.Symbol("&"))), true
		if bindTill == -1 {
			bindTill, hasVariadicParam = len(bindings), false
		}

		if !hasVariadicParam && len(bindings) != len(argList) {
			return nil, fmt.Errorf("expected %d arguments, got %d", len(bindings), len(argList))
		}

		values, err := evalMultiple(env, argList...)
		if err != nil {
			return nil, err
		}

		for i := 0; i < bindTill; i++ {
			fnEnv.Set(string(bindings[i].(types.Symbol)), values[i])
		}

		if hasVariadicParam {
			fnEnv.Set(string(bindings[bindTill+1].(types.Symbol)), types.List(values[bindTill:]))
		}

		return EVAL(args[1], fnEnv)
	}), nil
}

func evalMultiple(env *environ.Env, exprs ...types.MalValue) ([]types.MalValue, error) {
	values := make([]types.MalValue, 0, len(exprs))

	for _, form := range exprs {
		value, err := EVAL(form, env)
		if err != nil {
			return nil, err
		}

		values = append(values, value)
	}

	return values, nil
}
