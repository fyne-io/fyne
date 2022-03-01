package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// AndroidBuildToolsPath tries to find the location of the "build-tools" directory.
// This depends on ANDROID_HOME variable, so you should call RequireAndroidSDK() first.
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
	return strings.HasPrefix(os, "android")
}

// IsIOS returns true if the given os parameter represents one of the iOS targets (ios, iossimulator)
func IsIOS(os string) bool {
	return strings.HasPrefix(os, "ios")
}

// IsMobile returns true if the given os parameter represents a platform handled by gomobile.
func IsMobile(os string) bool {
	return IsIOS(os) || IsAndroid(os)
}

// RequireAndroidSDK will return an error if it cannot establish the location of a valid Android SDK installation.
// This is currently deduced using ANDROID_HOME environment variable.
func RequireAndroidSDK() error {
	if env, ok := os.LookupEnv("ANDROID_HOME"); !ok || env == "" {
		return fmt.Errorf("could not find android tools, missing ANDROID_HOME")
	}

	return nil
}
