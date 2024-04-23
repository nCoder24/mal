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

func EVAL(mal types.MalValue) (types.MalValue, error) {
	return mal, nil
}

func PRINT(mal types.MalValue) string {
	return printer.PrStr(mal)
}

func rep(exp string) string {
	mal, err := READ(exp)
	if err != nil {
		return err.Error()
	}

	mal, err = EVAL(mal)
	if err != nil {
		return err.Error()
	}

	return PRINT(mal)
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
