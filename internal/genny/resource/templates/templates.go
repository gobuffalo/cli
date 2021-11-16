package templates

import (
	"embed"
	"io/fs"
)

//go:embed core/* core/templates/folder-name/_form.plush.html.tmpl
var core embed.FS

func Core() fs.FS {
	sub, err := fs.Sub(core, "core")
	if err != nil {
		return nil
	}
	return sub
}

//go:embed standard/*
var standard embed.FS

func Standard() fs.FS {
	sub, err := fs.Sub(standard, "standard")
	if err != nil {
		return nil
	}
	return sub
}

//go:embed use_model/*
var useModel embed.FS

func UseModel() fs.FS {
	sub, err := fs.Sub(useModel, "use_model")
	if err != nil {
		return nil
	}
	return sub
}
