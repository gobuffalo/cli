package templates

import (
	"embed"
	"io/fs"
)

//go:embed * _flash.plush.html
var files embed.FS

func FS() fs.FS {
	return shadowFS{files}
}

type shadowFS struct {
	embed.FS
}

func (f shadowFS) Open(name string) (fs.File, error) {
	if name == "embed.go" {
		return nil, fs.ErrNotExist
	}
	return f.FS.Open(name)
}

func (f shadowFS) ReadDir(name string) ([]fs.DirEntry, error) {
	entries, err := f.FS.ReadDir(name)
	if name == "." {
		for i, entry := range entries {
			if entry.Name() == "embed.go" {
				entries = append(entries[:i], entries[i+1:]...)
				break
			}
		}
	}
	return entries, err
}

func (f shadowFS) ReadFile(name string) ([]byte, error) {
	if name == "embed.go" {
		return nil, fs.ErrNotExist
	}
	return f.FS.ReadFile(name)
}
