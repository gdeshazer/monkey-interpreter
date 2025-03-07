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

	checkParserErrors(t, parsr)

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

func TestReturnStatement(t *testing.T) {
	input := `
return 5;
return 10;
return 993322;
`

	var lex = lexer.New(input)
	var parsr = parser.New(lex)

	var program = parsr.ParseProgram()

	checkParserErrors(t, parsr)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements.  got=%d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStatement, ok := stmt.(*ast.ReturnStatement)

		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement, got %T", stmt)
			continue
		}

		if returnStatement.TokenLiteral() != "return" {
			t.Errorf("returnStatment.TokenLiteral not 'return' got %q", returnStatement.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	lex := lexer.New(input)
	parsr := parser.New(lex)

	program := parsr.ParseProgram()
	checkParserErrors(t, parsr)

	if len(program.Statements) != 1 {
		t.Fatalf("program does not have enough statements, got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("prgram.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)

	if !ok {
		t.Fatalf("stmt.Expression is not ast.Identifier. got=%T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("ident.Value was not %s, got=%s", "foobar", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral was not %s, got=%s", "foobar", ident.TokenLiteral())
	}
}

func checkParserErrors(t *testing.T, parsr *parser.Parser) {
	errors := parsr.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))

	for _, msg := range errors {
		t.Errorf("parser error %q", msg)
	}

	t.FailNow()
}
