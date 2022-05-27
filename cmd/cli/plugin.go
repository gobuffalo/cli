package cli

type Plugin interface {
	Name() string
}

type Plugins []Plugin

func (pl Plugins) FindCommand(name string) Command {
	for _, p := range pl {
		cmd, ok := p.(Command)
		if !ok || p.Name() != name {
			continue
		}

		return cmd
	}

	return nil
}
