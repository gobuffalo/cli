package cli

import (
	"context"
)

// Command is a struct that can be invoked by the main app
// to do so it would use the Name method to identify it.
type Command interface {
	Name() string
	Main(ctx context.Context, pwd string, args []string) error
}

// Commands holds a set of useful methods for working with
// a group of commands.
type Commands []Command

func (cc Commands) Contains(name string) bool {
	return cc.Find(name) != nil
}

func (cc Commands) Find(name string) Command {
	for _, v := range cc {
		if v.Name() == name {
			return v
		}
	}

	return nil
}
