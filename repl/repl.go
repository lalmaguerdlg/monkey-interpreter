// Package repl
package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/lexer"
	"monkey/parser"
	"monkey/token"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)

		fmt.Fprintln(out, "Token:")
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Fprintf(out, "%+v\n", tok)
		}
		l.Reset()
		p := parser.New(l)
		program := p.ParseProgram()

		errors := p.Errors()
		if len(errors) > 0 {
			fmt.Fprintln(out, "parsing errors:")
			for _, error := range errors {
				fmt.Fprintln(out, error)
			}
		} else {
			fmt.Fprintln(out, "AST:")
			fmt.Fprintln(out, program.ASTDebugString())
		}
	}
}
