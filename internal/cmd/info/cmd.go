package info

import (
	"github.com/gobuffalo/clara/v2/genny/rx"
	"github.com/gobuffalo/cli/internal/genny/info"
	"github.com/gobuffalo/meta"
	"github.com/spf13/cobra"
)

var (
	app         = meta.New(".")
	infoOptions = struct {
		Clara *rx.Options
		Info  *info.Options
	}{
		Clara: &rx.Options{
			App: app,
		},
		Info: &info.Options{
			App: app,
		},
	}
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info",
		Short: "Print diagnostic information (useful for debugging)",
		RunE:  runE,
	}

	return cmd
}
