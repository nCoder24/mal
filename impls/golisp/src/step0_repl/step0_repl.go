package main

import (
	"bufio"
	"fmt"
	"os"
)

func READ(s string) string {
	return s
}

func EVAL(s string) string {
	return s
}

func PRINT(s string) string {
	return s
}

func rep(exp string) string {
	return PRINT(EVAL(READ(exp)))
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
