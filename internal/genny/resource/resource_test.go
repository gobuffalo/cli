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
	"github.com/gobuffalo/cli/internal/genny/newapp/api"
	"github.com/gobuffalo/cli/internal/genny/newapp/core"
	"github.com/gobuffalo/cli/internal/genny/newapp/web"
	"github.com/gobuffalo/cli/internal/genny/testrunner"
	"github.com/gobuffalo/flect/name"
	"github.com/gobuffalo/meta"
	"github.com/stretchr/testify/require"
)

type pass struct {
	Name    string
	Options Options
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
			r := require.New(st)

			opts := &web.Options{}
			opts.Options = &core.Options{App: app}
			run, err := testrunner.WebApp(opts)
			r.NoError(err)

			tt.Options.App = app
			g, err := New(&tt.Options)
			r.NoError(err)

			r.NoError(run.With(g))
			r.NoError(run.Run())

			res := run.Results()
			r.Len(res.Commands, 1)

			c := res.Commands[0]
			r.Equal("buffalo-pop pop g model widget name desc:nulls.Text", strings.Join(c.Args, " "))
			r.Len(res.Files, 30)

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
			r := require.New(st)

			opts := &web.Options{}
			opts.Options = &core.Options{App: app}
			run, err := testrunner.WebApp(opts)
			r.NoError(err)

			tt.Options.App = app
			tt.Options.SkipTemplates = true
			g, err := New(&tt.Options)
			r.NoError(err)

			r.NoError(run.With(g))
			r.NoError(run.Run())

			res := run.Results()

			r.Len(res.Commands, 1)
			for _, s := range []string{"_form", "edit", "index", "new", "show"} {
				p := path.Join("templates", tt.Name, s+".html")
				_, err = res.Find(p)
				r.Error(err)
			}

			r.Len(res.Files, 24)
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
			r := require.New(st)

			opts := &api.Options{}
			opts.Options = &core.Options{App: app}
			run, err := testrunner.ApiApp(opts)
			r.NoError(err)

			tt.Options.App = app
			g, err := New(&tt.Options)
			r.NoError(err)

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

			r.Len(res.Files, 18)
		})
	}
}

func Test_New_UseModel(t *testing.T) {
	r := require.New(t)

	app := meta.New(".")
	app.PackageRoot("github.com/markbates/coke")

	opts := &web.Options{}
	opts.Options = &core.Options{App: app}
	run, err := testrunner.WebApp(opts)
	r.NoError(err)

	ats, err := attrs.ParseArgs("name", "desc:nulls.Text")
	r.NoError(err)

	g, err := New(&Options{
		App:   app,
		Name:  "Widget",
		Model: "User",
		Attrs: ats,
	})
	r.NoError(err)
	r.NoError(run.With(g))
	r.NoError(run.Run())

	res := run.Results()
	r.Len(res.Commands, 1)

	c := res.Commands[0]
	r.Equal("buffalo-pop pop g model user name desc:nulls.Text", strings.Join(c.Args, " "))
	r.Len(res.Files, 30)

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

	opts := &web.Options{}
	opts.Options = &core.Options{App: app}
	run, err := testrunner.WebApp(opts)
	r.NoError(err)

	g, err := New(&Options{
		App:       app,
		Name:      "Widget",
		SkipModel: true,
	})
	r.NoError(err)
	r.NoError(run.With(g))
	r.NoError(run.Run())

	res := run.Results()
	r.Len(res.Commands, 0)
	r.Len(res.Files, 30)

	f, err := res.Find("actions/widgets.go")
	r.NoError(err)
	actions := []string{"List", "Show", "Create", "Update", "Destroy", "New", "Edit"}
	for _, action := range actions {
		r.Contains(f.String(), fmt.Sprintf("func (v WidgetsResource) %v(c buffalo.Context) error {", action))
	}
}
