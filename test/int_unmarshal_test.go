package test

import (
	"math"
	"testing"

	"github.com/indeedhat/icl"
	"github.com/stretchr/testify/require"
)

// int
type intTarget struct {
	I   int   `icl:"i"`
	I8  int8  `icl:"i8"`
	I16 int16 `icl:"i16"`
	I32 int32 `icl:"i32"`
	I64 int64 `icl:"i64"`
}

var intUnmarshalTests = map[string]unmarshalTest{
	"int valid": {
		`i = 9223372036854775807`,
		intTarget{I: math.MaxInt},
		"",
	},
	"int negative": {
		`i = -9223372036854775808`,
		intTarget{I: math.MinInt},
		"",
	},
	"int bad type": {
		`i = "bad"`,
		intTarget{I: 0},
		"",
	},
	"int8 valid": {
		`i8 = 127`,
		intTarget{I8: math.MaxInt8},
		"",
	},
	"int8 negative": {
		`i8 = -128`,
		intTarget{I8: math.MinInt8},
		"",
	},
	"int8 bad type": {
		`i8 = "bad"`,
		intTarget{I8: 0},
		"",
	},
	"int8 too large": {
		`i8 = 129`,
		intTarget{I8: 0},
		"",
	},
	"int16 valid": {
		`i16 = 32767`,
		intTarget{I16: math.MaxInt16},
		"",
	},
	"int16 negative": {
		`i16 = -32768`,
		intTarget{I16: math.MinInt16},
		"",
	},
	"int16 bad type": {
		`i16 = "bad"`,
		intTarget{I16: 0},
		"",
	},
	"int16 too large": {
		`i16 = 32768`,
		intTarget{I16: 0},
		"",
	},
	"int32 valid": {
		`i32 = 2147483647`,
		intTarget{I32: math.MaxInt32},
		"",
	},
	"int32 negative": {
		`i32 = -2147483648`,
		intTarget{I32: math.MinInt32},
		"",
	},
	"int32 bad type": {
		`i32 = "bad"`,
		intTarget{I32: 0},
		"",
	},
	"int32 too large": {
		`i32 = 2147483648`,
		intTarget{I32: 0},
		"",
	},
	"int64 valid": {
		`i64 = 9223372036854775807`,
		intTarget{I64: math.MaxInt64},
		"",
	},
	"int64 negative": {
		`i64 = -9223372036854775808`,
		intTarget{I64: math.MinInt64},
		"",
	},
	"int64 bad type": {
		`i64 = "bad"`,
		intTarget{I16: 0},
		"",
	},
}

func TestUnmarshalInt(t *testing.T) {
	for key, test := range intUnmarshalTests {
		t.Run(key, func(t *testing.T) {
			tgt := intTarget{}
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

type intPtrTarget struct {
	I   *int   `icl:"i"`
	I8  *int8  `icl:"i8"`
	I16 *int16 `icl:"i16"`
	I32 *int32 `icl:"i32"`
	I64 *int64 `icl:"i64"`
}

var intPtrUnmarshalTests = map[string]unmarshalTest{
	"int valid": {
		`i = 9223372036854775807`,
		intPtrTarget{I: ptr(math.MaxInt)},
		"",
	},
	"int negative": {
		`i = -9223372036854775808`,
		intPtrTarget{I: ptr(math.MinInt)},
		"",
	},
	"int bad type": {
		`i = "bad"`,
		intPtrTarget{},
		"",
	},
	"int8 valid": {
		`i8 = 127`,
		intPtrTarget{I8: ptr(int8(math.MaxInt8))},
		"",
	},
	"int8 negative": {
		`i8 = -128`,
		intPtrTarget{I8: ptr(int8(math.MinInt8))},
		"",
	},
	"int8 bad type": {
		`i8 = "bad"`,
		intPtrTarget{},
		"",
	},
	"int8 too large": {
		`i8 = 129`,
		intPtrTarget{},
		"",
	},
	"int16 valid": {
		`i16 = 32767`,
		intPtrTarget{I16: ptr(int16(math.MaxInt16))},
		"",
	},
	"int16 negative": {
		`i16 = -32768`,
		intPtrTarget{I16: ptr(int16(math.MinInt16))},
		"",
	},
	"int16 bad type": {
		`i16 = "bad"`,
		intPtrTarget{},
		"",
	},
	"int16 too large": {
		`i16 = 32768`,
		intPtrTarget{},
		"",
	},
	"int32 valid": {
		`i32 = 2147483647`,
		intPtrTarget{I32: ptr(int32(math.MaxInt32))},
		"",
	},
	"int32 negative": {
		`i32 = -2147483648`,
		intPtrTarget{I32: ptr(int32(math.MinInt32))},
		"",
	},
	"int32 bad type": {
		`i32 = "bad"`,
		intPtrTarget{},
		"",
	},
	"int32 too large": {
		`i32 = 2147483648`,
		intPtrTarget{},
		"",
	},
	"int64 valid": {
		`i64 = 9223372036854775807`,
		intPtrTarget{I64: ptr(int64(math.MaxInt64))},
		"",
	},
	"int64 negative": {
		`i64 = -9223372036854775808`,
		intPtrTarget{I64: ptr(int64(math.MinInt64))},
		"",
	},
	"int64 bad type": {
		`i64 = "bad"`,
		intPtrTarget{},
		"",
	},
}

func TestUnmarshalIntPtr(t *testing.T) {
	for key, test := range intPtrUnmarshalTests {
		t.Run(key, func(t *testing.T) {
			tgt := intPtrTarget{}
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

// uint
// float
// string
// slice
// map
// struct
