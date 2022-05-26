// Copyright 2015 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mobile

import (
	"flag"
	"fmt"
	"os"
)

var (
	gomobileName = "gomobile"
)

type command struct {
	run   func(*command) error
	Flag  flag.FlagSet
	Name  string
	Usage string
	Short string
	Long  string

	IconPath, AppName      string
	Version, Cert, Profile string
	Build                  int
}

func (cmd *command) usage() {
	fmt.Fprintf(os.Stdout, "usage: %s %s %s\n%s", gomobileName, cmd.Name, cmd.Usage, cmd.Long)
}
