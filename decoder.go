package icl

import (
	"errors"
	"fmt"
	"log"
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

type Decoder struct {
	ast          Ast
	target       reflect.Value
	paramCounter int
	blockMap     map[reflect.Value]map[string]struct{}
	recover      error
	line         int
	pos          int
}

func NewDecoder(a Ast, target reflect.Value) *Decoder {
	return &Decoder{
		ast:      a,
		target:   target,
		blockMap: make(map[reflect.Value]map[string]struct{}),
	}
}

func (d *Decoder) decode() error {
	for _, node := range d.ast.Nodes {
		if err := d.node(node, d.target, ""); err != nil {
			return err
		}
	}

	return nil
}

func (d *Decoder) assign(node *AssignNode, target reflect.Value, path string) error {
	d.line = node.Token.Line
	d.line = node.Token.Pos

	v, f, tag, err := d.findTargetField(node.Name, target)
	if err != nil {
		if errors.Is(errFieldNotFound, err) {
			return nil
		}
		return err
	}

	rv, rf := *v, *f
	path += "." + tag.key

	defer func() {
		// if err := recover(); err != nil {
		// 	d.recover = fmt.Errorf("%s: %v", path, err)
		// }
	}()

	var setErr error
	var originalRv reflect.Value
	rk := rf.Type.Kind()
	if rk == reflect.Pointer {
		if _, ok := node.Value.(*NullNode); ok {
			return nil
		}

		originalRv = rv
		rk = rf.Type.Elem().Kind()
		if rv.IsZero() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}
		rv = rv.Elem()

		defer func() {
			if setErr != nil && rv.IsZero() {
				originalRv.Set(reflect.Zero(originalRv.Type()))
			}
		}()
	}

	switch rk {
	// primatives
	case reflect.String,
		reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:

		d.line = node.Value.Tkn().Line
		d.line = node.Value.Tkn().Pos

		setErr = d.assignPrimitiveNode(node.Value, rv, rk, false)

	// complex types
	case reflect.Slice:
		val, ok := node.Value.(*SliceNode)
		if !ok {
			return errors.New("node is not a slice")
		}

		for _, entry := range val.Elements {
			d.line = entry.Tkn().Line
			d.line = entry.Tkn().Pos

			if err := d.assignPrimitiveNode(entry, rv, rv.Type().Elem().Kind(), true); err != nil {
				setErr = err
				return nil
			}
		}

	case reflect.Map:
		val, ok := node.Value.(*MapNode)
		if !ok {
			return errors.New("node is not a map")
		}

		for key, value := range val.Elements {
			d.line = value.Tkn().Line
			d.line = value.Tkn().Pos

			t := reflect.New(rv.Type().Elem())
			log.Print(rv.Type().Elem(), t)
			rv.SetMapIndex(reflect.ValueOf(key.(*StringNode).Value), t)

			var isSlice bool
			rt := rv.Type()

			if rt.Kind() == reflect.Slice {
				rt = rt.Elem()
				isSlice = true
			}

			if err := d.assignPrimitiveNode(value, t, rt.Elem().Kind(), isSlice); err != nil {
				return err
			}

			// rv.SetMapIndex(reflect.ValueOf(key.(*StringNode).Value), t)
		}

	default:
		setErr = errors.New("unknown")
	}

	return nil
}

func (d *Decoder) block(node *BlockNode, rv reflect.Value, path string) error {
	pc := 0
	d.paramCounter = pc

	// track block assignments to make sure we aren't trying to re assign
	if rv.Kind() != reflect.Slice {
		if _, ok := d.blockMap[rv][node.TokenLiteral()]; ok {
			return fmt.Errorf("multiple \"%s\" blocks found for field that is not a slice", node.TokenLiteral())
		}

		if _, ok := d.blockMap[rv]; !ok {
			d.blockMap[rv] = make(map[string]struct{})
		}
		d.blockMap[rv][node.TokenLiteral()] = struct{}{}
	}

	var originalTarget reflect.Value
	if rv.Kind() == reflect.Slice {
		newTgt := reflect.New(rv.Type().Elem()).Elem()
		originalTarget = rv
		rv = newTgt
	} else if rv.Kind() == reflect.Pointer && rv.IsNil() {
		rv.Set(reflect.New(rv.Type().Elem()))
		rv = rv.Elem()
	}

	// params
	for _, param := range node.Parameters {
		d.line = param.Line
		d.line = param.Pos

		v, _, _, err := d.findTargetField(&Identifier{Value: ".param"}, rv)
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
			return errors.New(path + ": .param fields must be a string")
		}

		rv.SetString(param.Literal)
	}

	// body
	for _, node := range node.Body.Nodes {
		if err := d.node(node, rv, path); err != nil {
			return err
		}
	}

	if originalTarget.Kind() == reflect.Slice {
		originalTarget.Set(reflect.Append(originalTarget, rv))
	}

	return nil
}

func (d *Decoder) node(node Node, target reflect.Value, path string) error {
	var err error

	switch n := node.(type) {
	case *AssignNode:
		err = d.assign(n, target, path)
	case *BlockNode:
		v, _, tag, err := d.findTargetField(&Identifier{Value: n.Token.Literal}, target)
		if err != nil {
			if errors.Is(errFieldNotFound, err) {
				return nil
			}
			return err
		}

		err = d.block(n, *v, path+"."+tag.key)
	default:
		// NB to keep icl as fault tollerent as possible any other node types in the ast will be ignored
	}

	if r := d.recover; r != nil {
		d.recover = nil
		err = r
	}

	if err != nil {
		return fmt.Errorf("%w\nline(%d) pos(%d)", err, d.line, d.pos)
	}

	return nil
}

func (d *Decoder) findTargetField(
	ident *Identifier,
	target reflect.Value,
) (
	*reflect.Value,
	*reflect.StructField,
	*tags,
	error,
) {
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
			return nil, nil, nil, err
		}

		if tag.key != ident.Value && ident.Value != ".param" {
			continue
		}

		if tag.isParam && paramCounter < d.paramCounter {
			paramCounter++
			continue
		}

		return &rv, &rf, tag, nil
	}

	return nil, nil, nil, errFieldNotFound
}

func (d *Decoder) assignPrimitiveNode(node Node, rv reflect.Value, rk reflect.Kind, isSlice bool) error {
	d.line = node.Tkn().Line
	d.pos = node.Tkn().Pos

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
			v, err := strconv.ParseFloat(val, bs)
			if err != nil {
				return err
			}
			rv.SetFloat(v)
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
				log.Print(rv, val, rv.CanSet())
				rv.SetFloat(val)
			}
		default:
			return errors.New("invalid type " + rk.String())
		}
	default:
		return errors.New("invalid type " + rk.String())
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
