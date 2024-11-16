package icl

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

type encoder struct {
}

func NewEncoder(v any) (*Ast, error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Struct || (rv.Kind() == reflect.Pointer && rv.Elem().Kind() != reflect.Struct) {
		return nil, errors.New("can only encode struct and *struct values")
	}

	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}

	a := &Ast{}
	e := encoder{}

	for i := 0; i < rv.NumField(); i++ {
		value := rv.Field(i)
		field := rv.Type().Field(i)

		tag := field.Tag.Get(`icl`)
		if tag == "" {
			continue
		}

		n, err := e.buildNode(tag, field, value)
		if err != nil {
			return nil, err
		}

		a.Nodes = append(a.Nodes, n)
	}

	return a, nil
}

func (e encoder) buildNode(tag string, field reflect.StructField, value reflect.Value) (Node, error) {
	kind := field.Type.Kind()
	if kind == reflect.Pointer {
		if value.IsNil() {
			return &AssignNode{
				Name:  &Identifier{Token: Token{Type: TknIdent, Literal: tag}, Value: tag},
				Value: &NullNode{},
			}, nil
		}
		kind = field.Type.Elem().Kind()
		value = value.Elem()
	}

	switch kind {
	// primatives
	case reflect.String,
		reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		v, err := e.buildPrimativeNode(tag, kind, value)
		if err != nil {
			return nil, err
		}
		return &AssignNode{
			Name:  &Identifier{Token: Token{Type: TknIdent, Literal: tag}, Value: tag},
			Value: v,
		}, nil

	// complex types
	case reflect.Slice:
		var elems []Node

		for i := 0; i < value.Len(); i++ {
			node, err := e.buildPrimativeNode(tag, value.Type().Elem().Kind(), value.Index(i))
			if err != nil {
				return nil, err
			}
			elems = append(elems, node)
		}

		return &AssignNode{
			Name:  &Identifier{Token: Token{Type: TknIdent, Literal: tag}, Value: tag},
			Value: &SliceNode{Elements: elems},
		}, nil
	case reflect.Map:
		panic("not implemented")
	case reflect.Struct:
		panic("not implemented")
	}

	return nil, errors.New("cant convert " + kind.String())
}

func (e encoder) buildPrimativeNode(tag string, kind reflect.Kind, value reflect.Value) (Node, error) {
	switch kind {
	case reflect.String:
		return &StringNode{Value: value.Interface().(string)}, nil
	case reflect.Bool:
		return &BooleanNode{Value: value.Interface().(bool)}, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return &NumberNode{Value: strconv.FormatInt(value.Int(), 10)}, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return &NumberNode{Value: strconv.FormatUint(value.Uint(), 10)}, nil
	case reflect.Float32, reflect.Float64:
		precision := -1
		if pos := strings.LastIndex(tag, "."); pos != -1 {
			i, err := strconv.Atoi(tag[pos+1:])
			if err != nil || i < 1 {
				return nil, errors.New("invalid tag " + tag)
			}
			precision = i
			tag = tag[:pos]
		}

		return &NumberNode{Value: strconv.FormatFloat(value.Float(), 'f', precision, 64)}, nil
	}

	return nil, errors.New("invalid kind " + kind.String())
}
