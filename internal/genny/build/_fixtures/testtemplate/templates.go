package testtemplate

import (
	"embed"
	"io/fs"
)

//go:embed good/*
var good embed.FS

func Good() fs.FS {
	sub, err := fs.Sub(good, "good")
	if err != nil {
		return nil
	}
	return sub
}

//go:embed bad/*
var bad embed.FS

func Bad() fs.FS {
	sub, err := fs.Sub(bad, "bad")
	if err != nil {
		return nil
	}
	return sub
}
