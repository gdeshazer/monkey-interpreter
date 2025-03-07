package ast

import "bytes"

type Program struct {
	Statements []Statement
}

func (prog *Program) TokenLiteral() string {
	if len(prog.Statements) > 0 {
		return prog.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (prog *Program) String() string {
	var out bytes.Buffer

	for _, s := range prog.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}
