package icl

import "bytes"

type Ast struct {
	Statements []Statement
}

// String implements Node
func (n *Ast) String() string {
	var buf bytes.Buffer

	for _, stmt := range n.Statements {
		buf.WriteString(stmt.String())
	}

	return buf.String()
}

// TokenLiteral implements Node
func (n *Ast) TokenLiteral() string {
	if len(n.Statements) == 0 {
		return ""
	}

	return n.Statements[0].TokenLiteral()
}

var _ Node = (*Ast)(nil)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Identifier struct {
	Token
	Value string
}

// String implements Statement
func (n *Identifier) String() string {
	return n.Value
}

// expressionNode implements Expression
func (*Identifier) expressionNode() {
}

// TokenLiteral implements Node
func (n *Identifier) TokenLiteral() string {
	return n.Token.Literal
}

var _ Expression = (*Identifier)(nil)
