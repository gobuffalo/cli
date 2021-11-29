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

	tcases := []struct {
		name    string
		args    []string
		content string
	}{
		{name: "Plain text", args: []string{"version"}, content: "version"},
		{name: "JSON", args: []string{"version", "--json"}, content: "\"version\":"},
	}

	for _, tcase := range tcases {
		t.Run(tcase.name, func(tx *testing.T) {
			rx := require.New(tx)
			out, err := testhelpers.RunBuffaloCMD(t, tcase.args)

			rx.NoError(err)
			rx.Contains(out, tcase.content)
		})
	}
}
