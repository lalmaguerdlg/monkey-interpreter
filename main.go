// Package main
package main

import (
	"fmt"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"monkey/repl"
	"os"
	"os/user"
)

func main() {
	args := os.Args[1:]

	if len(args) > 0 {
		filename := args[0]
		runFile(filename)
	} else {
		runRepl()
	}
}

func runFile(filename string) {
	dat, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Could not open file %s\n", filename)
		return
	}
	script := string(dat)
	l := lexer.New(script)
	p := parser.New(l)
	program := p.ParseProgram()
	errors := p.Errors()
	if len(errors) > 0 {
		for _, msg := range errors {
			fmt.Println(msg)
		}
		return
	}

	env := object.NewEnvironment()
	evaluation := evaluator.Eval(program, env)
	if err, ok := evaluation.(*object.Error); ok {
		fmt.Println(err.Message)
	}
}

func runRepl() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the Monkey programming language!\n", user.Username)
	fmt.Printf("Feel free to type in commands\n")
	repl.Start(os.Stdin, os.Stdout)
}
