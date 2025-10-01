//go:build !ci && android && !noos && !tinygo

package app

import (
	"log"
	"os"
	"path/filepath"
)

func rootConfigDir() string {
	filesDir := os.Getenv("FILESDIR")
	if filesDir == "" {
		log.Println("FILESDIR env was not set by android native code")
		return "/data/data" // probably won't work, but we can't make a better guess
	}

	return filepath.Join(filesDir, "fyne")
}
