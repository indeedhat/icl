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
		".i: invalid uint type string\nline(0) pos(6)",
	},
	"uint8 valid": {
		`i8 = 255`,
		uintTarget{I8: uint8(0xFF)},
		"",
	},
	"uint8 bad type": {
		`i8 = "bad"`,
		uintTarget{},
		".i8: invalid uint8 type string\nline(0) pos(7)",
	},
	"uint8 too large": {
		`i8 = 256`,
		uintTarget{},
		".i8: strconv.ParseUint: parsing \"256\": value out of range\nline(0) pos(5)",
	},
	"uint16 valid": {
		`i16 = 65535`,
		uintTarget{I16: uint16(0xFFFF)},
		"",
	},
	"uint16 bad type": {
		`i16 = "bad"`,
		uintTarget{},
		".i16: invalid uint16 type string\nline(0) pos(8)",
	},
	"uint16 too large": {
		`i16 = 65536`,
		uintTarget{},
		".i16: strconv.ParseUint: parsing \"65536\": value out of range\nline(0) pos(6)",
	},
	"uint32 valid": {
		`i32 = 4294967295`,
		uintTarget{I32: uint32(0xFFFFFFFF)},
		"",
	},
	"uint32 bad type": {
		`i32 = "bad"`,
		uintTarget{},
		".i32: invalid uint32 type string\nline(0) pos(8)",
	},
	"uint32 too large": {
		`i32 = 4294967296`,
		uintTarget{},
		".i32: strconv.ParseUint: parsing \"4294967296\": value out of range\nline(0) pos(6)",
	},
	"uint64 valid": {
		`i64 = 18446744073709551615`,
		uintTarget{I64: uint64(0xFFFFFFFFFFFFFFFF)},
		"",
	},
	"uint64 bad type": {
		`i64 = "bad"`,
		uintTarget{},
		".i64: invalid uint64 type string\nline(0) pos(8)",
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
		".i: invalid uint type string\nline(0) pos(6)",
	},
	"uint8 valid": {
		`i8 = 255`,
		uintPtrTarget{I8: ptr(uint8(0xFF))},
		"",
	},
	"uint8 bad type": {
		`i8 = "bad"`,
		uintPtrTarget{},
		".i8: invalid uint8 type string\nline(0) pos(7)",
	},
	"uint8 too large": {
		`i8 = 256`,
		uintPtrTarget{},
		".i8: strconv.ParseUint: parsing \"256\": value out of range\nline(0) pos(5)",
	},
	"uint16 valid": {
		`i16 = 65535`,
		uintPtrTarget{I16: ptr(uint16(0xFFFF))},
		"",
	},
	"uint16 bad type": {
		`i16 = "bad"`,
		uintPtrTarget{},
		".i16: invalid uint16 type string\nline(0) pos(8)",
	},
	"uint16 too large": {
		`i16 = 65536`,
		uintPtrTarget{},
		".i16: strconv.ParseUint: parsing \"65536\": value out of range\nline(0) pos(6)",
	},
	"uint32 valid": {
		`i32 = 4294967295`,
		uintPtrTarget{I32: ptr(uint32(0xFFFFFFFF))},
		"",
	},
	"uint32 bad type": {
		`i32 = "bad"`,
		uintPtrTarget{},
		".i32: invalid uint32 type string\nline(0) pos(8)",
	},
	"uint32 too large": {
		`i32 = 4294967296`,
		uintPtrTarget{},
		".i32: strconv.ParseUint: parsing \"4294967296\": value out of range\nline(0) pos(6)",
	},
	"uint64 valid": {
		`i64 = 18446744073709551615`,
		uintPtrTarget{I64: ptr(uint64(0xFFFFFFFFFFFFFFFF))},
		"",
	},
	"uint64 bad type": {
		`i64 = "bad"`,
		uintPtrTarget{},
		".i64: invalid uint64 type string\nline(0) pos(8)",
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
