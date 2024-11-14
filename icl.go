package icl

func Parse(data string) (*Ast, error) {
	p := NewParser(newLexer(data))

	return p.Parse(), nil
}
