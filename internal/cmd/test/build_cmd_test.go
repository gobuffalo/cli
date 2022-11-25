package test

import (
	"reflect"
	"testing"
)

func TestBuildCmd(t *testing.T) {
	tcases := []struct {
		name string
		args []string
		want []string
	}{
		{
			name: "no args",
			args: []string{},
			want: []string{"go", "test", "-p", "1", "-tags", "development", "./..."},
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
			want: []string{"go", "test", "-p", "1", "-tags", "development", "--count=1", "-v", "./..."},
		},

		{
			name: "flags but no path",
			args: []string{"--count=1", "-v"},
			want: []string{"go", "test", "-p", "1", "-tags", "development", "--count=1", "-v", "./..."},
		},

		{
			name: "flags no path with -tags",
			args: []string{"--count=1", "-tags", "foo,bar"},
			want: []string{"go", "test", "-p", "1", "-tags", "development", "--count=1", "-tags", "foo,bar", "./..."},
		},

		{
			name: "force migrations should be removed",
			args: []string{"--force-migrations", "-tags", "foo,bar"},
			want: []string{"go", "test", "-p", "1", "-tags", "development", "-tags", "foo,bar", "./..."},
		},
	}

	for _, tc := range tcases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := buildCmd(tc.args)
			if len(cmd.Args) != len(tc.want) {
				t.Errorf("want %d args, got %d", len(tc.want), len(cmd.Args))
			}

			if !reflect.DeepEqual(cmd.Args, tc.want) {
				t.Errorf("want %s, got %s", tc.want, cmd.Args)
			}
		})
	}
}
