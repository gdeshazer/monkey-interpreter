package lexer

import (
	"monkeyInterpreter/pkg/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func New(input string) *Lexer {
	var lexer = &Lexer{input: input}

	lexer.readChar()

	return lexer
}

func (lexer *Lexer) readChar() {
	if lexer.readPosition >= len(lexer.input) {
		// set current character to ASCII code 0 (NUL) when we are at the limit of the input length
		lexer.ch = 0
	} else {
		lexer.ch = lexer.input[lexer.readPosition]
	}

	lexer.position = lexer.readPosition
	lexer.readPosition += 1
}

func (lexer *Lexer) NextToken() token.Token {
	var tok token.Token

	lexer.skipWhitespace()

	switch lexer.ch {
	case '=':
		tok = token.NewToken(token.ASSIGN, lexer.ch)
	case ';':
		tok = token.NewToken(token.SEMICOLON, lexer.ch)
	case '(':
		tok = token.NewToken(token.LPAREN, lexer.ch)
	case ')':
		tok = token.NewToken(token.RPAREN, lexer.ch)
	case ',':
		tok = token.NewToken(token.COMMA, lexer.ch)
	case '+':
		tok = token.NewToken(token.PLUS, lexer.ch)
	case '{':
		tok = token.NewToken(token.LBRACE, lexer.ch)
	case '}':
		tok = token.NewToken(token.RBRACE, lexer.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(lexer.ch) {
			tok.Literal = lexer.readIdentifier()
			tok.Type = token.LookUpIdent(tok.Literal)
			return tok
		} else if isInteger(lexer.ch) {
			tok.Type = token.INT
			tok.Literal = lexer.readNumber()
			return tok
		} else {
			tok = token.NewToken(token.ILLEGAL, lexer.ch)
		}
	}

	lexer.readChar()

	return tok
}

func (lexer *Lexer) skipWhitespace() {
	for lexer.ch == ' ' || lexer.ch == '\t' || lexer.ch == '\n' || lexer.ch == '\r' {
		lexer.readChar()
	}
}

func (lexer *Lexer) readIdentifier() string {
	var start = lexer.position
	for isLetter(lexer.ch) {
		lexer.readChar()
	}

	return lexer.input[start:lexer.position]
}

func (lexer *Lexer) readNumber() string {
	var start = lexer.position

	for isInteger(lexer.ch) {
		lexer.readChar()
	}

	return lexer.input[start:lexer.position]
}

func isLetter(ch byte) bool {
	var isLowerCase = 'a' <= ch && ch <= 'z'
	var isUpperCase = 'A' <= ch && ch <= 'Z'
	var isUnderScore = ch == '_'

	return isLowerCase || isUpperCase || isUnderScore
}

// only handling basic integer types to simplify things
func isInteger(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
