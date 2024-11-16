package main

import (
	"fmt"
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/indeedhat/icl"
)

func main() {
	// parse()
	marshal()
}

func parse() {
	ast, err := icl.Parse(`
		str_var = "things"
		null_var = null
		true_var = true
		false_var = false
		int_var = 39283
		float_var = 3.14.
		array_var = ["some", "data"]
		# comment
		map_var = {
			ident: "val",
			"string": 123,
			array_val: ["1","2","3", "4"], # trailing comma
			map_val: {data: true,},
		}
		my_block "with" data {
			# comment
			inner_data = true
			# comment
			with_map = {data: true}

			inner_block {
				# comment
			}
		}
		my_envar = env(SECRET_TOKEN)
	`)
	if err != nil {
		log.Fatal(err)
	}

	spew.Dump(ast)
}

type t struct {
	String   string   `icl:"string"`
	Boolean  bool     `icl:"boolean"`
	Nil      *string  `icl:"nullable"`
	NilWith  *string  `icl:"nullable_with_val"`
	Int      int8     `icl:"integer"`
	Uint     uint     `icl:"unsigned_integer"`
	Float    float32  `icl:"float.2"`
	StrSlice []string `icl:"string_slice"`
	NoInc    string
}

func marshal() {
	str := "my string"
	v := t{
		String:      "yeppas",
		Boolean:     false,
		Int:         -17,
		Uint:        28362,
		Float:       38.2223,
		NoInc:       "noooop",
		StrSlice:    []string{"one", "two", "three"},
		StrPtrSlice: []string{&str},
		NilWith:     &str,
	}

	fmt.Print(icl.NewEncoder(v))
}
