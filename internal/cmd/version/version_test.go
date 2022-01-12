// +build integration

package version

import (
	"testing"

	"github.com/gobuffalo/cli/internal/testhelpers"
	"github.com/stretchr/testify/require"
)

func TestVersion(t *testing.T) {
	r := require.New(t)
	r.NoError(testhelpers.EnsureBuffaloCMD(t))

	tt := []struct {
		name    string
		args    []string
		content string
	}{
		{name: "Plain text", args: []string{"version"}, content: "version"},
		{name: "JSON", args: []string{"version", "--json"}, content: "\"version\":"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rx := require.New(t)
			out, err := testhelpers.RunBuffaloCMD(t, tc.args)

			rx.NoError(err)
			rx.Contains(out, tc.content)
		})
	}
}
