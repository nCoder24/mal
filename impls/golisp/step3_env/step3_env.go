package main

import (
	"bufio"
	"fmt"
	"os"

	environment "github.com/nCoder24/mal/impls/golisp/env"
	"github.com/nCoder24/mal/impls/golisp/printer"
	"github.com/nCoder24/mal/impls/golisp/reader"
	"github.com/nCoder24/mal/impls/golisp/types"
)

func READ(str string) (types.MalValue, error) {
	return reader.ReadStr(str)
}

func EVAL(mal types.MalValue, env *environment.Env) (types.MalValue, error) {
	list, isList := mal.(types.List)
	if !isList || len(list) == 0 {
		return evalAST(mal, env)
	}

	switch list[0].(types.Symbol) {
	case "def!":
		return evalDef(mal, env)
	case "let*":
		return evalLet(mal, env)
	}

	return evalList(mal, env)
}

func evalList(mal types.MalValue, env *environment.Env) (types.MalValue, error) {
	resolved, err := evalAST(mal, env)
	if err != nil {
		return nil, err
	}

	resolvedList := resolved.(types.List)
	if f, ok := resolvedList[0].(types.Func); ok {
		return f(resolvedList[1:])
	}

	return nil, fmt.Errorf("cannot call '%s'", mal.(types.List)[0])
}

func evalDef(mal types.MalValue, env *environment.Env) (types.MalValue, error) {
	list := mal.(types.List)

	val, err := EVAL(list[2], env)
	if err != nil {
		return nil, err
	}

	env.Set(string(list[1].(types.Symbol)), val)

	return val, nil
}

func evalLet(mal types.MalValue, env *environment.Env) (types.MalValue, error) {
	letEnv := environment.NewWith(env)
	list := mal.(types.List)

	bindings, ok := types.Seq(list[1])
	if !ok {
		return nil, fmt.Errorf("bindings must be sequence")
	}

	for i := 0; i < len(bindings); i += 2 {
		val, err := EVAL(bindings[i+1], letEnv)
		if err != nil {
			return nil, err
		}

		letEnv.Set(string(bindings[i].(types.Symbol)), val)
	}

	return EVAL(list[2], letEnv)
}

func evalAST(mal types.MalValue, env *environment.Env) (types.MalValue, error) {
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
	coll T, env *environment.Env,
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

func resolveSymbol(symbol types.Symbol, env *environment.Env) (types.MalValue, error) {
	return env.Get(string(symbol))
}

func PRINT(mal types.MalValue) string {
	return printer.PrStr(mal)
}

func rep(exp string, env *environment.Env) string {
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
	env := newReplEnv()

	for prompt(scanner) {
		fmt.Println(rep(scanner.Text(), env))
	}
}
