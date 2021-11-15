package testtemplate

import (
	"embed"
	"io/fs"
)

//go:embed config/*
var config embed.FS

func Config() fs.FS {
	sub, err := fs.Sub(config, "config")
	if err != nil {
		return nil
	}
	return sub
}

//go:embed module/*
var module embed.FS

func Module() fs.FS {
	sub, err := fs.Sub(module, "module")
	if err != nil {
		return nil
	}
	return modBakRename{sub}
}

type modBakRename struct {
	fs.FS
}

func (m modBakRename) Open(name string) (fs.File, error) {
	if name == "go.mod" {
		name = "go.mod.bak"
	}
	return m.FS.Open(name)
}
