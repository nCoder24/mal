package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"runtime"
	"slices"

	"github.com/nCoder24/mal/impls/golisp/core"
	environ "github.com/nCoder24/mal/impls/golisp/env"
	"github.com/nCoder24/mal/impls/golisp/printer"
	"github.com/nCoder24/mal/impls/golisp/reader"
	"github.com/nCoder24/mal/impls/golisp/types"
)

var debugging *bool

func init() {
	debugging = flag.Bool("debug", false, "Prints debug logs and stacktrace on panic")
	flag.Parse()
}

func handlePanic() {
	if err := recover(); err != nil {
		fmt.Print("panic: ", err)

		if *debugging {
			buf := make([]byte, 10000)
			runtime.Stack(buf, false)
			fmt.Printf("Stack trace : %s ", string(buf))
		}
	}
}

func READ(str string) (types.MalValue, error) {
	return reader.ReadStr(str)
}

func EVAL(ast types.MalValue, env *environ.Env) (types.MalValue, error) {
	for {
		list, ok := ast.(types.List)
		if !ok || len(list) == 0 {
			return evalAST(ast, env)
		}

		fnExpr, exprs := list[0], list[1:]
		fnSymbol, ok := fnExpr.(types.Symbol)

		var err error
		switch fnSymbol {
		case "def!":
			return execDef(env, exprs)
		case "let*":
			ast, env, err = execLet(env, exprs)
			if err != nil {
				return nil, err
			}
		case "do":
			ast, err = execDo(env, exprs)
			if err != nil {
				return nil, err
			}
		case "if":
			ast, err = execIf(env, exprs)
			if err != nil {
				return nil, err
			}
		case "fn*":
			return execFn(env, exprs)
		default:
			fn, err := evalSymbol(fnSymbol, env)
			if err != nil {
				return nil, err
			}

			args, err := evalExprs(exprs, env)
			if err != nil {
				return nil, err
			}

			if fn, ok := fn.(types.DefinedFunc); ok {
				ast, env, err = evalDefinedFunc(fn, args)
				if err != nil {
					return nil, err
				}

				continue
			}

			return evalFunc(fn, args)
		}
	}
}

func evalFunc(fn types.MalValue, args []types.MalValue) (types.MalValue, error) {
	if fn, ok := fn.(types.Func); ok {
		return fn(args)
	}

	return nil, fmt.Errorf("'%v' is not callable", fn)
}

func evalDefinedFunc(fn types.DefinedFunc, args []types.MalValue) (types.MalValue, *environ.Env, error) {
	if !slices.Contains(fn.Bindings, "&") && len(fn.Bindings) != len(args) {
		return nil, nil, fmt.Errorf("expected %d arguments, got %d", len(fn.Bindings), len(args))
	}

	env := environ.New(environ.WithOuterEnv(fn.Env.(*environ.Env)), environ.WithBindings(fn.Bindings, args))

	return fn.Body, env, nil
}

func evalAST(mal types.MalValue, env *environ.Env) (types.MalValue, error) {
	switch v := mal.(type) {
	case types.Symbol:
		return evalSymbol(v, env)
	case types.List:
		return evalExprs(v, env)
	case types.Vector:
		return evalExprs(v, env)
	case types.Map:
		return evalExprs(v, env)
	default:
		return mal, nil
	}
}

func evalExprs[T types.List | types.Vector | types.Map | []types.MalValue](
	exprs T, env *environ.Env,
) (T, error) {
	values := make([]types.MalValue, 0, len(exprs))
	for _, value := range exprs {
		resolved, err := EVAL(value, env)
		if err != nil {
			return nil, err
		}

		values = append(values, resolved)
	}

	return T(values), nil
}

func evalSymbol(symbol types.Symbol, env *environ.Env) (types.MalValue, error) {
	return env.Get(string(symbol))
}

func PRINT(mal types.MalValue) string {
	return printer.PrStr(mal, true)
}

func rep(exp string, env *environ.Env) string {
	defer handlePanic()

	mal, err := READ(exp)
	if err != nil {
		return fmt.Sprintf("syntax error: %v", err)
	}

	mal, err = EVAL(mal, env)
	if err != nil {
		return fmt.Sprintf("runtime error: %v", err)
	}

	return PRINT(mal)
}

func prompt(scanner *bufio.Scanner) bool {
	fmt.Print("user> ")
	return scanner.Scan()
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	env := environ.New(environ.WithLookup(core.Namespace))

	// builtin definitions written in language itself
	rep("(def! not (fn* (a) (if a false true)))", env)

	for prompt(scanner) {
		fmt.Println(rep(scanner.Text(), env))
	}
}
