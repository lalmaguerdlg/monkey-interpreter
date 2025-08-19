// Package repl
package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	env := object.NewEnvironment()

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)

		// for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
		// fmt.Fprintln(out, "Token:")
		// fmt.Fprintf(out, "%+v\n", tok)
		// }
		// l.Reset()
		p := parser.New(l)
		program := p.ParseProgram()
		errors := p.Errors()
		if len(errors) > 0 {
			printParseErrors(out, errors)
			continue
		}

		result := evaluator.Eval(program, env)
		if result != nil {
			fmt.Fprintln(out, result.Inspect())
		}
	}
}

func printParseErrors(out io.Writer, errors []string) {
	fmt.Fprintln(out, "parsing errors:")
	for _, error := range errors {
		fmt.Fprintf(out, "\t%s\n", error)
	}
}
