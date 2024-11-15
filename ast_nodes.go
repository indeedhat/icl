package icl

import (
	"bytes"
	"fmt"
	"strings"
)

type IntegerNode struct {
	Token Token
	Value int64
}

// String implements Node
func (n *IntegerNode) String() string {
	return n.Token.Literal
}

// TokenNode implements Node
func (n *IntegerNode) TokenLiteral() string {
	return n.Token.Literal
}

var _ Node = (*IntegerNode)(nil)

type StringNode struct {
	Token Token
	Value string
}

// String implements Node
func (n *StringNode) String() string {
	return `"` + n.Token.Literal + `"`
}

// TokenNode implements Node
func (n *StringNode) TokenLiteral() string {
	return n.Token.Literal
}

var _ Node = (*StringNode)(nil)

type BooleanNode struct {
	Token Token
	Value bool
}

// String implements Node
func (n *BooleanNode) String() string {
	return n.Token.Literal
}

// TokenNode implements Node
func (n *BooleanNode) TokenLiteral() string {
	return n.Token.Literal
}

var _ Node = (*BooleanNode)(nil)

type NullNode struct {
	Token Token
}

// String implements Node
func (n *NullNode) String() string {
	return "NULL"
}

// TokenLiteral implements Node
func (n *NullNode) TokenLiteral() string {
	return n.Token.Literal
}

var _ Node = (*NullNode)(nil)

type ArrayNode struct {
	Token    Token
	Elements []Node
}

// String implements Node
func (a *ArrayNode) String() string {
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

// TokenLiteral implements Node
func (a *ArrayNode) TokenLiteral() string {
	return a.Token.Literal
}

var _ Node = (*ArrayNode)(nil)

type MapNode struct {
	Token    Token
	Elements map[Node]Node
}

// String implements Node
func (n *MapNode) String() string {
	var buf bytes.Buffer

	buf.WriteString("{\n")

	for key, val := range n.Elements {
		buf.WriteString(indent(fmt.Sprintf("%s: %s,", key.String(), val.String())) + "\n")
	}

	buf.WriteString("}")

	return buf.String()
}

// TokenLiteral implements Node
func (n *MapNode) TokenLiteral() string {
	return n.Token.Literal
}

var _ Node = (*MapNode)(nil)

type AssignNode struct {
	Token Token
	Name  *Identifier
	Value Node
}

// String implements Node
func (n *AssignNode) String() string {
	var buf bytes.Buffer

	buf.WriteString(n.Name.Value)
	buf.WriteString(" = ")

	if n.Value != nil {
		buf.WriteString(n.Value.String())
	}

	return buf.String() + "\n"
}

// statementNode implements Node
func (n *AssignNode) statementNode() {
}

// TokenLiteral implements Node
func (n *AssignNode) TokenLiteral() string {
	return n.Token.Literal
}

var _ Node = (*AssignNode)(nil)

type BlockNode struct {
	Token      Token
	Parameters []Token
	Body       *BlockBodyNode
}

// String implements Node
func (n *BlockNode) String() string {
	var buf bytes.Buffer

	buf.WriteString(n.TokenLiteral())
	for _, p := range n.Parameters {
		buf.WriteString(" " + p.Literal)
	}
	buf.WriteString(" ")

	buf.WriteString(n.Body.String())

	return buf.String()
}

// TokenLiteral implements Node
func (n *BlockNode) TokenLiteral() string {
	return n.Token.Literal
}

var _ Node = (*BlockNode)(nil)

type BlockBodyNode struct {
	Token Token
	Nodes []Node
}

// String implements Node
func (n *BlockBodyNode) String() string {
	var buf bytes.Buffer

	buf.WriteString("{\n")

	for _, stmt := range n.Nodes {
		buf.WriteString(indent(stmt.String()) + "\n")
	}

	buf.WriteString("}")

	return buf.String()
}

// TokenLiteral implements Node
func (n *BlockBodyNode) TokenLiteral() string {
	return n.Token.Literal
}

var _ Node = (*BlockBodyNode)(nil)

func indent(s string) string {
	if strings.Contains(s, "\n") {
		parts := strings.Split(s, "\n")
		s = strings.Join(parts, "\n    ")
	}

	return "    " + s
}
