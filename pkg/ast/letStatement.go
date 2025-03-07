package ast

import "monkeyInterpreter/pkg/token"

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	value Expression
}

func (ls *LetStatement) statementNode() {}

func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}
