package docker

import (
	"fmt"
	"testing"

	"github.com/gobuffalo/genny/v2/gentest"
	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	r := require.New(t)

	g, err := New(&Options{})
	r.NoError(err)

	run := gentest.NewRunner()
	r.NoError(run.With(g))
	r.NoError(run.Run())

	res := run.Results()
	r.Len(res.Commands, 0)

	for _, v := range res.Files {

		fmt.Println("FILE >>>> ", v.Name())
	}

	r.Len(res.Files, 2)

	f := res.Files[0]
	r.Equal(".dockerignore", f.Name())

	f = res.Files[1]
	r.Equal("Dockerfile", f.Name())
	r.Contains(f.String(), "multi-stage")
}
