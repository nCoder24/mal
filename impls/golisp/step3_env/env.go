package main

import (
	"fmt"

	"github.com/nCoder24/mal/impls/golisp/env"
	"github.com/nCoder24/mal/impls/golisp/types"
)

func newReplEnv() *env.Env {
	e := env.New()

	e.Set("+", types.Func(func(args []types.MalValue) (types.MalValue, error) {
		ints, err := toInts(args)
		if err != nil {
			return nil, err
		}

		res := types.Int(0)
		for _, i := range ints {
			res += i
		}

		return res, nil
	}))

	e.Set("-", types.Func(func(args []types.MalValue) (types.MalValue, error) {
		ints, err := toInts(args)
		if err != nil {
			return nil, err
		}

		res := ints[0]
		for _, i := range ints[1:] {
			res -= i
		}

		return res, nil
	}))

	e.Set("*", types.Func(func(args []types.MalValue) (types.MalValue, error) {
		ints, err := toInts(args)
		if err != nil {
			return nil, err
		}

		res := types.Int(1)
		for _, i := range ints {
			res *= i
		}

		return res, nil
	}))

	e.Set("/", types.Func(func(args []types.MalValue) (types.MalValue, error) {
		ints, err := toInts(args)
		if err != nil {
			return nil, err
		}

		res := ints[0]
		for _, i := range ints[1:] {
			res /= i
		}

		return res, nil
	}))

	return e
}

func toInts(mals []types.MalValue) ([]types.Int, error) {
	ints := make([]types.Int, 0, len(mals))

	for _, arg := range mals {
		i, ok := arg.(types.Int)
		if !ok {
			return nil, fmt.Errorf("expected int, got %v", arg)
		}

		ints = append(ints, i)
	}

	return ints, nil
}
