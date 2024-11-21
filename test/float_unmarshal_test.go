package test

import (
	"testing"

	"github.com/indeedhat/icl"
	"github.com/stretchr/testify/require"
)

type floatTarget struct {
	F32 float32 `icl:"f32"`
	F64 float64 `icl:"f64"`
}

var floatUnmarshalTests = map[string]unmarshalTest{
	"float32 valid": {
		`f32 = 1283.1`,
		floatTarget{F32: 1283.1},
		"",
	},
	"float32 from int": {
		`f32 = 128`,
		floatTarget{F32: 128},
		"",
	},
	"float32 bad type": {
		`f32 = "bad"`,
		floatTarget{},
		"",
	},
	"float64 valid": {
		`f64 = 1283.1`,
		floatTarget{F64: 1283.1},
		"",
	},
	"float64 from int": {
		`f64 = 128`,
		floatTarget{F64: 128},
		"",
	},
	"float64 bad type": {
		`f64 = "bad"`,
		floatTarget{},
		"",
	},
}

func TestUnmarshalFloat(t *testing.T) {
	for key, test := range floatUnmarshalTests {
		t.Run(key, func(t *testing.T) {
			tgt := floatTarget{}
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

type floatPtrTarget struct {
	F32 *float32 `icl:"f32"`
	F64 *float64 `icl:"f64"`
}

var floatPtrUnmarshalTests = map[string]unmarshalTest{
	"float32 valid": {
		`f32 = 1283.1`,
		floatPtrTarget{F32: ptr[float32](1283.1)},
		"",
	},
	"float32 from int": {
		`f32 = 128`,
		floatPtrTarget{F32: ptr[float32](128)},
		"",
	},
	"float32 bad type": {
		`f32 = "bad"`,
		floatPtrTarget{},
		"",
	},
	"float64 valid": {
		`f64 = 1283.1`,
		floatPtrTarget{F64: ptr(1283.1)},
		"",
	},
	"float64 from int": {
		`f64 = 128`,
		floatPtrTarget{F64: ptr[float64](128)},
		"",
	},
	"float64 bad type": {
		`f64 = "bad"`,
		floatPtrTarget{},
		"",
	},
}

func TestUnmarshalFloatPtr(t *testing.T) {
	for key, test := range floatPtrUnmarshalTests {
		t.Run(key, func(t *testing.T) {
			tgt := floatPtrTarget{}
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
