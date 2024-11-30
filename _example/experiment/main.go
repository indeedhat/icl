package main

import (
	"log"
	"reflect"

	"github.com/davecgh/go-spew/spew"
)

type Target struct {
	String         string
	StringPtr      *string
	StringSlice    []string
	StringPtrSlice []*string
	StringMap      map[string]string
	StringPtrMap   map[string]*string
}

var data = map[string]any{
	"string":           "my string",
	"string_ptr":       ptr("my string"),
	"string_slice":     []string{"my", "string"},
	"string_ptr_slice": []*string{ptr("my"), ptr("string")},
	"string_map":       map[string]string{"my": "my", "string": "string"},
	"string_ptr_map":   map[string]*string{"my": ptr("my"), "string": ptr("string")},
}

func ptr[T any](v T) *T {
	return &v
}

func main() {
	t := Target{}
	assignString(&t, "my string")
	assignStringPtr(&t, ptr("my string"))
	assignStringSlice(&t, []string{"my", "string"})
	assignStringPtrSlice(&t, []*string{ptr("my"), ptr("string")})
	assignStringMap(&t, map[string]string{"my": "my", "string": "string"})
	assignStringPtrMap(&t, map[string]*string{"my": ptr("my"), "string": ptr("string")})

	spew.Dump(t)
}

func assignString(t *Target, s string) {
	rv := reflect.ValueOf(t)

	field := rv.Elem().FieldByName("String")
	field.SetString(s)
}

func assignStringPtr(t *Target, s *string) {
	rv := reflect.ValueOf(t)

	field := rv.Elem().FieldByName("StringPtr")

	v := reflect.New(field.Type())
	v.Elem().Set(reflect.ValueOf(s))

	field.Set(v.Elem())
}

func assignStringSlice(t *Target, s []string) {
	rv := reflect.ValueOf(t)

	field := rv.Elem().FieldByName("StringSlice")

	for _, e := range s {
		field.Set(reflect.Append(field, reflect.ValueOf(e)))
	}
}

func assignStringPtrSlice(t *Target, s []*string) {
	rv := reflect.ValueOf(t)

	field := rv.Elem().FieldByName("StringPtrSlice")

	for _, e := range s {
		field.Set(reflect.Append(field, reflect.ValueOf(e)))
	}
}

func assignStringMap(t *Target, s map[string]string) {
	rv := reflect.ValueOf(t)

	field := rv.Elem().FieldByName("StringMap")
	log.Print(field.Kind())

	if field.IsZero() {
		field.Set(reflect.MakeMap(field.Type()))
	}

	for k, v := range s {
		field.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v))
	}
}

func assignStringPtrMap(t *Target, s map[string]*string) {
	rv := reflect.ValueOf(t)

	field := rv.Elem().FieldByName("StringPtrMap")

	if field.IsZero() {
		field.Set(reflect.MakeMap(field.Type()))
	}

	for k, v := range s {
		field.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v))
	}
}
