package repl

import (
	"bufio"
	"fmt"
	"io"

	"necronet.info/interpreter/evaluator"
	"necronet.info/interpreter/lexer"
	"necronet.info/interpreter/object"
	"necronet.info/interpreter/parser"
)

const PROMPT = "$"

const MONKEY_FACE = `
.--. .-" "-. .--.
/..\/ .-..-. \/..\ | | '| / Y \ |' | | |\\\0|0///| \ '- ,\.-"""""""-./, -' / ''-' /_ ^ ^ _\ '-''
| \._ _./ |
\ \'~'/ /
'._ '-=-' _.' '-----'
`

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()
        env := object.NewEnvironment()

		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

        evaluated := evaluator.Eval(program, env)
        if evaluated != nil {
            io.WriteString(out, evaluated.Inspect())
            io.WriteString(out, "\n")
	    }
    }
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, MONKEY_FACE)
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
