package lexer

import (
	"monkeyInterpreter/pkg/token"
)

type Lexer struct {

	// input represents the source code that the Lexer will parse into tokens.
	input string

	// position tracks the current reading position in the input string.
	position int

	// readPosition tracks the next read position in the input string
	readPosition int

	// ch is the current character to process
	ch byte
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
		if lexer.peakChar() == '=' {
			var ch = lexer.ch

			// read next char in input
			lexer.readChar()
			var literal = string(ch) + string(lexer.ch)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = token.NewToken(token.ASSIGN, lexer.ch)
		}
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
	case '-':
		tok = token.NewToken(token.MINUS, lexer.ch)
	case '*':
		tok = token.NewToken(token.ASTERISK, lexer.ch)
	case '/':
		tok = token.NewToken(token.SLASH, lexer.ch)
	case '!':
		if lexer.peakChar() == '=' {
			var ch = lexer.ch

			// read next char in input
			lexer.readChar()
			var literal = string(ch) + string(lexer.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			tok = token.NewToken(token.BANG, lexer.ch)
		}
	case '>':
		tok = token.NewToken(token.GREATERTHAN, lexer.ch)
	case '<':
		tok = token.NewToken(token.LESSTHAN, lexer.ch)
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

func (lexer *Lexer) peakChar() byte {
	if lexer.readPosition >= len(lexer.input) {
		return 0
	} else {
		return lexer.input[lexer.readPosition]
	}
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
