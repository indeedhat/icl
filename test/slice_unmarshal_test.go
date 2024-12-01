package test

import (
	"testing"

	"github.com/indeedhat/icl"
	"github.com/stretchr/testify/require"
)

type sliceTarget struct {
	IntSlice     []int     `icl:"int_slice"`
	Float64Slice []float64 `icl:"float64_slice"`
	StringSlice  []string  `icl:"string_slice"`
}

var sliceUnmarshalTests = map[string]unmarshalTest{
	"int slice valid": {
		`int_slice = [1, 2, 3]`,
		sliceTarget{IntSlice: []int{1, 2, 3}},
		"",
	},
	"int slice invalid type": {
		`int_slice = ["bad", 2, 3]`,
		sliceTarget{},
		".int_slice: invalid int type string\nline(0) pos(15)",
	},
	"int slice empty": {
		`int_slice = []`,
		sliceTarget{},
		"",
	},
	"float64 slice valid": {
		`float64_slice = [1.1, 2.2, 3.3]`,
		sliceTarget{Float64Slice: []float64{1.1, 2.2, 3.3}},
		"",
	},
	"float64 slice invalid": {
		`float64_slice = ["bad", 2.2, 3.3]`,
		sliceTarget{},
		".float64_slice: invalid float64 type string\nline(0) pos(19)",
	},
	"string slice valid": {
		`string_slice = ["a", "b", "c"]`,
		sliceTarget{StringSlice: []string{"a", "b", "c"}},
		"",
	},
	"string slice invalid": {
		`string_slice = [1, 2, 3]`,
		sliceTarget{},
		".string_slice: invalid type NUMBER : string\nline(0) pos(16)",
	},
}

func TestUnmarshalSlices(t *testing.T) {
	for key, test := range sliceUnmarshalTests {
		t.Run(key, func(t *testing.T) {
			tgt := sliceTarget{}
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
