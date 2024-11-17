package icl

import (
	"fmt"
	"os"
)

// Parsean icl a byte array into an Ast
func Parse(data []byte) (*Ast, error) {
	p := NewParser(newLexer(string(data)))

	return p.Parse(), nil
}

// ParseString an icl string into an Ast
func ParseString(data string) (*Ast, error) {
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

// UnMarshal unmarshals a byte array value into a struct
func UnMarshal(data []byte, v any) error {
	a, err := Parse(data)
	if err != nil {
		return err
	}

	return a.Unmarshal(v)
}

// UnMarshalString unmarshals a string value into a struct
func UnMarshalString(s string, v any) error {
	a, err := ParseString(s)
	if err != nil {
		return err
	}

	return a.Unmarshal(v)
}

// UnMarshalFile unmarshals a file path value directly into a struct
func UnMarshalFile(path string, v any) error {
	a, err := ParseFile(path)
	if err != nil {
		return err
	}

	return a.Unmarshal(v)
}

// UnmarshalVersion takes a map of possible version targets and unmarshels the document int the appropriate one
// If no appropriate target is found then nothing will be unmarshaled
func UnmarshalVersion(data []byte, versions map[int]any) (int, any, error) {
	a, err := Parse(data)
	if err != nil {
		return 0, nil, err
	}

	t, ok := versions[a.Version()]
	if !ok {
		return 0, nil, fmt.Errorf("no target was provided for version %d", a.Version())
	}

	err = a.Unmarshal(t)

	return a.Version(), t, err
}
