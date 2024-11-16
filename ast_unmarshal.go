package icl

import (
	"fmt"
	"reflect"
)

type InvalidUnmarshalError struct {
	Type reflect.Type
}

func (e *InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "icl: Unmarshal(nil)"
	}

	if e.Type.Kind() != reflect.Pointer {
		return "icl: Unmarshal(non-pointer " + e.Type.String() + ")"
	}
	return "icl: Unmarshal(nil " + e.Type.String() + ")"
}

func (a Ast) Unmarshal(v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() || rv.Elem().Kind() != reflect.Struct {
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}

	d := decoder{a, rv.Elem()}
	return d.decode()
}

type decoder struct {
	ast    Ast
	target reflect.Value
}

func (d decoder) decode() error {
	for _, node := range d.ast.Nodes {
		var err error

		switch n := node.(type) {
		case *AssignNode:
			err = d.assign(n, d.target)
		default:
			err = fmt.Errorf("invalid node type: %T", n)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (d decoder) assign(node *AssignNode, target reflect.Value) error {
	panic("not implemented")
}
