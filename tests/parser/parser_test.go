package parser

import (
	"monkeyInterpreter/pkg/ast"
	"monkeyInterpreter/pkg/lexer"
	"monkeyInterpreter/pkg/parser"
	"testing"
)

func TestLetStatements(t *testing.T) {
	var input = `
let x = 5;
let y = 10;
let foobar = 838383;
`
	var lex = lexer.New(input)
	var parsr = parser.New(lex)

	var program = parsr.ParseProgram()

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements, got=%d", len(program.Statements))
	}

	var tests = []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		statemnt := program.Statements[i]

		if !testLetStatement(t, statemnt, tt.expectedIdentifier) {
			return
		}
	}
}

func testLetStatement(t *testing.T, statemnt ast.Statement, name string) bool {
	if statemnt.TokenLiteral() != "let" {
		t.Errorf("statemnt.TokenLiteral not 'let', got=%q", statemnt.TokenLiteral())
		return false
	}

	letStmnt, ok := statemnt.(*ast.LetStatement)

	if !ok {
		t.Errorf("statemnt not *ast.LetStatement got=%T", statemnt)
	}

	if letStmnt.Name.TokenLiteral() != name {
		t.Errorf("letStmnt.Name.TokenLiteral() not %s, got=%s", name, letStmnt.Name.TokenLiteral())
		return false
	}

	return true
}
