package test

import (
	"testing"

	"github.com/indeedhat/icl"
	"github.com/stretchr/testify/require"
)

type MarshalTarget struct {
	U     uint                `icl:"u"`
	Up    *uint               `icl:"up"`
	Us    []uint              `icl:"us"`
	Ups   []*uint             `icl:"ups"`
	Um    map[string]uint     `icl:"um"`
	Upm   map[string]*uint    `icl:"upm"`
	U8    uint8               `icl:"u8"`
	Up8   *uint8              `icl:"up8"`
	Us8   []uint8             `icl:"us8"`
	Ups8  []*uint8            `icl:"ups8"`
	Um8   map[string]uint8    `icl:"um8"`
	Upm8  map[string]*uint8   `icl:"upm8"`
	U16   uint16              `icl:"u16"`
	Up16  *uint16             `icl:"up16"`
	Us16  []uint16            `icl:"us16"`
	Ups16 []*uint16           `icl:"ups16"`
	Um16  map[string]uint16   `icl:"um16"`
	Upm16 map[string]*uint16  `icl:"upm16"`
	U32   uint32              `icl:"u32"`
	Up32  *uint32             `icl:"up32"`
	Us32  []uint32            `icl:"us32"`
	Ups32 []*uint32           `icl:"ups32"`
	Um32  map[string]uint32   `icl:"um32"`
	Upm32 map[string]*uint32  `icl:"upm32"`
	U64   uint64              `icl:"u64"`
	Up64  *uint64             `icl:"up64"`
	Us64  []uint64            `icl:"us64"`
	Ups64 []*uint64           `icl:"ups64"`
	Um64  map[string]uint64   `icl:"um64"`
	Upm64 map[string]*uint64  `icl:"upm64"`
	I     int                 `icl:"i"`
	Ip    *int                `icl:"ip"`
	Is    []int               `icl:"is"`
	Ips   []*int              `icl:"ips"`
	Im    map[string]int      `icl:"im"`
	Ipm   map[string]*int     `icl:"ipm"`
	I8    int8                `icl:"i8"`
	Ip8   *int8               `icl:"ip8"`
	Is8   []int8              `icl:"is8"`
	Ips8  []*int8             `icl:"ips8"`
	Im8   map[string]int8     `icl:"im8"`
	Ipm8  map[string]*int8    `icl:"ipm8"`
	I16   int16               `icl:"i16"`
	Ip16  *int16              `icl:"ip16"`
	Is16  []int16             `icl:"is16"`
	Ips16 []*int16            `icl:"ips16"`
	Im16  map[string]int16    `icl:"im16"`
	Ipm16 map[string]*int16   `icl:"ipm16"`
	I32   int32               `icl:"i32"`
	Ip32  *int32              `icl:"ip32"`
	Is32  []int32             `icl:"is32"`
	Ips32 []*int32            `icl:"ips32"`
	Im32  map[string]int32    `icl:"im32"`
	Ipm32 map[string]*int32   `icl:"ipm32"`
	I64   int64               `icl:"i64"`
	Ip64  *int64              `icl:"ip64"`
	Is64  []int64             `icl:"is64"`
	Ips64 []*int64            `icl:"ips64"`
	Im64  map[string]int64    `icl:"im64"`
	Ipm64 map[string]*int64   `icl:"ipm64"`
	S     string              `icl:"s"`
	Sp    *string             `icl:"sp"`
	Ss    []string            `icl:"ss"`
	Sps   []*string           `icl:"sps"`
	Sm    map[string]string   `icl:"sm"`
	Spm   map[string]*string  `icl:"spm"`
	B     bool                `icl:"b"`
	Bp    *bool               `icl:"bp"`
	Bs    []bool              `icl:"bs"`
	Bps   []*bool             `icl:"bps"`
	Bm    map[string]bool     `icl:"bm"`
	Bpm   map[string]*bool    `icl:"bpm"`
	F32   float32             `icl:"f32.2"`
	Fp32  *float32            `icl:"fp32.2"`
	Fs32  []float32           `icl:"fs32.2"`
	Fps32 []*float32          `icl:"fps32.2"`
	Fm32  map[string]float32  `icl:"fm32.2"`
	Fpm32 map[string]*float32 `icl:"fpm32.2"`
	F64   float64             `icl:"f64.2"`
	Fp64  *float64            `icl:"fp64.2"`
	Fs64  []float64           `icl:"fs64.2"`
	Fps64 []*float64          `icl:"fps64.2"`
	Fm64  map[string]float64  `icl:"fm64.2"`
	Fpm64 map[string]*float64 `icl:"fpm64.2"`
	Sb    SubBlock            `icl:"sb"`
	Sbp   *SubBlock           `icl:"sbp"`
	Sbs   []SubBlock          `icl:"sbs"`
	Sbw   SubBlockWrapper     `icl:"sbw"`
}

type SubBlockWrapper struct {
	Sb  SubBlock   `icl:"sb"`
	Sbp *SubBlock  `icl:"sbp"`
	Sbs []SubBlock `icl:"sbs"`
}
type SubBlock struct {
	P1   string `icl:".param"`
	P2   string `icl:".param"`
	Data string `icl:"data"`
}

const expectedMarshalDocument = `u = 120
up = 120
us = [1, 2, 3]
ups = [1, 2, 3]
um = {
    "1": 1,
    "2": 2,
    "3": 3,
}
upm = {
    "1": 1,
    "2": 2,
    "3": 3,
}
u8 = 120
up8 = 120
us8 = [1, 2, 3]
ups8 = [1, 2, 3]
um8 = {
    "1": 1,
    "2": 2,
    "3": 3,
}
upm8 = {
    "1": 1,
    "2": 2,
    "3": 3,
}
u16 = 120
up16 = 120
us16 = [1, 2, 3]
ups16 = [1, 2, 3]
um16 = {
    "1": 1,
    "2": 2,
    "3": 3,
}
upm16 = {
    "1": 1,
    "2": 2,
    "3": 3,
}
u32 = 120
up32 = 120
us32 = [1, 2, 3]
ups32 = [1, 2, 3]
um32 = {
    "1": 1,
    "2": 2,
    "3": 3,
}
upm32 = {
    "1": 1,
    "2": 2,
    "3": 3,
}
u64 = 120
up64 = 120
us64 = [1, 2, 3]
ups64 = [1, 2, 3]
um64 = {
    "1": 1,
    "2": 2,
    "3": 3,
}
upm64 = {
    "1": 1,
    "2": 2,
    "3": 3,
}
i = 120
ip = 120
is = [1, 2, 3]
ips = [1, 2, 3]
im = {
    "1": 1,
    "2": 2,
    "3": 3,
}
ipm = {
    "1": 1,
    "2": 2,
    "3": 3,
}
i8 = 120
ip8 = 120
is8 = [1, 2, 3]
ips8 = [1, 2, 3]
im8 = {
    "1": 1,
    "2": 2,
    "3": 3,
}
ipm8 = {
    "1": 1,
    "2": 2,
    "3": 3,
}
i16 = 120
ip16 = 120
is16 = [1, 2, 3]
ips16 = [1, 2, 3]
im16 = {
    "1": 1,
    "2": 2,
    "3": 3,
}
ipm16 = {
    "1": 1,
    "2": 2,
    "3": 3,
}
i32 = 120
ip32 = 120
is32 = [1, 2, 3]
ips32 = [1, 2, 3]
im32 = {
    "1": 1,
    "2": 2,
    "3": 3,
}
ipm32 = {
    "1": 1,
    "2": 2,
    "3": 3,
}
i64 = 120
ip64 = 120
is64 = [1, 2, 3]
ips64 = [1, 2, 3]
im64 = {
    "1": 1,
    "2": 2,
    "3": 3,
}
ipm64 = {
    "1": 1,
    "2": 2,
    "3": 3,
}
s = "str"
sp = "str"
ss = ["str1", "str2", "str3"]
sps = ["str1", "str2", "str3"]
sm = {
    "1": "str1",
    "2": "str2",
    "3": "str3",
}
spm = {
    "1": "str1",
    "2": "str2",
    "3": "str3",
}
b = true
bp = true
bs = [true, false, true]
bps = [true, false, true]
bm = {
    "1": true,
    "2": false,
    "3": true,
}
bpm = {
    "1": true,
    "2": false,
    "3": true,
}
f32 = 120.00
fp32 = 120.00
fs32 = [1.00, 2.00, 3.00]
fps32 = [1.00, 2.00, 3.00]
fm32 = {
    "1": 1.00,
    "2": 2.00,
    "3": 3.00,
}
fpm32 = {
    "1": 1.00,
    "2": 2.00,
    "3": 3.00,
}
f64 = 120.00
fp64 = 120.00
fs64 = [1.00, 2.00, 3.00]
fps64 = [1.00, 2.00, 3.00]
fm64 = {
    "1": 1.00,
    "2": 2.00,
    "3": 3.00,
}
fpm64 = {
    "1": 1.00,
    "2": 2.00,
    "3": 3.00,
}
sb "param1" "param2" {
    data = "data"
}
sbp "param1" "param2" {
    data = "data"
}
sbs "param1" "param2" {
    data = "data"
}
sbs "param1" "param2" {
    data = "data2"
}
sbs "param1" "param2" {
    data = "data3"
}
sbw {
    sb "param1" "param2" {
        data = "data"
    }
    sbp "param1" "param2" {
        data = "data"
    }
    sbs "param1" "param2" {
        data = "data"
    }
    sbs "param1" "param2" {
        data = "data2"
    }
    sbs "param1" "param2" {
        data = "data3"
    }
}
`

var marshalTarget = MarshalTarget{
	U:     120,
	Up:    ptr[uint](120),
	Us:    []uint{1, 2, 3},
	Ups:   []*uint{ptr[uint](1), ptr[uint](2), ptr[uint](3)},
	Um:    map[string]uint{"1": 1, "2": 2, "3": 3},
	Upm:   map[string]*uint{"1": ptr[uint](1), "2": ptr[uint](2), "3": ptr[uint](3)},
	U8:    120,
	Up8:   ptr[uint8](120),
	Us8:   []uint8{1, 2, 3},
	Ups8:  []*uint8{ptr[uint8](1), ptr[uint8](2), ptr[uint8](3)},
	Um8:   map[string]uint8{"1": 1, "2": 2, "3": 3},
	Upm8:  map[string]*uint8{"1": ptr[uint8](1), "2": ptr[uint8](2), "3": ptr[uint8](3)},
	U16:   120,
	Up16:  ptr[uint16](120),
	Us16:  []uint16{1, 2, 3},
	Ups16: []*uint16{ptr[uint16](1), ptr[uint16](2), ptr[uint16](3)},
	Um16:  map[string]uint16{"1": 1, "2": 2, "3": 3},
	Upm16: map[string]*uint16{"1": ptr[uint16](1), "2": ptr[uint16](2), "3": ptr[uint16](3)},
	U32:   120,
	Up32:  ptr[uint32](120),
	Us32:  []uint32{1, 2, 3},
	Ups32: []*uint32{ptr[uint32](1), ptr[uint32](2), ptr[uint32](3)},
	Um32:  map[string]uint32{"1": 1, "2": 2, "3": 3},
	Upm32: map[string]*uint32{"1": ptr[uint32](1), "2": ptr[uint32](2), "3": ptr[uint32](3)},
	U64:   120,
	Up64:  ptr[uint64](120),
	Us64:  []uint64{1, 2, 3},
	Ups64: []*uint64{ptr[uint64](1), ptr[uint64](2), ptr[uint64](3)},
	Um64:  map[string]uint64{"1": 1, "2": 2, "3": 3},
	Upm64: map[string]*uint64{"1": ptr[uint64](1), "2": ptr[uint64](2), "3": ptr[uint64](3)},
	I:     120,
	Ip:    ptr(120),
	Is:    []int{1, 2, 3},
	Ips:   []*int{ptr(1), ptr(2), ptr(3)},
	Im:    map[string]int{"1": 1, "2": 2, "3": 3},
	Ipm:   map[string]*int{"1": ptr(1), "2": ptr(2), "3": ptr(3)},
	I8:    120,
	Ip8:   ptr[int8](120),
	Is8:   []int8{1, 2, 3},
	Ips8:  []*int8{ptr[int8](1), ptr[int8](2), ptr[int8](3)},
	Im8:   map[string]int8{"1": 1, "2": 2, "3": 3},
	Ipm8:  map[string]*int8{"1": ptr[int8](1), "2": ptr[int8](2), "3": ptr[int8](3)},
	I16:   120,
	Ip16:  ptr[int16](120),
	Is16:  []int16{1, 2, 3},
	Ips16: []*int16{ptr[int16](1), ptr[int16](2), ptr[int16](3)},
	Im16:  map[string]int16{"1": 1, "2": 2, "3": 3},
	Ipm16: map[string]*int16{"1": ptr[int16](1), "2": ptr[int16](2), "3": ptr[int16](3)},
	I32:   120,
	Ip32:  ptr[int32](120),
	Is32:  []int32{1, 2, 3},
	Ips32: []*int32{ptr[int32](1), ptr[int32](2), ptr[int32](3)},
	Im32:  map[string]int32{"1": 1, "2": 2, "3": 3},
	Ipm32: map[string]*int32{"1": ptr[int32](1), "2": ptr[int32](2), "3": ptr[int32](3)},
	I64:   120,
	Ip64:  ptr[int64](120),
	Is64:  []int64{1, 2, 3},
	Ips64: []*int64{ptr[int64](1), ptr[int64](2), ptr[int64](3)},
	Im64:  map[string]int64{"1": 1, "2": 2, "3": 3},
	Ipm64: map[string]*int64{"1": ptr[int64](1), "2": ptr[int64](2), "3": ptr[int64](3)},
	S:     "str",
	Sp:    ptr("str"),
	Ss:    []string{"str1", "str2", "str3"},
	Sps:   []*string{ptr("str1"), ptr("str2"), ptr("str3")},
	Sm:    map[string]string{"1": "str1", "2": "str2", "3": "str3"},
	Spm:   map[string]*string{"1": ptr("str1"), "2": ptr("str2"), "3": ptr("str3")},
	B:     true,
	Bp:    ptr(true),
	Bs:    []bool{true, false, true},
	Bps:   []*bool{ptr(true), ptr(false), ptr(true)},
	Bm:    map[string]bool{"1": true, "2": false, "3": true},
	Bpm:   map[string]*bool{"1": ptr(true), "2": ptr(false), "3": ptr(true)},
	F32:   120,
	Fp32:  ptr[float32](120),
	Fs32:  []float32{1, 2, 3},
	Fps32: []*float32{ptr[float32](1), ptr[float32](2), ptr[float32](3)},
	Fm32:  map[string]float32{"1": 1, "2": 2, "3": 3},
	Fpm32: map[string]*float32{"1": ptr[float32](1), "2": ptr[float32](2), "3": ptr[float32](3)},
	F64:   120,
	Fp64:  ptr[float64](120),
	Fs64:  []float64{1, 2, 3},
	Fps64: []*float64{ptr[float64](1), ptr[float64](2), ptr[float64](3)},
	Fm64:  map[string]float64{"1": 1, "2": 2, "3": 3},
	Fpm64: map[string]*float64{"1": ptr[float64](1), "2": ptr[float64](2), "3": ptr[float64](3)},
	Sb: SubBlock{
		P1:   "param1",
		P2:   "param2",
		Data: "data",
	},
	Sbp: &SubBlock{
		P1:   "param1",
		P2:   "param2",
		Data: "data",
	},
	Sbs: []SubBlock{
		{
			P1:   "param1",
			P2:   "param2",
			Data: "data",
		},
		{
			P1:   "param1",
			P2:   "param2",
			Data: "data2",
		},
		{
			P1:   "param1",
			P2:   "param2",
			Data: "data3",
		},
	},
	Sbw: SubBlockWrapper{
		Sb: SubBlock{
			P1:   "param1",
			P2:   "param2",
			Data: "data",
		},
		Sbp: &SubBlock{
			P1:   "param1",
			P2:   "param2",
			Data: "data",
		},
		Sbs: []SubBlock{
			{
				P1:   "param1",
				P2:   "param2",
				Data: "data",
			},
			{
				P1:   "param1",
				P2:   "param2",
				Data: "data2",
			},
			{
				P1:   "param1",
				P2:   "param2",
				Data: "data3",
			},
		},
	},
}

func TestMarshalDocument(t *testing.T) {
	document, err := icl.MarshalString(marshalTarget)

	require.Nil(t, err)
	require.Equal(t, expectedMarshalDocument, document)
}
