package core

import (
	"fmt"

	"github.com/nCoder24/mal/impls/golisp/types"
)

var Namespace = map[string]types.MalValue{
	"+": types.Func(func(args []types.MalValue) (types.MalValue, error) {
		ints, err := Numbers(args)
		if err != nil {
			return nil, err
		}

		res := types.Number(0)
		for _, i := range ints {
			res += i
		}

		return res, nil
	}),
	"-": types.Func(func(args []types.MalValue) (types.MalValue, error) {
		ints, err := Numbers(args)
		if err != nil {
			return nil, err
		}

		res := ints[0]
		for _, i := range ints[1:] {
			res -= i
		}

		return res, nil
	}),
	"*": types.Func(func(args []types.MalValue) (types.MalValue, error) {
		ints, err := Numbers(args)
		if err != nil {
			return nil, err
		}

		res := types.Number(1)
		for _, i := range ints {
			res *= i
		}

		return res, nil
	}),
	"/": types.Func(func(args []types.MalValue) (types.MalValue, error) {
		ints, err := Numbers(args)
		if err != nil {
			return nil, err
		}

		res := ints[0]
		for _, i := range ints[1:] {
			res /= i
		}

		return res, nil
	}),
}

func Numbers(mals []types.MalValue) ([]types.Number, error) {
	ints := make([]types.Number, 0, len(mals))

	for _, arg := range mals {
		i, ok := arg.(types.Number)
		if !ok {
			return nil, fmt.Errorf("expected number, got '%v'", arg)
		}

		ints = append(ints, i)
	}

	return ints, nil
}
