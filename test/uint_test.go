package test

import (
	"math"
	"testing"

	"github.com/indeedhat/icl"
	"github.com/stretchr/testify/require"
)

// uint
type uintTarget struct {
	I   uint   `icl:"i"`
	I8  uint8  `icl:"i8"`
	I16 uint16 `icl:"i16"`
	I32 uint32 `icl:"i32"`
	I64 uint64 `icl:"i64"`
}

var uintUnmarshalTests = map[string]unmarshalTest{
	"uint valid": {
		`i = 18446744073709551615`,
		uintTarget{I: uint(math.MaxUint)},
		"",
	},
	"uint bad type": {
		`i = "bad"`,
		uintTarget{},
		"",
	},
	"uint8 valid": {
		`i8 = 255`,
		uintTarget{I8: uint8(0xFF)},
		"",
	},
	"uint8 bad type": {
		`i8 = "bad"`,
		uintTarget{},
		"",
	},
	"uint8 too large": {
		`i8 = 256`,
		uintTarget{},
		"",
	},
	"uint16 valid": {
		`i16 = 65535`,
		uintTarget{I16: uint16(0xFFFF)},
		"",
	},
	"uint16 bad type": {
		`i16 = "bad"`,
		uintTarget{},
		"",
	},
	"uint16 too large": {
		`i16 = 65536`,
		uintTarget{},
		"",
	},
	"uint32 valid": {
		`i32 = 4294967295`,
		uintTarget{I32: uint32(0xFFFFFFFF)},
		"",
	},
	"uint32 bad type": {
		`i32 = "bad"`,
		uintTarget{},
		"",
	},
	"uint32 too large": {
		`i32 = 4294967296`,
		uintTarget{},
		"",
	},
	"uint64 valid": {
		`i64 = 18446744073709551615`,
		uintTarget{I64: uint64(0xFFFFFFFFFFFFFFFF)},
		"",
	},
	"uint64 bad type": {
		`i64 = "bad"`,
		uintTarget{},
		"",
	},
}

func TestUnmarshalUint(t *testing.T) {
	for key, test := range uintUnmarshalTests {
		t.Run(key, func(t *testing.T) {
			tgt := uintTarget{}
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

type uintPtrTarget struct {
	I   *uint   `icl:"i"`
	I8  *uint8  `icl:"i8"`
	I16 *uint16 `icl:"i16"`
	I32 *uint32 `icl:"i32"`
	I64 *uint64 `icl:"i64"`
}

var uintPtrUnmarshalTests = map[string]unmarshalTest{
	"uint valid": {
		`i = 18446744073709551615`,
		uintPtrTarget{I: ptr(uint(math.MaxUint))},
		"",
	},
	"uint bad type": {
		`i = "bad"`,
		uintPtrTarget{},
		"",
	},
	"uint8 valid": {
		`i8 = 255`,
		uintPtrTarget{I8: ptr(uint8(0xFF))},
		"",
	},
	"uint8 bad type": {
		`i8 = "bad"`,
		uintPtrTarget{},
		"",
	},
	"uint8 too large": {
		`i8 = 256`,
		uintPtrTarget{},
		"",
	},
	"uint16 valid": {
		`i16 = 65535`,
		uintPtrTarget{I16: ptr(uint16(0xFFFF))},
		"",
	},
	"uint16 bad type": {
		`i16 = "bad"`,
		uintPtrTarget{},
		"",
	},
	"uint16 too large": {
		`i16 = 65536`,
		uintPtrTarget{},
		"",
	},
	"uint32 valid": {
		`i32 = 4294967295`,
		uintPtrTarget{I32: ptr(uint32(0xFFFFFFFF))},
		"",
	},
	"uint32 bad type": {
		`i32 = "bad"`,
		uintPtrTarget{},
		"",
	},
	"uint32 too large": {
		`i32 = 4294967296`,
		uintPtrTarget{},
		"",
	},
	"uint64 valid": {
		`i64 = 18446744073709551615`,
		uintPtrTarget{I64: ptr(uint64(0xFFFFFFFFFFFFFFFF))},
		"",
	},
	"uint64 bad type": {
		`i64 = "bad"`,
		uintPtrTarget{},
		"",
	},
}

func TestUnmarshalUintPtr(t *testing.T) {
	for key, test := range uintPtrUnmarshalTests {
		t.Run(key, func(t *testing.T) {
			tgt := uintPtrTarget{}
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

// uuint
// float
// string
// slice
// map
// struct
