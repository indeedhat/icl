package icl

import "bytes"

type Ast struct {
	Nodes []Node
}

// String implements Node
func (n *Ast) String() string {
	var buf bytes.Buffer

	for _, stmt := range n.Nodes {
		buf.WriteString(stmt.String())
	}

	return buf.String()
}

// TokenLiteral implements Node
func (n *Ast) TokenLiteral() string {
	if len(n.Nodes) == 0 {
		return ""
	}

	return n.Nodes[0].TokenLiteral()
}

var _ Node = (*Ast)(nil)

type Node interface {
	TokenLiteral() string
	String() string
}

type Identifier struct {
	Token
	Value string
}

// String implements Node
func (n *Identifier) String() string {
	return n.Value
}

// TokenLiteral implements Node
func (n *Identifier) TokenLiteral() string {
	return n.Token.Literal
}

var _ Node = (*Identifier)(nil)
