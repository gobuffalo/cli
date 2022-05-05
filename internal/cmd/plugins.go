package cmd

import (
	"sync"

	"github.com/gobuffalo/cli/internal/plugins"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var _plugs plugins.List
var initPlugsOnce sync.Once

func plugs() plugins.List {
	initPlugsOnce.Do(func() {
		var err error
		_plugs, err = plugins.Available()
		if err == nil {
			return
		}

		_plugs = plugins.List{}
		logrus.Errorf("error loading plugins %s", err)
	})
	return _plugs
}

func decorate(name string, cmd *cobra.Command) {
	pugs := plugs()
	for _, c := range pugs[name] {
		anywhereCommands = append(anywhereCommands, c.Name)
		cc := plugins.Decorate(c)
		cmd.AddCommand(cc)
	}
}
