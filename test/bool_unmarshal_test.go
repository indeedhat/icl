package test

import (
	"testing"

	"github.com/indeedhat/icl"
	"github.com/stretchr/testify/require"
)

// uint
type boolTarget struct {
	B  bool  `icl:"b"`
	BP *bool `icl:"bp"`
}

var boolUnmarshalTests = map[string]unmarshalTest{
	"bool true": {
		`b = true`,
		boolTarget{B: true},
		"",
	},
	"bool false": {
		`b = false`,
		boolTarget{B: false},
		"",
	},
	"bool invalid": {
		`b = ""`,
		boolTarget{},
		".b: invalid bool type string\nline(0) pos(6)",
	},
	"*bool true": {
		`bp = true`,
		boolTarget{BP: ptr(true)},
		"",
	},
	"*bool false": {
		`bp = false`,
		boolTarget{BP: ptr(false)},
		"",
	},
	"*bool invalid": {
		`bp = ""`,
		boolTarget{},
		".bp: invalid bool type string\nline(0) pos(7)",
	},
}

func TestUnmarshalBool(t *testing.T) {
	for key, test := range boolUnmarshalTests {
		t.Run(key, func(t *testing.T) {
			tgt := boolTarget{}
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
