package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkeyInterpreter/pkg/lexer"
	"monkeyInterpreter/pkg/token"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	var scanner = bufio.NewScanner(in)

	for {
		fmt.Fprintf(out, PROMPT)

		var scanned = scanner.Scan()

		if !scanned {
			return
		}

		var line = scanner.Text()
		var lex = lexer.New(line)

		for tok := lex.NextToken(); tok.Type != token.EOF; tok = lex.NextToken() {
			fmt.Fprintf(out, "%+v\n", tok)
		}
	}
}
