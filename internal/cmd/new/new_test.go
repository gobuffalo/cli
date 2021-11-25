package new

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	dir := os.TempDir()
	err := os.Chdir(dir)
	if err != nil {
		t.Fatal("error changing to tmp dir")
	}

	err = os.RemoveAll(filepath.Join(dir, "app"))
	if err != nil {
		t.Fatal("error removing new dir")
	}

	tcases := []struct {
		name  string
		args  []string
		check func() error
	}{
		{
			name: "skip docker",
			args: []string{"app", "--api", "--skip-docker", "-f", "--vcs", "none"},
			check: func() error {
				_, err := os.Stat(filepath.Join("app", "Dockerfile"))
				if err == nil {
					return fmt.Errorf("Dockerfile should not be generated")
				}

				return nil
			},
		},

		{
			name: "docker there",
			args: []string{"app", "--api", "-f", "--vcs", "none"},
			check: func() error {
				_, err := os.Stat(filepath.Join("app", "Dockerfile"))
				if err != nil {
					return fmt.Errorf("Dockerfile should be there")
				}

				return nil
			},
		},
	}

	for _, v := range tcases {
		t.Run(v.name, func(t *testing.T) {
			Cmd.SetArgs(v.args)
			err = Cmd.Execute()
			if err != nil {
				t.Fatal("error running new")
			}

			err = v.check()
			if err != nil {
				t.Fatal(err)
			}
		})
	}

}
