package icl

import (
	"errors"
	"os"
	"reflect"
	"strconv"
)

var errFieldNotFound = errors.New("field not found")

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

type decoder struct {
	ast          Ast
	target       reflect.Value
	paramCounter int
}

func (d *decoder) decode() error {
	for _, node := range d.ast.Nodes {
		var err error

		switch n := node.(type) {
		case *AssignNode:
			err = d.assign(n, d.target)
		case *BlockNode:
			v, _, err := d.findTargetField(&Identifier{Value: n.Token.Literal}, d.target)
			if err != nil {
				if errors.Is(errFieldNotFound, err) {
					return nil
				}
				return err
			}

			err = d.block(n, *v)
		default:
			// NB to keep icl as fault tollerent as possible any other node types in the ast will be ignored
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (d *decoder) assign(node *AssignNode, target reflect.Value) error {
	v, f, err := d.findTargetField(node.Name, target)
	if err != nil {
		if errors.Is(errFieldNotFound, err) {
			return nil
		}
		return err
	}

	rv, rf := *v, *f

	rk := rf.Type.Kind()
	if rk == reflect.Pointer {
		if _, ok := node.Value.(*NullNode); ok {
			return nil
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

		d.assignPrimitiveNode(node.Value, rv, rv.Kind(), false)

	// complex types
	case reflect.Slice:
		val, ok := node.Value.(*SliceNode)
		if !ok {
			return errors.New("node is not a slice")
		}
		if rv.Kind() != reflect.Slice {
			return errors.New("invalid node " + rv.Kind().String())
		}

		for _, entry := range val.Elements {
			if err := d.assignPrimitiveNode(entry, rv, rv.Type().Elem().Kind(), true); err != nil {
				return nil
			}
		}

	case reflect.Struct:
		panic("not implemented")
	}

	return nil
}

func (d *decoder) block(node *BlockNode, target reflect.Value) error {
	pc := 0
	d.paramCounter = pc

	// params
	for _, param := range node.Parameters {
		v, _, err := d.findTargetField(&Identifier{Value: ".param"}, target)
		if err != nil {
			if errors.Is(errFieldNotFound, err) {
				continue
			}
			return err
		}

		rv := *v
		pc++
		d.paramCounter = pc

		if rv.Kind() != reflect.String {
			return errors.New(".param fields must be a string")
		}

		rv.SetString(param.Literal)
	}

	// body
	for _, node := range node.Body.Nodes {
		var err error

		switch n := node.(type) {
		case *AssignNode:
			err = d.assign(n, target)
		case *BlockNode:
			v, _, err := d.findTargetField(&Identifier{Value: n.Token.Literal}, target)
			if err != nil {
				if errors.Is(errFieldNotFound, err) {
					return nil
				}
				return err
			}

			err = d.block(n, *v)
		default:
			// NB to keep icl as fault tollerent as possible any other node types in the ast will be ignored
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (d *decoder) findTargetField(ident *Identifier, target reflect.Value) (*reflect.Value, *reflect.StructField, error) {
	var paramCounter int
	for i := 0; i < target.NumField(); i++ {
		// TODO: this needs a cache so i don't have to keep parsing the struct fields each time

		rv := target.Field(i)
		rf := target.Type().Field(i)

		tagString := rf.Tag.Get(`icl`)
		if tagString == "" {
			continue
		}

		tag, err := parseTags(tagString)
		if err != nil {
			return nil, nil, err
		}

		if tag.key != ident.Value && ident.Value != ".param" {
			continue
		}

		if tag.isParam && paramCounter < d.paramCounter {
			paramCounter++
			continue
		}

		return &rv, &rf, nil
	}

	return nil, nil, errFieldNotFound
}

func (d *decoder) assignPrimitiveNode(node Node, rv reflect.Value, rk reflect.Kind, isSlice bool) error {
	switch v := node.(type) {
	case *EnvarNode:
		val := os.Getenv(v.Identifier.Value)
		switch rk {
		case reflect.String:
			rv.SetString(val)
		case reflect.Bool:
			rv.SetBool(val == "true")
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			val, err := parseIntKind(val, rk, false)
			if err != nil {
				return err
			}
			rv.SetInt(val.(int64))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			val, err := parseUintKind(val, rk, false)
			if err != nil {
				return err
			}
			rv.SetUint(val.(uint64))
		case reflect.Float32, reflect.Float64:
			bs := 32
			if rk == reflect.Float64 {
				bs = 64
			}
			val, err := strconv.ParseFloat(val, bs)
			if err != nil {
				return err
			}
			rv.SetFloat(val)
		}
	case *StringNode:
		if rk != reflect.String {
			return errors.New("invalid type " + rk.String())
		}
		if isSlice {
			rv.Set(reflect.Append(rv, reflect.ValueOf(v.Value)))
		} else {
			rv.SetString(v.Value)
		}
	case *BooleanNode:
		if rk != reflect.Bool {
			return errors.New("invalid type " + rk.String())
		}
		if isSlice {
			rv.Set(reflect.Append(rv, reflect.ValueOf(v.Value)))
		} else {
			rv.SetBool(v.Value)
		}
	case *NumberNode:
		switch rk {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			val, err := parseIntKind(v.Value, rk, isSlice)
			if err != nil {
				return err
			}
			if isSlice {
				rv.Set(reflect.Append(rv, reflect.ValueOf(val)))
			} else {
				rv.SetInt(val.(int64))
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			val, err := parseUintKind(v.Value, rk, isSlice)
			if err != nil {
				return err
			}
			if isSlice {
				rv.Set(reflect.Append(rv, reflect.ValueOf(val)))
			} else {
				rv.SetUint(val.(uint64))
			}
		case reflect.Float32, reflect.Float64:
			bs := 32
			if rk == reflect.Float64 {
				bs = 64
			}
			val, err := strconv.ParseFloat(v.Value, bs)
			if err != nil {
				return err
			}
			if isSlice {
				rv.Set(reflect.Append(rv, reflect.ValueOf(val)))
			} else {
				rv.SetFloat(val)
			}
		default:
			return errors.New("invalid type " + rk.String())
		}
	}

	return nil
}

func parseIntKind(s string, k reflect.Kind, typed bool) (any, error) {
	var i int64
	var err error
	switch k {
	case reflect.Int8:
		i, err = strconv.ParseInt(s, 10, 8)
		if typed {
			return int8(i), err
		}
	case reflect.Int16:
		i, err = strconv.ParseInt(s, 10, 16)
		if typed {
			return int16(i), err
		}
	case reflect.Int32:
		i, err = strconv.ParseInt(s, 10, 32)
		if typed {
			return int32(i), err
		}
	case reflect.Int64:
		return strconv.ParseInt(s, 10, 64)
	case reflect.Int:
		i, err = strconv.ParseInt(s, 10, 64)
		if typed {
			return int(i), err
		}
	default:
		return 0, errors.New("unreachable")
	}

	return i, err
}

func parseUintKind(s string, k reflect.Kind, typed bool) (any, error) {
	var i uint64
	var err error
	switch k {
	case reflect.Uint8:
		i, err = strconv.ParseUint(s, 10, 8)
		if typed {
			return uint8(i), err
		}
	case reflect.Uint16:
		i, err = strconv.ParseUint(s, 10, 16)
		if typed {
			return uint16(i), err
		}
	case reflect.Uint32:
		i, err = strconv.ParseUint(s, 10, 32)
		if typed {
			return uint32(i), err
		}
	case reflect.Uint64:
		return strconv.ParseUint(s, 10, 64)
	case reflect.Uint:
		i, err = strconv.ParseUint(s, 10, 64)
		if typed {
			return uint(i), err
		}
	default:
		return 0, errors.New("unreachable")
	}

	return i, err
}
