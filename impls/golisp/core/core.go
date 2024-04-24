package core

import (
	"fmt"

	"github.com/nCoder24/mal/impls/golisp/printer"
	"github.com/nCoder24/mal/impls/golisp/types"
)

var Namespace = map[string]types.MalValue{
	"+": types.Func(func(args []types.MalValue) (types.MalValue, error) {
		nums, err := Numbers(args)
		if err != nil {
			return nil, err
		}

		res := types.Number(0)
		for _, num := range nums {
			res += num
		}

		return res, nil
	}),
	"-": types.Func(func(args []types.MalValue) (types.MalValue, error) {
		nums, err := Numbers(args)
		if err != nil {
			return nil, err
		}

		res := nums[0]
		for _, num := range nums[1:] {
			res -= num
		}

		return res, nil
	}),
	"*": types.Func(func(args []types.MalValue) (types.MalValue, error) {
		nums, err := Numbers(args)
		if err != nil {
			return nil, err
		}

		res := types.Number(1)
		for _, num := range nums {
			res *= num
		}

		return res, nil
	}),
	"/": types.Func(func(args []types.MalValue) (types.MalValue, error) {
		nums, err := Numbers(args)
		if err != nil {
			return nil, err
		}

		res := nums[0]
		for _, num := range nums[1:] {
			res /= num
		}

		return res, nil
	}),
	"=": types.Func(func(args []types.MalValue) (types.MalValue, error) {
		return deepEqual(args[0], args[1]), nil
	}),
	"<": types.Func(func(args []types.MalValue) (types.MalValue, error) {
		nums, err := Numbers(args)
		if err != nil {
			return nil, err
		}

		for i := 1; i < len(nums); i++ {
			if !(nums[i-1] < nums[i]) {
				return types.Bool(false), nil
			}
		}

		return types.Bool(true), nil
	}),
	">": types.Func(func(args []types.MalValue) (types.MalValue, error) {
		nums, err := Numbers(args)
		if err != nil {
			return nil, err
		}

		for i := 1; i < len(nums); i++ {
			if !(nums[i-1] > nums[i]) {
				return types.Bool(false), nil
			}
		}

		return types.Bool(true), nil
	}),
	"<=": types.Func(func(args []types.MalValue) (types.MalValue, error) {
		nums, err := Numbers(args)
		if err != nil {
			return nil, err
		}

		for i := 1; i < len(nums); i++ {
			if !(nums[i-1] <= nums[i]) {
				return types.Bool(false), nil
			}
		}

		return types.Bool(true), nil
	}),
	">=": types.Func(func(args []types.MalValue) (types.MalValue, error) {
		nums, err := Numbers(args)
		if err != nil {
			return nil, err
		}

		for i := 1; i < len(nums); i++ {
			if !(nums[i-1] >= nums[i]) {
				return types.Bool(false), nil
			}
		}

		return types.Bool(true), nil
	}),
	"list": types.Func(func(args []types.MalValue) (types.MalValue, error) {
		return types.List(args), nil
	}),
	"list?": types.Func(func(args []types.MalValue) (types.MalValue, error) {
		_, ok := args[0].(types.List)
		return types.Bool(ok), nil
	}),
	"empty?": types.Func(func(args []types.MalValue) (types.MalValue, error) {
		seq, err := types.Seq(args[0])
		return types.Bool(len(seq) == 0), err
	}),
	"count": types.Func(func(args []types.MalValue) (types.MalValue, error) {
		if _, ok := args[0].(types.NilPtr); ok {
			return types.Number(0), nil
		}

		seq, err := types.Seq(args[0])
		return types.Number(len(seq)), err
	}),
	"prn": types.Func(func(args []types.MalValue) (types.MalValue, error) {
		prStrs := make([]any, 0, len(args))
		for _, arg := range args {
			prStrs = append(prStrs, printer.PrStr(arg))
		}

		fmt.Println(prStrs...)

		return types.Nil, nil
	}),
}

func Numbers(mals []types.MalValue) ([]types.Number, error) {
	nums := make([]types.Number, 0, len(mals))

	for _, arg := range mals {
		i, ok := arg.(types.Number)
		if !ok {
			return nil, fmt.Errorf("expected number, got '%v'", arg)
		}

		nums = append(nums, i)
	}

	return nums, nil
}

func deepEqual(a, b types.MalValue) types.Bool {
	seqA, seqErrA := types.Seq(a)
	seqB, seqErrB := types.Seq(b)

	if seqErrA != nil && seqErrB != nil {
		return a == b
	}

	if seqErrA != nil || seqErrB != nil || len(seqA) != len(seqB) {
		return false
	}

	for i := range seqA {
		if !deepEqual(seqA[i], seqB[i]) {
			return false
		}
	}

	return true
}
