package icl

import (
	"errors"
	"fmt"
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

	v, _, tag, err := d.findTargetField(node.Name, target)
	if err != nil {
		if errors.Is(errFieldNotFound, err) {
			return nil
		}
		return err
	}

	rv := *v
	path += "." + tag.key

	rk := rv.Kind()
	if rk == reflect.Ptr {
		rk = rv.Type().Elem().Kind()
	}

	defer func() {
		if err := recover(); err != nil {
			d.recover = fmt.Errorf("%s: %v", path, err)
		}
	}()

	var setErr error
switcher:
	switch rk {
	// primatives
	case reflect.String,
		reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:

		d.line = node.Value.Tkn().Line
		d.line = node.Value.Tkn().Pos

		setErr = d.assignPrimitiveNode(node.Value, rv, false)

	// complex types
	case reflect.Slice:
		val, ok := node.Value.(*SliceNode)
		if !ok {
			return errors.New("node is not a slice")
		}

		for _, entry := range val.Elements {
			d.line = entry.Tkn().Line
			d.line = entry.Tkn().Pos

			if rv.Kind() == reflect.Ptr {
				if err := d.assignPrimitiveNode(entry, rv.Elem(), true); err != nil {
					setErr = err
					break switcher
				}
			} else {
				if err := d.assignPrimitiveNode(entry, rv, true); err != nil {
					setErr = err
					break switcher
				}
			}
		}

	case reflect.Map:
		val, ok := node.Value.(*MapNode)
		if !ok {
			setErr = errors.New("node is not a map")
			break switcher
		}

		if !rv.IsValid() {
			setErr = errors.New("invalid map")
			break switcher
		}

		if rv.IsNil() {
			rv.Set(reflect.MakeMap(rv.Type()))
		}

		for key, value := range val.Elements {
			d.line = value.Tkn().Line
			d.line = value.Tkn().Pos

			t := reflect.New(rv.Type().Elem())

			var isSlice bool
			rt := rv.Type()

			if rt.Kind() == reflect.Slice {
				rt = rt.Elem()
				isSlice = true
			}

			if err := d.assignPrimitiveNode(value, t.Elem(), isSlice); err != nil {
				setErr = err
				break switcher
			}

			// NB: this looks like pointless repetition but go uses the common interface type if i
			//     try to do this with a single case of *StringNode, *Identifier
			switch k := key.(type) {
			case *StringNode:
				rv.SetMapIndex(reflect.ValueOf(k.Value), t.Elem())
			case *Identifier:
				rv.SetMapIndex(reflect.ValueOf(k.Value), t.Elem())
			default:
				setErr = errors.New("Map keys must be a string")
				break switcher
			}
		}

	default:
		setErr = errors.New(path + ": unknown type " + rk.String())
	}

	if setErr != nil {
		return fmt.Errorf("%s: %w", path, setErr)
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

func (d *Decoder) assignPrimitiveNode(node Node, rv reflect.Value, isSlice bool) error {
	d.line = node.Tkn().Line
	d.pos = node.Tkn().Pos

	switch v := node.(type) {
	case *EnvarNode:
		val := os.Getenv(v.Identifier.Value)
		rk := rv.Kind()
		if rk == reflect.Ptr {
			rk = rv.Elem().Kind()
		}
		switch rk {
		case reflect.String:
			assignReflectValue(rv, val, false)
		case reflect.Bool:
			assignReflectValue(rv, val == "true", false)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			val, err := parseIntKind(val, rk)
			if err != nil {
				return err
			}
			assignReflectValue(rv, val, false)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			val, err := parseUintKind(val, rk)
			if err != nil {
				return err
			}
			assignReflectValue(rv, val, false)
		case reflect.Float32, reflect.Float64:
			bs := 32
			if rk == reflect.Float64 {
				bs = 64
			}

			v, err := strconv.ParseFloat(val, bs)
			if err != nil {
				return err
			}

			if rk == reflect.Float32 {
				assignReflectValue(rv, float32(v), false)
			} else {
				assignReflectValue(rv, v, false)
			}
		}
	case *NullNode:
		if rv.Kind() != reflect.Ptr {
			return fmt.Errorf("invalid %v type null", baseKind(rv))
		}
	case *StringNode:
		if !checkReflectKind(rv, reflect.String, isSlice) {
			return fmt.Errorf("invalid %v type string", baseKind(rv))
		}

		assignReflectValue(rv, v.Value, isSlice)
	case *BooleanNode:
		if !checkReflectKind(rv, reflect.Bool, isSlice) {
			return fmt.Errorf("invalid %v type bool", baseKind(rv))
		}

		assignReflectValue(rv, v.Value, isSlice)
	case *NumberNode:
		rk := baseKind(rv)

		switch rk {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			val, err := parseIntKind(v.Value, rk)
			if err != nil {
				return err
			}

			assignReflectValue(rv, val, isSlice)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			val, err := parseUintKind(v.Value, rk)
			if err != nil {
				return err
			}

			assignReflectValue(rv, val, isSlice)
		case reflect.Float32, reflect.Float64:
			bs := 32
			if rk == reflect.Float64 {
				bs = 64
			}

			val, err := strconv.ParseFloat(v.Value, bs)
			if err != nil {
				return err
			}

			if rk == reflect.Float32 {
				assignReflectValue(rv, float32(val), isSlice)
			} else {
				assignReflectValue(rv, val, isSlice)
			}
		default:
			return errors.New("invalid type " + string(v.Tkn().Type) + " : " + rk.String())
		}
	default:
		return errors.New("invalid node type " + string(v.Tkn().Type))
	}

	return nil
}

func parseIntKind(s string, k reflect.Kind) (any, error) {
	switch k {
	case reflect.Int8:
		i, err := strconv.ParseInt(s, 10, 8)
		return int8(i), err
	case reflect.Int16:
		i, err := strconv.ParseInt(s, 10, 16)
		return int16(i), err
	case reflect.Int32:
		i, err := strconv.ParseInt(s, 10, 32)
		return int32(i), err
	case reflect.Int64:
		return strconv.ParseInt(s, 10, 64)
	case reflect.Int:
		i, err := strconv.ParseInt(s, 10, 64)
		return int(i), err
	}

	return 0, errors.New("invilid int type")
}

func parseUintKind(s string, k reflect.Kind) (any, error) {
	switch k {
	case reflect.Uint8:
		i, err := strconv.ParseUint(s, 10, 8)
		return uint8(i), err
	case reflect.Uint16:
		i, err := strconv.ParseUint(s, 10, 16)
		return uint16(i), err
	case reflect.Uint32:
		i, err := strconv.ParseUint(s, 10, 32)
		return uint32(i), err
	case reflect.Uint64:
		return strconv.ParseUint(s, 10, 64)
	case reflect.Uint:
		i, err := strconv.ParseUint(s, 10, 64)
		return uint(i), err
	}

	return 0, errors.New("invilid int type")
}

func checkReflectKind(rv reflect.Value, expected reflect.Kind, isSlice bool) bool {
	if isSlice {
		if rv.Kind() != reflect.Slice {
			return false
		}

		if rv.Type().Elem().Kind() == expected {
			return true
		}

		if rv.Type().Elem().Kind() == reflect.Ptr && rv.Type().Elem().Elem().Kind() == expected {
			return true
		}

		return false
	}

	if rv.Kind() == expected {
		return true
	}

	if rv.Kind() == reflect.Ptr && rv.Type().Elem().Kind() == expected {
		return true
	}

	return false
}

func assignReflectValue[T any](rv reflect.Value, val T, isSlice bool) {
	if isSlice {
		if rv.Type().Elem().Kind() == reflect.Ptr {
			rv.Set(reflect.Append(rv, reflect.ValueOf(&val)))
		} else {
			rv.Set(reflect.Append(rv, reflect.ValueOf(val)))
		}
	} else if rv.Kind() == reflect.Ptr {
		ptr := reflect.New(reflect.TypeOf(val))
		ptr.Elem().Set(reflect.ValueOf(val))
		rv.Set(ptr)
	} else {
		rv.Set(reflect.ValueOf(val))
	}
}

func baseKind(rv reflect.Value) reflect.Kind {
	rk := rv.Kind()
	if rk == reflect.Ptr {
		rk = rv.Elem().Kind()
		if rk == reflect.Invalid {
			rk = rv.Type().Elem().Kind()
		}
	} else if rk == reflect.Slice {
		rk = rv.Type().Elem().Kind()
		if rk == reflect.Ptr {
			rk = rv.Type().Elem().Elem().Kind()
		}
	}

	return rk
}
