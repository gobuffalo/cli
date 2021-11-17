package resource

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gobuffalo/attrs"
	"github.com/gobuffalo/flect/name"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gentest"
	"github.com/gobuffalo/meta"
	"github.com/stretchr/testify/require"
)

type pass struct {
	Name    string
	Options Options
}

func runner(r *require.Assertions) *genny.Runner {
	run := gentest.NewRunner()
	fsys := os.DirFS("./_fixtures/coke")
	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		f, err := fsys.Open(path)
		if err != nil {
			return err
		}
		path = strings.TrimSuffix(path, ".tmpl")
		run.Disk.Add(genny.NewFile(path, f))
		return nil
	})
	r.NoError(err)
	return run
}

func Test_New(t *testing.T) {
	ats, err := attrs.ParseArgs("name", "desc:nulls.Text")
	if err != nil {
		t.Fatal(err)
	}

	table := []pass{
		{"default", Options{Name: "widget", Attrs: ats}},
		{"nested", Options{Name: "admin/widget", Attrs: ats}},
	}

	app := meta.New(".")
	app.PackageRoot("github.com/markbates/coke")

	for _, tt := range table {
		t.Run(tt.Name, func(st *testing.T) {
			tt.Options.App = app
			r := require.New(st)
			g, err := New(&tt.Options)
			r.NoError(err)

			run := runner(r)
			r.NoError(run.With(g))
			r.NoError(run.Run())

			res := run.Results()
			r.Len(res.Commands, 1)

			c := res.Commands[0]
			r.Equal("buffalo-pop pop g model widget name desc:nulls.Text", strings.Join(c.Args, " "))

			r.Len(res.Files, 9)

			nn := name.New(tt.Options.Name).Pluralize().String()
			actions := []string{"_form", "index", "show", "new", "edit"}
			for _, s := range actions {
				p := path.Join("templates", nn, s+".plush.html")
				_, err = res.Find(p)
				r.NoError(err)
			}

			fsys := os.DirFS(filepath.Join("_fixtures", tt.Name))
			err = fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if d.IsDir() {
					return nil
				}

				f, err := res.Find(strings.TrimSuffix(path, ".tmpl"))
				r.NoError(err)

				s, err := fs.ReadFile(fsys, path)
				r.NoError(err)

				clean := func(s string) string {
					s = strings.TrimSpace(s)
					s = strings.Replace(s, "\n", "", -1)
					s = strings.Replace(s, "\t", "", -1)
					s = strings.Replace(s, "\r", "", -1)
					return s
				}
				r.Equal(clean(string(s)), clean(f.String()))
				return nil
			})
			r.NoError(err)
		})
	}
}

func Test_New_SkipTemplates(t *testing.T) {
	ats, err := attrs.ParseArgs("name", "desc:nulls.Text")
	if err != nil {
		t.Fatal(err)
	}
	table := []pass{
		{"default", Options{Name: "widget", Attrs: ats}},
		{"nested", Options{Name: "admin/widget", Attrs: ats}},
	}

	app := meta.New(".")
	app.PackageRoot("github.com/markbates/coke")

	for _, tt := range table {
		t.Run(tt.Name, func(st *testing.T) {
			tt.Options.App = app
			tt.Options.SkipTemplates = true
			r := require.New(st)
			g, err := New(&tt.Options)
			r.NoError(err)

			run := runner(r)
			r.NoError(run.With(g))
			r.NoError(run.Run())

			res := run.Results()

			r.Len(res.Commands, 1)

			for _, s := range []string{"_form", "edit", "index", "new", "show"} {
				p := path.Join("templates", tt.Name, s+".html")
				_, err = res.Find(p)
				r.Error(err)
			}

			r.Len(res.Files, 3)
		})
	}
}

func Test_New_API(t *testing.T) {
	ats, err := attrs.ParseArgs("name", "desc:nulls.Text")
	if err != nil {
		t.Fatal(err)
	}
	table := []pass{
		{"default", Options{Name: "widget", Attrs: ats}},
		{"nested", Options{Name: "admin/widget", Attrs: ats}},
	}

	app := meta.New(".")
	app.PackageRoot("github.com/markbates/coke")
	app.AsAPI = true

	for _, tt := range table {
		t.Run(tt.Name, func(st *testing.T) {
			tt.Options.App = app
			r := require.New(st)
			g, err := New(&tt.Options)
			r.NoError(err)

			run := runner(r)
			r.NoError(run.With(g))
			r.NoError(run.Run())

			res := run.Results()

			r.Len(res.Commands, 1)

			nn := name.New(tt.Options.Name).Pluralize().String()
			for _, s := range []string{"_form", "edit", "index", "new", "show"} {
				p := path.Join("templates", nn, s+".html")
				_, err = res.Find(p)
				r.Error(err)
			}

			r.Len(res.Files, 3)
		})
	}
}

func Test_New_UseModel(t *testing.T) {
	r := require.New(t)

	ats, err := attrs.ParseArgs("name", "desc:nulls.Text")
	r.NoError(err)

	app := meta.New(".")
	app.PackageRoot("github.com/markbates/coke")

	opts := &Options{
		App:   app,
		Name:  "Widget",
		Model: "User",
		Attrs: ats,
	}
	g, err := New(opts)
	r.NoError(err)

	run := runner(r)
	r.NoError(run.With(g))
	r.NoError(run.Run())

	res := run.Results()
	r.Len(res.Commands, 1)

	c := res.Commands[0]
	r.Equal("buffalo-pop pop g model user name desc:nulls.Text", strings.Join(c.Args, " "))

	r.Len(res.Files, 9)

	for _, s := range []string{"_form", "edit", "index", "new", "show"} {
		p := path.Join("templates", "widgets", s+".plush.html")
		_, err = res.Find(p)
		r.NoError(err)
	}

	f, err := res.Find("actions/widgets.go")
	r.NoError(err)
	r.Contains(f.String(), "users := &models.Users{}")

}

func Test_New_SkipModel(t *testing.T) {
	r := require.New(t)

	app := meta.New(".")
	app.PackageRoot("github.com/markbates/coke")

	opts := &Options{
		App:       app,
		Name:      "Widget",
		SkipModel: true,
	}

	g, err := New(opts)
	r.NoError(err)

	run := runner(r)
	r.NoError(run.With(g))
	r.NoError(run.Run())

	res := run.Results()
	r.Len(res.Commands, 0)
	r.Len(res.Files, 9)

	f, err := res.Find("actions/widgets.go")
	r.NoError(err)
	actions := []string{"List", "Show", "Create", "Update", "Destroy", "New", "Edit"}
	for _, action := range actions {
		r.Contains(f.String(), fmt.Sprintf("func (v WidgetsResource) %v(c buffalo.Context) error {", action))
	}
}
