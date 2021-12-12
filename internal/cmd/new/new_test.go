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

	tt := []struct {
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

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			testhelpers.RunWithinTempFolder(t, func(t *testing.T) {
				r := require.New(t)
				out, err := testhelpers.RunBuffaloCMD(t, tc.args)
				tc.check(r, out, err)
			})
		})
	}
}

func TestNewAppAPIContent(t *testing.T) {
	r := require.New(t)
	r.NoError(testhelpers.EnsureBuffaloCMD(t))

	tt := []struct {
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

	testhelpers.RunWithinTempFolder(t, func(t *testing.T) {
		r := require.New(t)
		_, err := testhelpers.RunBuffaloCMD(t, []string{"new", "apicontent", "--api", "-f", "--vcs", "none"})
		r.NoError(err)

		for _, tc := range tt {
			t.Run(tc.path, func(t *testing.T) {
				r := require.New(t)
				exists := true

				b, err := os.ReadFile(tc.path)
				if err != nil && errors.Is(err, os.ErrNotExist) {
					exists = false
				}

				tc.check(r, string(b), exists)
			})
		}
	})
}

func TestNewAppTravis(t *testing.T) {
	r := require.New(t)
	r.NoError(testhelpers.EnsureBuffaloCMD(t))

	tt := []struct {
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

	testhelpers.RunWithinTempFolder(t, func(t *testing.T) {
		out, err := testhelpers.RunBuffaloCMD(t, []string{"new", "apitravis", "--api", "-f", "--vcs", "none", "--ci-provider", "travis", "--db-type", "sqlite3"})
		t.Log(out)
		r.NoError(err)

		r := require.New(t)
		for _, tc := range tt {
			t.Run(tc.path, func(t *testing.T) {
				b, err := os.ReadFile(tc.path)

				exists := true
				if err != nil && errors.Is(err, os.ErrNotExist) {
					exists = false
				}

				tc.check(r, string(b), exists)
			})
		}
	})
}
