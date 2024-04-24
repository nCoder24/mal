package main

import (
	"fmt"

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

	if isTrue, ok := pred.(bool); pred == nil || ok && !isTrue {
		if len(args) > 2 {
			return EVAL(args[2], env)
		}

		return nil, nil
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

		if len(bindings) != len(argList) {
			return nil, fmt.Errorf("expected %d arguments, got %d", len(bindings), len(argList))
		}

		for i, symbol := range bindings {
			val, err := EVAL(argList[i], env)
			if err != nil {
				return nil, err
			}

			fnEnv.Set(string(symbol.(types.Symbol)), val)
		}

		return EVAL(args[1], fnEnv)
	}), nil
}
