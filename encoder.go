package icl

import (
	"errors"
	"reflect"
	"strconv"
)

// Encoder handles the transation of a strinc into an Ast
type Encoder struct {
	ast *Ast
	rv  reflect.Value
}

// NewEncoder creates a new instance of the Encoder struct used to transalate a go struct into an icl Ast
func NewEncoder(v any) (*Encoder, error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Struct || (rv.Kind() == reflect.Pointer && rv.Elem().Kind() != reflect.Struct) {
		return nil, errors.New("can only encode struct and *struct values")
	}

	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}

	return &Encoder{&Ast{}, rv}, nil
}

// Encode runs the encoder logic and returns the resulting Ast
func (e Encoder) Encode(v any) (*Ast, error) {
	for i := 0; i < e.rv.NumField(); i++ {
		value := e.rv.Field(i)
		rf := e.rv.Type().Field(i)

		tagString := rf.Tag.Get(`icl`)
		if tagString == "" {
			continue
		}

		tag, err := parseTags(tagString)
		if err != nil {
			return nil, err
		}

		n, err := e.buildNode(tag, rf, value)
		if err != nil {
			return nil, err
		}

		e.ast.Nodes = append(e.ast.Nodes, n)
	}

	return e.ast, nil
}

func (e Encoder) buildNode(tag *tags, rf reflect.StructField, rv reflect.Value) (Node, error) {
	rk := rf.Type.Kind()
	if rk == reflect.Pointer {
		if rv.IsNil() {
			return &AssignNode{
				Name:  &Identifier{Token: Token{Type: TknIdent, Literal: tag.key}, Value: tag.key},
				Value: &NullNode{},
			}, nil
		}
		rk = rf.Type.Elem().Kind()
		rv = rv.Elem()
	}

	switch rk {
	// primitives
	case reflect.String,
		reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:

		if tag.isParam {
			return nil, errors.New("root struct cannot contain params")
		}

		v, err := e.buildPrimitiveNode(tag, rk, rv)
		if err != nil {
			return nil, err
		}

		return &AssignNode{
			Name:  &Identifier{Token: Token{Type: TknIdent, Literal: tag.key}, Value: tag.key},
			Value: v,
		}, nil

	// complex types
	case reflect.Slice:
		if tag.env != "" {
			return nil, errors.New("env() macro not allowed on slice field")
		}

		var elems []Node

		for i := 0; i < rv.Len(); i++ {
			var (
				node Node
				err  error
			)

			if rv.Type().Elem().Kind() == reflect.Struct {
				node, err = e.buildStructNode(tag, rv.Index(i))
			} else {
				node, err = e.buildPrimitiveNode(tag, rv.Type().Elem().Kind(), rv.Index(i))
			}
			if err != nil {
				return nil, err
			}
			elems = append(elems, node)
		}

		if rv.Type().Elem().Kind() == reflect.Struct {
			return &CollectionNode{
				Elements: elems,
			}, nil
		}

		return &AssignNode{
			Name:  &Identifier{Token: Token{Type: TknIdent, Literal: tag.key}, Value: tag.key},
			Value: &SliceNode{Elements: elems},
		}, nil

	case reflect.Struct:
		if tag.env != "" {
			return nil, errors.New("env() macro not allowed on struct field")
		}

		return e.buildStructNode(tag, rv)

	case reflect.Map:
		if tag.env != "" {
			return nil, errors.New("env() macro not allowed on map field")
		}

		elems := make(map[Node]Node)

		for _, key := range rv.MapKeys() {
			n, err := e.buildPrimitiveNode(tag, rv.Type().Elem().Kind(), rv.MapIndex(key))
			if err != nil {
				return nil, err
			}

			elems[&StringNode{Value: key.Interface().(string)}] = n
		}

		return &AssignNode{
			Name: &Identifier{Token: Token{Type: TknIdent, Literal: tag.key}, Value: tag.key},
			Value: &MapNode{
				Elements: elems,
			},
		}, nil
	}

	return nil, errors.New("cant convert " + rk.String())
}

func (e Encoder) buildPrimitiveNode(tag *tags, rk reflect.Kind, rv reflect.Value) (Node, error) {
	if tag.env != "" {
		return &EnvarNode{Identifier: &Identifier{Value: tag.env}}, nil
	}
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
		return &NumberNode{Value: strconv.FormatFloat(rv.Float(), 'f', tag.precision, 64)}, nil
	}

	return nil, errors.New("invalid kind " + rk.String())
}

func (e Encoder) buildStructNode(tag *tags, rv reflect.Value) (Node, error) {
	var (
		params []Token
		body   []Node
	)

	for i := 0; i < rv.NumField(); i++ {
		value := rv.Field(i)
		field := rv.Type().Field(i)

		tagString := field.Tag.Get(`icl`)
		if tagString == "" {
			continue
		}

		ftag, err := parseTags(tagString)
		if err != nil {
			return nil, err
		}

		if ftag.isParam {
			if field.Type.Kind() != reflect.String {
				return nil, errors.New("block params can only be of type string ")
			}
			params = append(params, Token{
				Literal: value.Interface().(string),
			})
			continue
		}

		n, err := e.buildNode(ftag, field, value)
		if err != nil {
			return nil, err
		}

		body = append(body, n)
	}

	return &BlockNode{
		Token:      Token{Literal: tag.key},
		Parameters: params,
		Body: &BlockBodyNode{
			Nodes: body,
		},
	}, nil
}
