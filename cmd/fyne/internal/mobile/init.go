// Copyright 2015 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mobile

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/sys/execabs"
)

var (
	goos = runtime.GOOS
)

func mkdir(dir string) error {
	if buildX || buildN {
		printcmd("mkdir -p %s", dir)
	}
	if buildN {
		return nil
	}
	return os.MkdirAll(dir, 0750)
}

func removeAll(path string) error {
	if buildX || buildN {
		printcmd(`rm -r -f "%s"`, path)
	}
	if buildN {
		return nil
	}

	// os.RemoveAll behaves differently in windows.
	// http://golang.org/issues/9606
	if goos == "windows" {
		err := resetReadOnlyFlagAll(path)
		if err != nil {
			return err
		}
	}

	return os.RemoveAll(path)
}

func resetReadOnlyFlagAll(path string) error {
	fi, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return os.Chmod(path, 0600)
	}
	fd, err := os.Open(filepath.Clean(path))
	if err != nil {
		return err
	}
	defer fd.Close()

	names, _ := fd.Readdirnames(-1)
	for _, name := range names {
		err := resetReadOnlyFlagAll(path + string(filepath.Separator) + name)
		if err != nil {
			return err
		}
	}
	return nil
}

func goEnv(name string) string {
	if val := os.Getenv(name); val != "" {
		return val
	}
	val, err := execabs.Command("go", "env", name).Output()
	if err != nil {
		panic(err) // the Go tool was tested to work earlier
	}
	return strings.TrimSpace(string(val))
}

func runCmd(cmd *execabs.Cmd) error {
	if buildX || buildN {
		dir := ""
		if cmd.Dir != "" {
			dir = "PWD=" + cmd.Dir + " "
		}
		env := strings.Join(cmd.Env, " ")
		if env != "" {
			env += " "
		}
		printcmd("%s%s%s", dir, env, strings.Join(cmd.Args, " "))
	}

	buf := new(bytes.Buffer)
	buf.WriteByte('\n')
	if buildV {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		cmd.Stdout = buf
		cmd.Stderr = buf
	}

	if buildWork {
		if goos == "windows" {
			cmd.Env = append(cmd.Env, `TEMP=`+tmpdir, `TMP=`+tmpdir)
		} else {
			cmd.Env = append(cmd.Env, `TMPDIR=`+tmpdir)
		}
	}

	if !buildN {
		cmd.Env = environ(cmd.Env)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("%s failed: %v%s", strings.Join(cmd.Args, " "), err, buf)
		}
	}
	return nil
}
