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
		rf := rv.Type().Field(i)

		tag := rf.Tag.Get(`icl`)
		if tag == "" {
			continue
		}

		n, err := e.buildNode(tag, rf, value)
		if err != nil {
			return nil, err
		}

		a.Nodes = append(a.Nodes, n)
	}

	return a, nil
}

func (e encoder) buildNode(tag string, rf reflect.StructField, rv reflect.Value) (Node, error) {
	rk := rf.Type.Kind()
	if rk == reflect.Pointer {
		if rv.IsNil() {
			return &AssignNode{
				Name:  &Identifier{Token: Token{Type: TknIdent, Literal: tag}, Value: tag},
				Value: &NullNode{},
			}, nil
		}
		rk = rf.Type.Elem().Kind()
		rv = rv.Elem()
	}

	switch rk {
	// primatives
	case reflect.String,
		reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		if tag == ".param" {
			return nil, errors.New("root struct cannot contain params")
		}
		v, err := e.buildPrimativeNode(tag, rk, rv)
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

		for i := 0; i < rv.Len(); i++ {
			node, err := e.buildPrimativeNode(tag, rv.Type().Elem().Kind(), rv.Index(i))
			if err != nil {
				return nil, err
			}
			elems = append(elems, node)
		}

		return &AssignNode{
			Name:  &Identifier{Token: Token{Type: TknIdent, Literal: tag}, Value: tag},
			Value: &SliceNode{Elements: elems},
		}, nil

	case reflect.Struct:
		var (
			params []Token
			body   []Node
		)

		for i := 0; i < rv.NumField(); i++ {
			value := rv.Field(i)
			field := rv.Type().Field(i)

			tag := field.Tag.Get(`icl`)
			if tag == "" {
				continue
			}

			if tag == ".param" {
				if field.Type.Kind() != reflect.String {
					return nil, errors.New("block params can only be of type string ")
				}
				params = append(params, Token{
					Literal: value.Interface().(string),
				})
				continue
			}

			n, err := e.buildNode(tag, field, value)
			if err != nil {
				return nil, err
			}

			body = append(body, n)
		}

		return &BlockNode{
			Token:      Token{Literal: tag},
			Parameters: params,
			Body: &BlockBodyNode{
				Nodes: body,
			},
		}, nil

	case reflect.Map:
		elems := make(map[Node]Node)

		for _, key := range rv.MapKeys() {
			n, err := e.buildPrimativeNode(tag, rv.Type().Elem().Kind(), rv.MapIndex(key))
			if err != nil {
				return nil, err
			}

			elems[&StringNode{Value: key.Interface().(string)}] = n
		}

		return &AssignNode{
			Name: &Identifier{Token: Token{Type: TknIdent, Literal: tag}, Value: tag},
			Value: &MapNode{
				Elements: elems,
			},
		}, nil
	}

	return nil, errors.New("cant convert " + rk.String())
}

func (e encoder) buildPrimativeNode(tag string, rk reflect.Kind, rv reflect.Value) (Node, error) {
	switch rk {
	case reflect.String:
		return &StringNode{Value: rv.Interface().(string)}, nil
	case reflect.Bool:
		return &BooleanNode{Value: rv.Interface().(bool)}, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return &NumberNode{Value: strconv.FormatInt(rv.Int(), 10)}, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return &NumberNode{Value: strconv.FormatUint(rv.Uint(), 10)}, nil
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

		return &NumberNode{Value: strconv.FormatFloat(rv.Float(), 'f', precision, 64)}, nil
	}

	return nil, errors.New("invalid kind " + rk.String())
}
