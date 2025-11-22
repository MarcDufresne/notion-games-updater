//go:build !dev
// +build !dev

package main

import (
	"embed"
)

//go:embed frontend/dist
var embeddedFrontend embed.FS

func init() {
	frontendFS = embeddedFrontend
}
