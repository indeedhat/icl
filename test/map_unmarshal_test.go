package test

import (
	"testing"

	"github.com/indeedhat/icl"
	"github.com/stretchr/testify/require"
)

type mapTarget struct {
	IntMap     map[string]int     `icl:"int_map"`
	Float64Map map[string]float64 `icl:"float64_map"`
	StringMap  map[string]string  `icl:"string_map"`
}

var mapUnmarshalTests = map[string]unmarshalTest{
	"int map valid": {
		`int_map = {"one": 1, "two": 2, "three": 3}`,
		mapTarget{IntMap: map[string]int{"one": 1, "two": 2, "three": 3}},
		"",
	},
	"int map invalid value type": {
		`int_map = {"one": "bad", "two": 2}`,
		mapTarget{IntMap: map[string]int{}},
		".int_map: invalid int type string\nline(0) pos(20)",
	},
	"int map empty": {
		`int_map = {}`,
		mapTarget{IntMap: map[string]int{}},
		"",
	},
	"float32 map valid": {
		`float64_map = {"pi": 3.14, "e": 2.71}`,
		mapTarget{Float64Map: map[string]float64{"pi": 3.14, "e": 2.71}},
		"",
	},
	"string map valid": {
		`string_map = {"key1": "value1", key2: "value2"}`,
		mapTarget{StringMap: map[string]string{"key1": "value1", "key2": "value2"}},
		"",
	},
	"string map invalid key type": {
		`string_map = {1: "value1"}`,
		mapTarget{StringMap: map[string]string{}},
		".string_map: Map keys must be a string\nline(0) pos(19)",
	},
}

func TestUnmarshalMaps(t *testing.T) {
	for key, test := range mapUnmarshalTests {
		t.Run(key, func(t *testing.T) {
			tgt := mapTarget{}
			err := icl.UnMarshalString(test.document, &tgt)

			if test.error != "" {
				require.NotNil(t, err)
				require.Equal(t, test.error, err.Error())
			} else {
				require.Nil(t, err)
			}

			require.Equal(t, test.output, tgt)
		})
	}
}
