package main

import (
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/indeedhat/icl"
)

func main() {
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
