//go:build dev
// +build dev

package main

import (
	"os"
)

func init() {
	// In dev mode, serve from the file system
	frontendFS = os.DirFS("frontend/dist")
}
