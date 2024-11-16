package icl

import "os"

// Parse an icl string into an Ast
func Parse(data string) (*Ast, error) {
	p := NewParser(newLexer(data))

	return p.Parse(), nil
}

// ParseFile parses the contents of a file into an Ast
func ParseFile(path string) (*Ast, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	p := NewParser(newLexer(string(data)))

	return p.Parse(), nil
}

// MarshalString marshals a struct value into a string
func MarshalString(v any) (string, error) {
	e, err := NewEncoder(v)
	if err != nil {
		return "", err
	}

	a, err := e.Encode(v)
	if err != nil {
		return "", err
	}

	return a.String(), nil
}

// Marshal marshals a strict value into a byte array
func Marshal(v any) ([]byte, error) {
	e, err := NewEncoder(v)
	if err != nil {
		return nil, err
	}

	a, err := e.Encode(v)
	if err != nil {
		return nil, err
	}

	return a.Bytes(), nil
}

// MarshalFile marshals a strict value directly into a file
func MarshalFile(v any, path string) error {
	e, err := NewEncoder(v)
	if err != nil {
		return err
	}

	a, err := e.Encode(v)
	if err != nil {
		return err
	}

	return os.WriteFile(path, a.Bytes(), 0644)
}
