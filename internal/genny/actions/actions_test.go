package actions

import (
	"fmt"
	"io/fs"
	"os"
	"strings"
	"testing"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gentest"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func compare(a, b string) bool {
	a = strings.TrimSpace(a)
	a = strings.Replace(a, "\r", "", -1)
	b = strings.TrimSpace(b)
	b = strings.Replace(b, "\r", "", -1)
	a = strings.TrimSpace(a)
	b = strings.TrimSpace(b)
	res := cmp.Equal(a, b)
	if !res {
		fmt.Println(cmp.Diff(a, b))
	}
	return res
}

func runner(r *require.Assertions) *genny.Runner {
	run := gentest.NewRunner()
	r.NoError(run.Disk.AddFS(os.DirFS("../actions/_fixtures/inputs/clean")))
	return run
}

func Test_New(t *testing.T) {
	r := require.New(t)

	g, err := New(&Options{
		Name:    "user",
		Actions: []string{"index"},
	})
	r.NoError(err)

	run := runner(r)
	r.NoError(run.With(g))
	r.NoError(run.Run())

	res := run.Results()
	r.Len(res.Commands, 0)
	r.Len(res.Files, 4)

	fsys := os.DirFS("../actions/_fixtures/outputs/clean")
	files := []string{"actions/user.go.tmpl", "actions/app.go.tmpl", "actions/user_test.go.tmpl", "templates/user/index.plush.html"}
	for _, s := range files {
		x, err := fs.ReadFile(fsys, s)
		r.NoError(err)
		f, err := res.Find(strings.TrimSuffix(s, ".tmpl"))
		r.NoError(err)
		fmt.Printf("\nfile %s", s)
		r.True(compare(string(x), f.String()))
	}
}

func Test_New_Multi(t *testing.T) {
	r := require.New(t)

	g, err := New(&Options{
		Name:    "user",
		Actions: []string{"show", "edit"},
	})
	r.NoError(err)

	run := runner(r)
	r.NoError(run.With(g))
	r.NoError(run.Run())

	res := run.Results()
	r.Len(res.Commands, 0)

	fsys := os.DirFS("../actions/_fixtures/outputs/multi")
	files := []string{"actions/user.go.tmpl", "actions/app.go.tmpl", "actions/user_test.go.tmpl", "templates/user/show.plush.html", "templates/user/edit.plush.html"}
	for _, s := range files {
		x, err := fs.ReadFile(fsys, s)
		r.NoError(err)
		f, err := res.Find(strings.TrimSuffix(s, ".tmpl"))
		r.NoError(err)
		fmt.Printf("\nfile %s\n", f.Name())
		r.True(compare(string(x), f.String()))
	}
}

func Test_New_Multi_Existing(t *testing.T) {
	r := require.New(t)

	g, err := New(&Options{
		Name:    "user",
		Actions: []string{"show", "edit"},
	})
	r.NoError(err)

	run := gentest.NewRunner()
	ins := os.DirFS("../actions/_fixtures/inputs/existing")
	err = fs.WalkDir(ins, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		f, err := ins.Open(path)
		r.NoError(err)
		path = strings.TrimSuffix(path, ".tmpl")
		run.Disk.Add(genny.NewFile(path, f))
		return nil
	})
	r.NoError(err)
	r.NoError(run.With(g))
	r.NoError(run.Run())

	res := run.Results()
	r.Len(res.Commands, 0)

	fsys := os.DirFS("../actions/_fixtures/outputs/existing")
	files := []string{"actions/user.go.tmpl", "actions/app.go.tmpl", "actions/user_test.go.tmpl", "templates/user/show.plush.html", "templates/user/edit.plush.html"}
	for _, s := range files {
		x, err := fs.ReadFile(fsys, s)
		r.NoError(err)
		f, err := res.Find(strings.TrimSuffix(s, ".tmpl"))
		r.NoError(err)
		r.True(compare(string(x), f.String()))
	}
}

func Test_New_SkipTemplates(t *testing.T) {
	r := require.New(t)

	g, err := New(&Options{
		Name:          "user",
		Actions:       []string{"index"},
		SkipTemplates: true,
	})
	r.NoError(err)

	run := runner(r)
	r.NoError(run.With(g))
	r.NoError(run.Run())

	res := run.Results()
	r.Len(res.Commands, 0)

	files := []string{"templates/user/index.html"}
	for _, s := range files {
		_, err := res.Find(s)
		r.Error(err)
	}
}
