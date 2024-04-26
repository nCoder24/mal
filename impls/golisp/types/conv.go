package types

import (
	"fmt"
	"strings"
)

func stringify(forms []MalValue) string {
	strs := make([]string, 0, len(forms))

	for _, v := range forms {
		strs = append(strs, fmt.Sprintf("%v", v))
	}

	return strings.Join(strs, " ")
}

func SymbolString(val MalValue) (string, error) {
	sym, ok := val.(Symbol)
	if !ok {
		return "", fmt.Errorf("'%v' is not a symbol", val)
	}

	return string(sym), nil
}

func SymbolStrings(vals MalValue) ([]string, error) {
	seq, err := Sequence(vals)
	if err != nil {
		return nil, err
	}

	symbols := make([]string, len(seq))
	for i, bind := range seq {
		symbol, err := SymbolString(bind)
		if err != nil {
			return nil, err
		}

		symbols[i] = symbol
	}

	return symbols, nil
}

func Numbers(mals []MalValue) ([]Number, error) {
	nums := make([]Number, 0, len(mals))

	for _, arg := range mals {
		i, ok := arg.(Number)
		if !ok {
			return nil, fmt.Errorf("expected number, got '%v'", arg)
		}

		nums = append(nums, i)
	}

	return nums, nil
}

func Sequence(val MalValue) ([]MalValue, error) {
	switch v := val.(type) {
	case List:
		return v, nil
	case Vector:
		return v, nil
	}

	return nil, fmt.Errorf("cannot convert '%v' to seq", val)
}
