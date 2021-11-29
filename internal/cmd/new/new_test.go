//go:build integration
// +build integration

package new

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/gobuffalo/cli/internal/testhelpers"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	r := require.New(t)
	r.NoError(testhelpers.EnsureBuffaloCMD(t))

	dir, err := os.MkdirTemp("", "buffalo-new-test-*")
	r.NoError(err)
	t.Cleanup(func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Logf("failed to delete temporary directory: %s", dir)
		}
	})

	r.NoError(os.Chdir(dir))

	tcases := []struct {
		name  string
		args  []string
		check func(*require.Assertions, string, error)
	}{
		{
			name: "no application name",
			args: []string{"new"},
			check: func(r *require.Assertions, out string, err error) {
				r.Error(err)
				r.Contains(out, "you must enter a name for your new application")
			},
		},
		{
			name: "skip docker",
			args: []string{"new", "nodocker", "--api", "--skip-docker", "-f", "--vcs", "none"},
			check: func(r *require.Assertions, out string, err error) {
				r.NoError(err)
				r.NoFileExists(filepath.Join("nodocker", "Dockerfile"))
			},
		},

		{
			name: "docker there",
			args: []string{"new", "wdocker", "--api", "-f", "--vcs", "none"},
			check: func(r *require.Assertions, out string, err error) {
				r.NoError(err)
				r.FileExists(filepath.Join("wdocker", "Dockerfile"))
			},
		},

		{
			name: "invalid db type",
			args: []string{"new", "api", "--api", "-f", "--db-type", "a"},
			check: func(r *require.Assertions, out string, err error) {
				r.Error(err)
				r.Contains(out, `unknown dialect`)
			},
		},

		{
			name: "forbidden application name",
			args: []string{"new", "buffalo", "-f", "--api"},
			check: func(r *require.Assertions, out string, err error) {
				r.Error(err)
				r.Contains(out, `name buffalo is not allowed, try a different application name`)
			},
		},
	}

	for _, v := range tcases {
		t.Run(v.name, func(t *testing.T) {
			wd, err := os.Getwd()
			r.NoError(err)
			defer os.Chdir(wd)

			dir := os.TempDir()
			r.NoError(os.Chdir(dir))

			r := require.New(t)
			out, err := testhelpers.RunBuffaloCMD(t, v.args)
			v.check(r, out, err)
		})
	}
}

func TestNewAppAPIContent(t *testing.T) {
	r := require.New(t)

	wd, err := os.Getwd()
	r.NoError(err)
	defer os.Chdir(wd)

	t.Log(wd)
	r.NoError(testhelpers.EnsureBuffaloCMD(t))

	dir := t.TempDir()
	r.NoError(os.Chdir(dir))

	_, err = testhelpers.RunBuffaloCMD(t, []string{"new", "apicontent", "--api", "-f", "--vcs", "none"})
	r.NoError(err)

	checks := []struct {
		path  string
		check func(*require.Assertions, string, bool)
	}{
		{
			path: filepath.Join("apicontent", "actions", "app.go"),
			check: func(r *require.Assertions, content string, exists bool) {
				r.True(exists)
				r.Contains(string(content), "app.Use(contenttype.Set(\"application/json\"))")
			},
		},

		{
			path: filepath.Join("apicontent", "actions", "home.go"),
			check: func(r *require.Assertions, content string, exists bool) {
				r.True(exists)

				r.Contains(string(content), "r.JSON")
				r.NotContains(string(content), "r.HTML")
			},
		},

		{
			path: filepath.Join("apicontent", "actions", "render.go"),
			check: func(r *require.Assertions, content string, exists bool) {
				r.True(exists)

				r.Contains(string(content), "DefaultContentType: \"application/json\",")
				r.NotContains(string(content), "HTMLLayout")
				r.NotContains(string(content), "TemplatesBox")
				r.NotContains(string(content), "Helpers")
			},
		},

		{
			path: filepath.Join("apicontent", "templates"),
			check: func(r *require.Assertions, content string, exists bool) {
				r.False(exists)
			},
		},

		{
			path: filepath.Join("apicontent", "assets"),
			check: func(r *require.Assertions, content string, exists bool) {
				r.False(exists)
			},
		},

		{
			path: filepath.Join("apicontent", "public"),
			check: func(r *require.Assertions, content string, exists bool) {
				r.False(exists)
			},
		},
	}

	for _, v := range checks {
		t.Run(v.path, func(t *testing.T) {
			r := require.New(t)
			exists := true

			b, err := os.ReadFile(v.path)
			if err != nil && errors.Is(err, os.ErrNotExist) {
				exists = false
			}

			v.check(r, string(b), exists)
		})
	}

}

func TestNewAppTravis(t *testing.T) {
	r := require.New(t)

	wd, err := os.Getwd()
	r.NoError(err)
	defer os.Chdir(wd)

	t.Log(wd)
	r.NoError(testhelpers.EnsureBuffaloCMD(t))

	dir := t.TempDir()
	r.NoError(os.Chdir(dir))

	out, err := testhelpers.RunBuffaloCMD(t, []string{"new", "apitravis", "--api", "-f", "--vcs", "none", "--ci-provider", "travis", "--db-type", "sqlite3"})
	t.Log(out)
	r.NoError(err)

	checks := []struct {
		path  string
		check func(*require.Assertions, string, bool)
	}{
		{
			path: filepath.Join("apitravis", ".travis.yml"),
			check: func(r *require.Assertions, content string, exists bool) {
				r.True(exists)
				r.Contains(string(content), "language: go")
				r.Contains(string(content), "1.11.x")
				r.Contains(string(content), "go_import_path:")
			},
		},

		{
			path: filepath.Join("apitravis", "database.yml"),
			check: func(r *require.Assertions, content string, exists bool) {
				r.True(exists)
				r.Contains(string(content), "dialect: \"sqlite3\"")
				r.Contains(string(content), "development:")
				r.Contains(string(content), "production:")
				r.Contains(string(content), "test:")
			},
		},
	}

	for _, v := range checks {
		t.Run(v.path, func(t *testing.T) {
			r := require.New(t)
			exists := true

			b, err := os.ReadFile(v.path)
			if err != nil && errors.Is(err, os.ErrNotExist) {
				exists = false
			}

			v.check(r, string(b), exists)
		})
	}
}
