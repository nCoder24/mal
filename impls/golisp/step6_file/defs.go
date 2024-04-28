package main

import (
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

func execLet(env *environ.Env, args []types.MalValue) (types.MalValue, *environ.Env, error) {
	letEnv := environ.New(environ.WithOuterEnv(env))

	bindings, err := types.Sequence(args[0])
	if err != nil {
		return nil, nil, err
	}

	for i := 0; i < len(bindings); i += 2 {
		val, err := EVAL(bindings[i+1], letEnv)
		if err != nil {
			return nil, nil, err
		}

		str, err := types.SymbolString(bindings[i])
		if err != nil {
			return nil, nil, err
		}

		letEnv.Set(str, val)
	}

	return args[1], letEnv, nil
}

func execDo(env *environ.Env, args []types.MalValue) (types.MalValue, error) {
	for _, form := range args[:len(args)-1] {
		_, err := EVAL(form, env)

		if err != nil {
			return nil, err
		}
	}

	return args[len(args)-1], nil
}

func execIf(env *environ.Env, args []types.MalValue) (types.MalValue, error) {
	pred, err := EVAL(args[0], env)
	if err != nil {
		return nil, err
	}

	if isTrue, ok := pred.(types.Bool); pred == types.Nil || ok && !bool(isTrue) {
		if len(args) > 2 {
			return args[2], nil
		}

		return types.Nil, nil
	}

	return args[1], nil
}

func execFn(env *environ.Env, args []types.MalValue) (types.MalValue, error) {
	symbols, err := types.SymbolStrings(args[0])
	if err != nil {
		return nil, err
	}

	fn := types.DefinedFunc{
		Env:      env,
		Bindings: symbols,
		Body:     args[1],
	}

	return fn, nil
}
