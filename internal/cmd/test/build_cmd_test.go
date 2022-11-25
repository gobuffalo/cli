package test

import (
	"strings"
	"testing"
)

func TestBuildCmd(t *testing.T) {
	tcases := []struct {
		name     string
		args     []string
		want     []string
		contains string
	}{
		{
			name: "no args",
			args: []string{},
			want: []string{"go", "test", "-p", "1", "-tags", "development"},
		},

		{
			name: "path specified",
			args: []string{"./cmd/..."},
			want: []string{"go", "test", "-p", "1", "-tags", "development", "./cmd/..."},
		},

		{
			name: "path specified and flags",
			args: []string{"--count=1", "./cmd/..."},
			want: []string{"go", "test", "-p", "1", "-tags", "development", "--count=1", "./cmd/..."},
		},

		{
			name: "flags but no path",
			args: []string{"--count=1", "-v"},
			want: []string{"go", "test", "-p", "1", "-tags", "development", "--count=1", "-v"},
		},

		{
			name: "flags but no path",
			args: []string{"--count=1", "-v"},
			want: []string{"go", "test", "-p", "1", "-tags", "development", "--count=1", "-v"},
		},

		{
			name: "flags no path with -tags",
			args: []string{"--count=1", "-tags", "foo,bar"},
			want: []string{"go", "test", "-p", "1", "-tags", "development", "--count=1", "-tags", "foo,bar"},
		},

		{
			name: "force migrations should be removed",
			args: []string{"--force-migrations", "-tags", "foo,bar"},
			want: []string{"go", "test", "-p", "1", "-tags", "development", "-tags", "foo,bar"},
		},

		{
			name: "force migrations should be removed",
			args: []string{"--force-migrations", "-m", "Something"},
			want: []string{"go", "test", "-p", "1", "-tags", "development"},
		},

		{
			name:     "testify.m should go at the end",
			args:     []string{"--force-migrations", "-testify.m", "Something"},
			want:     []string{"go", "test", "-p", "1", "-tags", "development"},
			contains: "-testify.m Something",
		},
	}

	for _, tc := range tcases {
		t.Run(tc.name, func(t *testing.T) {
			cmd, err := buildCmd(tc.args)
			if err != nil {
				t.Fatal(err)
			}

			wantcmd := strings.Join(tc.want, " ")
			resultcmd := strings.Join(cmd.Args, " ")

			prefixMatches := strings.HasPrefix(resultcmd, wantcmd)
			if !prefixMatches {
				t.Errorf("prefix `%s` not found in `%s`", wantcmd, resultcmd)
			}

			if tc.contains != "" && !strings.Contains(resultcmd, tc.contains) {
				t.Errorf("string `%s` not found in `%s`", tc.contains, resultcmd)
			}
		})
	}
}
