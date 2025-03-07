package parser

import (
	"monkeyInterpreter/pkg/ast"
	"monkeyInterpreter/pkg/lexer"
	"monkeyInterpreter/pkg/token"
)

type Parser struct {
	lex *lexer.Lexer

	// currentToken represents the token we are currently parsing
	currentToken token.Token

	// peekToken is the next token returned from lex.  This allows us to look ahead when forming the AST
	peekToken token.Token
}

func New(lex *lexer.Lexer) *Parser {
	var p = &Parser{lex: lex}

	// read two tokens to set current and peek
	p.nextToken()
	p.nextToken()

	return p
}

// nextToken advances the parser's tokens by setting the currentToken to the peekToken, and then setting peekToken to
// the next token retrieved by lex
func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lex.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}

	program.Statements = []ast.Statement{}

	for !p.currentTokenIs(token.EOF) {
		statemnt := p.parseStatement()

		if statemnt != nil {
			program.Statements = append(program.Statements, statemnt)
		}

		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.LET:
		return p.parseLetStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	state := &ast.LetStatement{Token: p.currentToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	state.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	//todo: skipping expressions for now
	for !p.currentTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return state
}

func (p *Parser) currentTokenIs(t token.TokenType) bool {
	return p.currentToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// note that this method has a side effect of advancing to the next token
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		return false
	}
}
