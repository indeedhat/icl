package icl

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// Node creates the shape for an Ast node
type Node interface {
	TokenLiteral() string
	String() string
}

// Ast contains the Abstract Syntax Tree of an icl ducument
type Ast struct {
	Nodes []Node
}

// Version returns the version of the ICL document contained in the Ast
func (n *Ast) Version() int {
	if len(n.Nodes) == 0 {
		return 0
	}

	assignment, ok := n.Nodes[0].(*AssignNode)
	if !ok || assignment.Name.Value != "version" {
		return 0
	}

	value, ok := assignment.Value.(*NumberNode)
	if !ok {
		return 0
	}

	i, err := strconv.ParseInt(value.Value, 10, 64)
	if err != nil {
		return 0
	}

	return int(i)
}

// Unmarshal fillso out the provided struct pointer with the data in the AST
func (a Ast) Unmarshal(v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() || rv.Elem().Kind() != reflect.Struct {
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}

	d := NewDecoder(a, rv.Elem())
	return d.decode()
}

// String implements Node
func (n *Ast) String() string {
	var buf bytes.Buffer

	for _, stmt := range n.Nodes {
		buf.WriteString(stmt.String() + "\n")
	}

	return buf.String()
}

// Bytes returns the byte array representaiton of the output icl string
func (n *Ast) Bytes() []byte {
	var buf bytes.Buffer

	for _, stmt := range n.Nodes {
		buf.WriteString(stmt.String() + "\n")
	}

	return buf.Bytes()
}

// TokenLiteral implements Node
func (n *Ast) TokenLiteral() string {
	if len(n.Nodes) == 0 {
		return ""
	}

	return n.Nodes[0].TokenLiteral()
}

var _ Node = (*Ast)(nil)

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

type NumberNode struct {
	Token Token
	Value string
}

// String implements Node
func (n *NumberNode) String() string {
	return n.Value
}

// TokenNode implements Node
func (n *NumberNode) TokenLiteral() string {
	return n.Token.Literal
}

var _ Node = (*NumberNode)(nil)

type StringNode struct {
	Token Token
	Value string
}

// String implements Node
func (n *StringNode) String() string {
	return strconv.Quote(n.Value)
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
	if n.Value {
		return "true"
	}

	return "false"
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
	return "null"
}

// TokenLiteral implements Node
func (n *NullNode) TokenLiteral() string {
	return n.Token.Literal
}

var _ Node = (*NullNode)(nil)

type SliceNode struct {
	Token    Token
	Elements []Node
}

// String implements Node
func (n *SliceNode) String() string {
	if len(n.Elements) == 0 {
		return "[]"
	}
	var buf bytes.Buffer
	buf.WriteString("[")

	for i, elem := range n.Elements {
		if i > 0 {
			buf.WriteString(", ")
		}

		buf.WriteString(elem.String())
	}

	buf.WriteString("]")

	return buf.String()
}

// TokenLiteral implements Node
func (n *SliceNode) TokenLiteral() string {
	return n.Token.Literal
}

var _ Node = (*SliceNode)(nil)

type CollectionNode struct {
	Elements []Node
}

// String implements Node
func (n *CollectionNode) String() string {
	if len(n.Elements) == 0 {
		return ""
	}

	var buf bytes.Buffer

	for i, elem := range n.Elements {
		buf.WriteString(elem.String())
		if i < len(n.Elements)-1 {
			buf.WriteString("\n")
		}
	}

	return buf.String()
}

// TokenLiteral implements Node
func (n *CollectionNode) TokenLiteral() string {
	return ""
}

var _ Node = (*CollectionNode)(nil)

type MapNode struct {
	Token    Token
	Elements map[Node]Node
}

// String implements Node
func (n *MapNode) String() string {
	var buf bytes.Buffer

	if len(n.Elements) == 0 {
		return "{}"
	}

	buf.WriteString("{\n")

	keys := make([]Node, 0, len(n.Elements))
	for key := range n.Elements {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].String() < keys[j].String()
	})

	for _, key := range keys {
		buf.WriteString(indent(fmt.Sprintf("%s: %s,", key.String(), n.Elements[key].String())) + "\n")
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

	return buf.String()
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
		buf.WriteString(" " + strconv.Quote(p.Literal))
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

type EnvarNode struct {
	Token      Token
	Identifier *Identifier
}

// String implements Node
func (n *EnvarNode) String() string {
	var buf bytes.Buffer

	buf.WriteString("env(")
	buf.WriteString(n.Identifier.Value)
	buf.WriteString(")")

	return buf.String()
}

// TokenLiteral implements Node
func (n *EnvarNode) TokenLiteral() string {
	return n.Token.Literal
}

var _ Node = (*EnvarNode)(nil)

func indent(s string) string {
	if strings.Contains(s, "\n") {
		parts := strings.Split(s, "\n")
		s = strings.Join(parts, "\n    ")
	}

	return "    " + s
}
