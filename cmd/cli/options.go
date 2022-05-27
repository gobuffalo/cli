package cli

import flag "github.com/spf13/pflag"

type Option func(*App)

// New cli application
func New(opts ...Option) *App {
	app := &App{
		IO:      &IO{},
		FlagSet: flag.NewFlagSet("application", flag.ContinueOnError),

		plugins: Plugins{},
	}

	for _, o := range opts {
		o(app)
	}

	return app
}

// WithPlugins sets the plugins for the app.
func WithPlugins(plugins ...Plugin) Option {
	return func(a *App) {
		a.plugins = append(a.plugins, plugins...)
	}
}

// WithIO sets the IO for the app.
func WithIO(io *IO) Option {
	return func(a *App) {
		a.IO = io
	}
}
