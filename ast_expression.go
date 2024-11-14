package icl

import (
	"bytes"
	"fmt"
)

type InfixExpression struct {
	Token    Token
	Left     Expression
	Operator string
	Right    Expression
}

// String implements Expression
func (n *InfixExpression) String() string {
	var buf bytes.Buffer

	buf.WriteString("(")
	buf.WriteString(n.Left.String())
	buf.WriteString(" " + n.Operator + " ")
	buf.WriteString(n.Right.String())
	buf.WriteString(")")

	return buf.String()
}

// TokenLiteral implements Expression
func (n *InfixExpression) TokenLiteral() string {
	return n.Token.Literal
}

// expressionNode implements Expression
func (*InfixExpression) expressionNode() {
}

var _ Expression = (*InfixExpression)(nil)

type IntegerLiteral struct {
	Token Token
	Value int64
}

// String implements Expression
func (n *IntegerLiteral) String() string {
	return n.Token.Literal
}

// TokenLiteral implements Expression
func (n *IntegerLiteral) TokenLiteral() string {
	return n.Token.Literal
}

// expressionNode implements Expression
func (*IntegerLiteral) expressionNode() {
}

var _ Expression = (*IntegerLiteral)(nil)

type StringLiteral struct {
	Token Token
	Value string
}

// String implements Expression
func (n *StringLiteral) String() string {
	return `"` + n.Token.Literal + `"`
}

// TokenLiteral implements Expression
func (n *StringLiteral) TokenLiteral() string {
	return n.Token.Literal
}

// expressionNode implements Expression
func (*StringLiteral) expressionNode() {
}

var _ Expression = (*StringLiteral)(nil)

type PrefixExpression struct {
	Token    Token
	Operator string
	Right    Expression
}

// String implements Expression
func (n *PrefixExpression) String() string {
	var buf bytes.Buffer

	buf.WriteString("(")
	buf.WriteString(n.Operator)
	buf.WriteString(n.Right.String())
	buf.WriteString(")")

	return buf.String()
}

// TokenLiteral implements Expression
func (n *PrefixExpression) TokenLiteral() string {
	return n.Token.Literal
}

// expressionNode implements Expression
func (*PrefixExpression) expressionNode() {
}

var _ Expression = (*PrefixExpression)(nil)

type BooleanLiteral struct {
	Token Token
	Value bool
}

// String implements Expression
func (n *BooleanLiteral) String() string {
	return n.Token.Literal
}

// TokenLiteral implements Expression
func (n *BooleanLiteral) TokenLiteral() string {
	return n.Token.Literal
}

// expressionNode implements Expression
func (*BooleanLiteral) expressionNode() {
}

var _ Expression = (*BooleanLiteral)(nil)

type NullLiteral struct {
	Token Token
}

// String implements Expression
func (n *NullLiteral) String() string {
	return "NULL"
}

// TokenLiteral implements Expression
func (n *NullLiteral) TokenLiteral() string {
	return n.Token.Literal
}

// expressionNode implements Expression
func (*NullLiteral) expressionNode() {
}

var _ Expression = (*NullLiteral)(nil)

type ArrayLiteral struct {
	Token    Token
	Elements []Expression
}

// String implements Expression
func (a *ArrayLiteral) String() string {
	var buf bytes.Buffer
	buf.WriteString("[")

	for i, elem := range a.Elements {
		if i > 0 {
			buf.WriteString(", ")
		}

		buf.WriteString(elem.String())
	}

	buf.WriteString("]")

	return buf.String()
}

// TokenLiteral implements Expression
func (a *ArrayLiteral) TokenLiteral() string {
	return a.Token.Literal
}

// expressionNode implements Expression
func (*ArrayLiteral) expressionNode() {
}

var _ Expression = (*ArrayLiteral)(nil)

type MapLiteral struct {
	Token    Token
	Elements map[Expression]Expression
}

// String implements Expression
func (n *MapLiteral) String() string {
	var buf bytes.Buffer

	buf.WriteString("{\n")

	for key, val := range n.Elements {
		buf.WriteString(fmt.Sprintf("    %s: %s,\n", key.String(), val.String()))
	}

	buf.WriteString("}")

	return buf.String()
}

// TokenLiteral implements Expression
func (n *MapLiteral) TokenLiteral() string {
	return n.Token.Literal
}

// expressionNode implements Expression
func (*MapLiteral) expressionNode() {
}

var _ Expression = (*MapLiteral)(nil)

type BlockLiteral struct {
	Token      Token
	Parameters []Token
	Body       *BlockStatement
}

// String implements Expression
func (n *BlockLiteral) String() string {
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
func (n *BlockLiteral) TokenLiteral() string {
	return n.Token.Literal
}

// expressionNode implements Expression
func (*BlockLiteral) expressionNode() {
}

var _ Expression = (*BlockLiteral)(nil)
