package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/nCoder24/mal/impls/golisp/core"
	environ "github.com/nCoder24/mal/impls/golisp/env"
	"github.com/nCoder24/mal/impls/golisp/printer"
	"github.com/nCoder24/mal/impls/golisp/reader"
	"github.com/nCoder24/mal/impls/golisp/types"
)

func READ(str string) (types.MalValue, error) {
	return reader.ReadStr(str)
}

func EVAL(mal types.MalValue, env *environ.Env) (types.MalValue, error) {
	list, ok := mal.(types.List)
	if !ok || len(list) == 0 {
		return evalAST(mal, env)
	}

	fn, args := list[0], list[1:]

	if symbol, ok := fn.(types.Symbol); ok {
		switch symbol {
		case "def!":
			return execDef(env, args)
		case "let*":
			return execLet(env, args)
		case "do":
			return execDo(env, args)
		case "if":
			return execIf(env, args)
		case "fn*":
			return execFn(env, args)
		}
	}

	return exec(env, fn, args)
}

func exec(env *environ.Env, fn types.MalValue, args []types.MalValue) (types.MalValue, error) {
	f, err := EVAL(fn, env)
	if err != nil {
		return nil, err
	}

	resolvedArgs := make([]types.MalValue, 0, len(args))
	for _, arg := range args {
		v, err := EVAL(arg, env)
		if err != nil {
			return nil, err
		}

		resolvedArgs = append(resolvedArgs, v)
	}

	if callable, ok := f.(types.Func); ok {
		return callable(resolvedArgs)
	}

	return nil, fmt.Errorf("'%v' is not callable", f)
}

func evalAST(mal types.MalValue, env *environ.Env) (types.MalValue, error) {
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
	coll T, env *environ.Env,
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

func resolveSymbol(symbol types.Symbol, env *environ.Env) (types.MalValue, error) {
	return env.Get(string(symbol))
}

func PRINT(mal types.MalValue) string {
	return printer.PrStr(mal)
}

func rep(exp string, env *environ.Env) string {
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
	env := environ.New(environ.WithData(core.Namespace))

	for prompt(scanner) {
		fmt.Println(rep(scanner.Text(), env))
	}
}
