package ast

import (
	"bytes"
	"monkeyInterpreter/pkg/token"
)

type InfixExpression struct {
	Token    token.Token
	Right    Expression
	Left     Expression
	Operator string
}

func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

func (ie *InfixExpression) expressionNode() {}
