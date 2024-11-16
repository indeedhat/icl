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

type config struct {
	String      string             `icl:"string"`
	Boolean     bool               `icl:"boolean"`
	Nil         *string            `icl:"nullable"`
	NilWith     *string            `icl:"nullable_with_val"`
	Int         int8               `icl:"integer"`
	Uint        uint               `icl:"unsigned_integer"`
	Float       float32            `icl:"float.2"`
	StrSlice    []string           `icl:"string_slice"`
	IntSlice    []int              `icl:"int_slice"`
	FloatSlice  []float64          `icl:"float_slice.3"`
	Struct      inner              `icl:"struct"`
	Map         map[string]float64 `icl:"float_map.2"`
	StructSlice []inner2           `icl:"multi_block"`
	NoInc       string
	Envar       string `icl:"envar,env(MY_ENVAR)"`
}

type inner struct {
	P1     string `icl:".param"`
	P2     string `icl:".param"`
	Data   bool   `icl:"data"`
	Inner2 inner2 `icl:"inner_2"`
}

type inner2 struct {
	P1   string `icl:".param"`
	P2   string `icl:".param"`
	Data bool   `icl:"data"`
}

func marshal() {
	str := "my string"
	v := config{
		String:     "yeppas",
		Boolean:    false,
		Int:        -17,
		Uint:       28362,
		Float:      38.2223,
		NoInc:      "noooop",
		StrSlice:   []string{"one", "two", "three"},
		IntSlice:   []int{1, 2, 3},
		FloatSlice: []float64{1, 2, 3},
		NilWith:    &str,
		Struct: inner{
			P1:   "Param 1",
			P2:   "Param 2",
			Data: true,
			Inner2: inner2{
				P2:   "Inner param 2",
				P1:   "Inner param 1",
				Data: false,
			},
		},
		Map: map[string]float64{
			"one": 1,
			"two": 2,
		},
		StructSlice: []inner2{
			{
				P1:   "p1",
				P2:   "p2",
				Data: true,
			},
			{
				P1:   "p1",
				P2:   "p2",
				Data: false,
			},
		},
	}

	fmt.Print(icl.MarshalString(v))
}
