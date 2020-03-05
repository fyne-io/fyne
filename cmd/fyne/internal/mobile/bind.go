// Copyright 2015 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mobile

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
)

// pkg, tmpdir in build.go

func copyFile(dst, src string) error {
	if buildX {
		printcmd("cp %s %s", src, dst)
	}
	return writeFile(dst, func(w io.Writer) error {
		if buildN {
			return nil
		}
		f, err := os.Open(src)
		if err != nil {
			return err
		}
		defer f.Close()

		if _, err := io.Copy(w, f); err != nil {
			return fmt.Errorf("cp %s %s failed: %v", src, dst, err)
		}
		return nil
	})
}

func writeFile(filename string, generate func(io.Writer) error) error {
	if buildV {
		fmt.Fprintf(os.Stderr, "write %s\n", filename)
	}

	err := mkdir(filepath.Dir(filename))
	if err != nil {
		return err
	}

	if buildN {
		return generate(ioutil.Discard)
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := f.Close(); err == nil {
			err = cerr
		}
	}()

	return generate(f)
}

func packagesConfig(targetOS string) *packages.Config {
	config := &packages.Config{}
	config.Env = append(os.Environ(), "GOARCH=arm", "GOOS="+targetOS)
	tags := buildTags
	if targetOS == "darwin" {
		tags = append(tags, "ios")
	}
	if len(tags) > 0 {
		config.BuildFlags = []string{"-tags=" + strings.Join(tags, ",")}
	}
	return config
}
