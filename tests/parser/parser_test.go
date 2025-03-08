package parser

import (
	"fmt"
	"monkeyInterpreter/pkg/ast"
	"monkeyInterpreter/pkg/lexer"
	"monkeyInterpreter/pkg/parser"
	"testing"
)

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		//{"let y = 10;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		var lex = lexer.New(tt.input)
		var parsr = parser.New(lex)

		var program = parsr.ParseProgram()

		checkParserErrors(t, parsr)

		stmt := program.Statements[0]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.LetStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return foobar;", "foobar"},
	}

	for _, tt := range tests {
		lex := lexer.New(tt.input)
		parse := parser.New(lex)
		program := parse.ParseProgram()
		checkParserErrors(t, parse)

		if len(program.Statements) != 1 {
			t.Fatalf("program does not have expected number of statements.  Expected 1, got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]
		returnStmt, ok := stmt.(*ast.ReturnStatement)

		if !ok {
			t.Fatalf("stmt not *ast.ReturnStatement. got=%T", stmt)
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Fatalf("returnStmt.TokenLiteral not 'return', got %q", returnStmt.TokenLiteral())
		}

		if testLiteralExpression(t, returnStmt.ReturnValue, tt.expectedValue) {
			return
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

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

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

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)

	if !ok {
		t.Fatalf("stmt.Expression is not ast.Identifier. got=%T", stmt.Expression)
	}

	if literal.Value != 5 {
		t.Errorf("literal.Value was not %d, got=%d", 5, literal.Value)
	}

	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral was not %s, got=%s", "5", literal.TokenLiteral())
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"-10;", "-", 10},
		//{"!true;", "!", true},
		//{"!false;", "!", false},
		{"!foobar;", "!", "foobar"},
		{"-foobar;", "-", "foobar"},
		//{"-(5 + 5);", "-", "5 + 5"},
	}

	for _, tt := range prefixTests {
		lex := lexer.New(tt.input)
		parsr := parser.New(lex)

		program := parsr.ParseProgram()
		checkParserErrors(t, parsr)

		if len(program.Statements) != 1 {
			t.Fatalf("program does not have expected number of statements.  Expected 1, got=%d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("prgram.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)

		if !ok {
			t.Fatalf("stmt.Expression is not ast.PrefixExpression. got=%T", stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Fatalf("ex.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}

		if !testLiteralExpression(t, exp.Right, tt.value) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"foobar + barfoo;", "foobar", "+", "barfoo"},
		{"foobar - barfoo;", "foobar", "-", "barfoo"},
		{"foobar * barfoo;", "foobar", "*", "barfoo"},
		{"foobar / barfoo;", "foobar", "/", "barfoo"},
		{"foobar > barfoo;", "foobar", ">", "barfoo"},
		{"foobar < barfoo;", "foobar", "<", "barfoo"},
		{"foobar == barfoo;", "foobar", "==", "barfoo"},
		{"foobar != barfoo;", "foobar", "!=", "barfoo"},
	}

	for _, tt := range infixTests {
		lex := lexer.New(tt.input)
		parsr := parser.New(lex)

		program := parsr.ParseProgram()
		checkParserErrors(t, parsr)

		if len(program.Statements) != 1 {
			t.Fatalf("program does not have expected number of statements.  Expected 1, got=%d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("prgram.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		if !testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b / c", "(a + (b / c))"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
	}

	for _, tt := range tests {
		lex := lexer.New(tt.input)
		parsr := parser.New(lex)

		program := parsr.ParseProgram()
		checkParserErrors(t, parsr)

		actual := program.String()

		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
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

func testInfixExpression(t *testing.T, expression ast.Expression, left interface{}, operator string, right interface{}) bool {
	opxExpres, ok := expression.(*ast.InfixExpression)

	if !ok {
		t.Errorf("expected expression to be an infix expression but got=%T(%s)", expression, expression)
		return false
	}

	if !testLiteralExpression(t, opxExpres.Left, left) {
		return false
	}

	if opxExpres.Operator != operator {
		t.Errorf("expected operator to be %s but got=%q", operator, opxExpres.Operator)
		return false
	}

	if !testLiteralExpression(t, opxExpres.Right, right) {
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, expression ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, expression, int64(v))
	case int64:
		return testIntegerLiteral(t, expression, v)
	case string:
		return testIdentifier(t, expression, v)
	}

	t.Errorf("expression type %T was not handled", expected)
	return false
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)

	if !ok {
		t.Errorf("expected il to be an IntegerLiteral. got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("expected integer value to be %d, but got=%d", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("expected integ.TokenLiteral to be %d, but got=%s", value, integ.TokenLiteral())
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)

	if !ok {
		t.Errorf("expected expression to be ast.Identifier, got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("expected identity value to be %s, got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("expected identity.TokenLiteral to be %s, got %s", value, ident.TokenLiteral())
		return false
	}

	return true
}
