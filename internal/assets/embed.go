package assets

import (
	"embed"
	"io/fs"
)

//go:embed dist/*
var embeddedFiles embed.FS

func GetFileSystem() fs.FS {
	fsys, err := fs.Sub(embeddedFiles, "dist")
	if err != nil {
		panic(err)
	}
	return fsys
}
