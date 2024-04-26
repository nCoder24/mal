package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/nCoder24/mal/impls/golisp/printer"
	"github.com/nCoder24/mal/impls/golisp/reader"
	"github.com/nCoder24/mal/impls/golisp/types"
)

func READ(str string) (types.MalValue, error) {
	return reader.ReadStr(str)
}

func EVAL(mal types.MalValue, env map[string]types.MalValue) (types.MalValue, error) {
	list, isList := mal.(types.List)
	resolved, err := evalAST(mal, env)

	if err != nil || !isList || len(list) == 0 {
		return resolved, err
	}

	resolvedList := resolved.(types.List)
	if f, ok := resolvedList[0].(types.Func); ok {
		return f(resolvedList[1:])
	}

	return nil, fmt.Errorf("cannot call '%s'", list[0])
}

func evalAST(mal types.MalValue, env map[string]types.MalValue) (types.MalValue, error) {
	switch v := mal.(type) {
	case types.Symbol:
		return resolveSymbol(v, env)
	case types.List:
		return resolveForms(v, env)
	case types.Vector:
		return resolveForms(v, env)
	case types.Map:
		return resolveForms(v, env)
	default:
		return mal, nil
	}
}

func resolveForms[T types.List | types.Vector | types.Map](
	coll T, env map[string]types.MalValue,
) (types.MalValue, error) {
	values := make([]types.MalValue, 0, len(coll))
	for _, value := range coll {
		resolved, err := EVAL(value, env)
		if err != nil {
			return nil, err
		}

		values = append(values, resolved)
	}

	return T(values), nil
}

func resolveSymbol(symbol types.Symbol, env map[string]types.MalValue) (types.MalValue, error) {
	if resolved, ok := env[string(symbol)]; ok {
		return resolved, nil
	}

	return nil, fmt.Errorf("could not resolve symbol '%s'", symbol)
}

func PRINT(mal types.MalValue) string {
	return printer.PrStr(mal, true)
}

func rep(exp string) string {
	mal, err := READ(exp)
	if err != nil {
		return errorString(err)
	}

	mal, err = EVAL(mal, env)
	if err != nil {
		return errorString(err)
	}

	return PRINT(mal)
}

func errorString(err error) string {
	return "ERROR: " + err.Error()
}

func prompt(scanner *bufio.Scanner) bool {
	fmt.Print("user> ")
	return scanner.Scan()
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for prompt(scanner) {
		fmt.Println(rep(scanner.Text()))
	}
}
