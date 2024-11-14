package icl

import (
	"bytes"
	"fmt"
)

type ExpressionStatement struct {
	Token      Token
	Expression Expression
}

// String implements Statement
func (n *ExpressionStatement) String() string {
	if n.Expression == nil {
		return ""
	}

	return n.Expression.String()
}

// TokenLiteral implements Statement
func (n *ExpressionStatement) TokenLiteral() string {
	return n.Token.Literal
}

// statementNode implements Statement
func (*ExpressionStatement) statementNode() {
}

var _ Statement = (*ExpressionStatement)(nil)

type AssignStatement struct {
	Token Token
	Name  *Identifier
	Value Expression
}

// String implements Statement
func (n *AssignStatement) String() string {
	var buf bytes.Buffer

	buf.WriteString(n.Name.Value)
	buf.WriteString(" = ")

	if n.Value != nil {
		buf.WriteString(n.Value.String())
	}

	return buf.String() + "\n"
}

// statementNode implements Statement
func (n *AssignStatement) statementNode() {
}

// TokenLiteral implements Node
func (n *AssignStatement) TokenLiteral() string {
	return n.Token.Literal
}

var _ Statement = (*AssignStatement)(nil)

type BlockStatement struct {
	Token      Token
	Parameters []Token
	Body       *BlockBodyStatement
}

// String implements Expression
func (n *BlockStatement) String() string {
	var buf bytes.Buffer

	buf.WriteString(n.TokenLiteral())
	for _, p := range n.Parameters {
		buf.WriteString(" " + p.Literal)
	}
	buf.WriteString(" ")

	buf.WriteString(n.Body.String())

	return buf.String()
}

// TokenLiteral implements Expression
func (n *BlockStatement) TokenLiteral() string {
	return n.Token.Literal
}

// expressionNode implements Expression
func (*BlockStatement) statementNode() {
}

var _ Statement = (*BlockStatement)(nil)

type BlockBodyStatement struct {
	Token      Token
	Statements []Statement
}

// String implements Statement
func (n *BlockBodyStatement) String() string {
	var buf bytes.Buffer

	buf.WriteString("{\n")

	for _, stmt := range n.Statements {
		buf.WriteString(fmt.Sprintf("    %s\n", stmt.String()))
	}

	buf.WriteString("}")

	return buf.String()
}

// statementNode implements Statement
func (n *BlockBodyStatement) statementNode() {
}

// TokenLiteral implements Node
func (n *BlockBodyStatement) TokenLiteral() string {
	return n.Token.Literal
}

var _ Statement = (*BlockBodyStatement)(nil)
