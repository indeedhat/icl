package test

import (
	"testing"

	"github.com/indeedhat/icl"
	"github.com/stretchr/testify/require"
)

// uint
type stringTarget struct {
	S  string  `icl:"s"`
	SP *string `icl:"sp"`
}

var stringUnmarshalTests = map[string]unmarshalTest{
	"string valid": {
		`s = "a string"`,
		stringTarget{S: "a string"},
		"",
	},
	"string invalid": {
		`s = []`,
		stringTarget{S: ""},
		"",
	},
	"*string valid": {
		`sp = "a string"`,
		stringTarget{SP: ptr("a string")},
		"",
	},
	"*string invalid": {
		`sp = []`,
		stringTarget{},
		"",
	},
}

func TestUnmarshalString(t *testing.T) {
	for key, test := range stringUnmarshalTests {
		t.Run(key, func(t *testing.T) {
			tgt := stringTarget{}
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
