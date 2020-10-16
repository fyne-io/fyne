package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func AndroidBuildToolsPath() string {
	env := os.Getenv("ANDROID_HOME")
	dir := filepath.Join(env, "build-tools")

	// this may contain a version subdir
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return dir
	}

	childDir := ""
	for _, f := range files {
		if f.Name() == "zipalign" { // no version subdir
			return dir
		}
		if f.IsDir() && childDir == "" {
			childDir = f.Name()
		}
	}

	if childDir == "" {
		return dir
	}
	return filepath.Join(dir, childDir)
}

// IsAndroid returns true if the given os parameter represents one of the Android targets.
func IsAndroid(os string) bool {
	return strings.Index(os, "android") == 0
}

// IsMobile returns true if the given os parameter represents a platform handled by gomobile.
func IsMobile(os string) bool {
	return os == "ios" || strings.HasPrefix(os, "android")
}

func RequireAndroidSDK() error {
	if env, ok := os.LookupEnv("ANDROID_HOME"); !ok || env == "" {
		return fmt.Errorf("could not find android tools, missing ANDROID_HOME")
	}

	return nil
}
