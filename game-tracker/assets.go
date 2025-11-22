package gametracker

import (
	"embed"
	"io/fs"
)

//go:embed frontend/dist/*
var frontendFS embed.FS

func GetFrontendAssets() fs.FS {
	f, err := fs.Sub(frontendFS, "frontend/dist")
	if err != nil {
		panic(err)
	}
	return f
}
