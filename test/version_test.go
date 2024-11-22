package test

import (
	"testing"

	"github.com/indeedhat/icl"
	"github.com/stretchr/testify/require"
)

type VersionTest struct {
	Document string
	Expected int
}

var versionTests = map[string]VersionTest{
	"valid": {
		Document: `version = 1`,
		Expected: 1,
	},
	"proceeding comment": {
		Document: `# comment
		version = 1`,
		Expected: 1,
	},
	"no field": {
		Document: `# comment`,
		Expected: 0,
	},
	"not first statement": {
		Document: `title = "nope"
		version = 1`,
		Expected: 0,
	},
	"bad type": {
		Document: `version = 3.2`,
		Expected: 0,
	},
}

func TestVersionParsing(t *testing.T) {
	for name, testCase := range versionTests {
		t.Run(name, func(t *testing.T) {
			ast, err := icl.ParseString(testCase.Document)
			require.Nil(t, err)

			require.Equal(t, testCase.Expected, ast.Version())
		})
	}
}
