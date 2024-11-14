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
		array_var = ["some", "data"]
		# comment
	`)
	if err != nil {
		log.Fatal(err)
	}

	spew.Dump(ast)
}
