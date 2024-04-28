package core

import (
	"fmt"
	"os"
	"strings"

	"github.com/nCoder24/mal/impls/golisp/printer"
	"github.com/nCoder24/mal/impls/golisp/reader"
	"github.com/nCoder24/mal/impls/golisp/types"
)

func sum(args []types.MalValue) (types.MalValue, error) {
	nums, err := types.Numbers(args)
	if err != nil {
		return nil, err
	}

	res := types.Number(0)
	for _, num := range nums {
		res += num
	}

	return res, nil
}

func sub(args []types.MalValue) (types.MalValue, error) {
	nums, err := types.Numbers(args)
	if err != nil {
		return nil, err
	}

	res := nums[0]
	for _, num := range nums[1:] {
		res -= num
	}

	return res, nil
}

func mul(args []types.MalValue) (types.MalValue, error) {
	nums, err := types.Numbers(args)
	if err != nil {
		return nil, err
	}

	res := types.Number(1)
	for _, num := range nums {
		res *= num
	}

	return res, nil
}

func divide(args []types.MalValue) (types.MalValue, error) {
	nums, err := types.Numbers(args)
	if err != nil {
		return nil, err
	}

	res := nums[0]
	for _, num := range nums[1:] {
		res /= num
	}

	return res, nil
}

func equal(args []types.MalValue) (types.MalValue, error) {
	return types.DeepEqual(args[0], args[1]), nil
}

func lt(args []types.MalValue) (types.MalValue, error) {
	nums, err := types.Numbers(args)
	if err != nil {
		return nil, err
	}

	for i := 1; i < len(nums); i++ {
		if !(nums[i-1] < nums[i]) {
			return types.Bool(false), nil
		}
	}

	return types.Bool(true), nil
}

func gt(args []types.MalValue) (types.MalValue, error) {
	nums, err := types.Numbers(args)
	if err != nil {
		return nil, err
	}

	for i := 1; i < len(nums); i++ {
		if !(nums[i-1] > nums[i]) {
			return types.Bool(false), nil
		}
	}

	return types.Bool(true), nil
}

func lte(args []types.MalValue) (types.MalValue, error) {
	isGt, err := gt(args)
	if err != nil {
		return nil, err
	}

	return !isGt.(types.Bool), nil
}

func gte(args []types.MalValue) (types.MalValue, error) {
	isLt, err := lt(args)
	if err != nil {
		return nil, err
	}

	return !isLt.(types.Bool), nil
}

func list(args []types.MalValue) (types.MalValue, error) {
	return types.List(args), nil
}

func isList(args []types.MalValue) (types.MalValue, error) {
	_, ok := args[0].(types.List)
	return types.Bool(ok), nil
}

func count(args []types.MalValue) (types.MalValue, error) {
	if _, ok := args[0].(types.NilPtr); ok {
		return types.Number(0), nil
	}

	seq, err := types.Sequence(args[0])
	return types.Number(len(seq)), err
}

func isEmpty(args []types.MalValue) (types.MalValue, error) {
	c, err := count(args)
	if err == nil {
		return c.(types.Number) == 0, nil
	}

	return types.Number(0), err
}

func prStr(args []types.MalValue) (types.MalValue, error) {
	strs := make([]string, len(args))
	for i, arg := range args {
		strs[i] = printer.PrStr(arg, true)
	}

	return types.String(strings.Join(strs, " ")), nil
}

func prn(args []types.MalValue) (types.MalValue, error) {
	s, _ := prStr(args)
	fmt.Println(s)

	return types.Nil, nil
}

func str(args []types.MalValue) (types.MalValue, error) {
	strs := make([]string, len(args))
	for i, arg := range args {
		strs[i] = printer.PrStr(arg, false)
	}

	return types.String(strings.Join(strs, "")), nil
}

func printLine(args []types.MalValue) (types.MalValue, error) {
	strs := make([]string, len(args))
	for i, arg := range args {
		strs[i] = printer.PrStr(arg, false)
	}

	fmt.Println(strings.Join(strs, " "))

	return types.Nil, nil
}

func readStr(args []types.MalValue) (types.MalValue, error) {
	if raw, ok := args[0].(types.String); ok {
		return reader.ReadStr(string(raw))
	}

	return nil, fmt.Errorf("argument must be a string")
}

func slurp(args []types.MalValue) (types.MalValue, error) {
	if path, ok := args[0].(types.String); ok {
		content, err := os.ReadFile(string(path))
		if err != nil {
			return nil, fmt.Errorf("error while reading file '%s': %w", content, err)
		}

		return types.String(content), nil
	}

	return nil, fmt.Errorf("argument must be a string")
}
