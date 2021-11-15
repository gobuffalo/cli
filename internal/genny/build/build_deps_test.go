package build

import (
	"testing"

	"github.com/gobuffalo/genny/v2/gentest"
	"github.com/gobuffalo/meta"
	"github.com/stretchr/testify/require"
)

func Test_buildDeps(t *testing.T) {
	r := require.New(t)

	opts := &Options{
		Tags: meta.BuildTags{"foo"},
	}

	run := gentest.NewRunner()
	r.NoError(run.WithNew(buildDeps(opts)))

	r.NoError(run.Run())

	res := run.Results()
	r.Len(res.Commands, 0)
}
