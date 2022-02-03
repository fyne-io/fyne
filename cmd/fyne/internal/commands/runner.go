package commands

import (
	"golang.org/x/sys/execabs"
)

type runner interface {
	runOutput(arg ...string) ([]byte, error)
	setDir(dir string)
	setEnv(env []string)
}

type command struct {
	exe string
	dir string
	env []string
}

func (c *command) setDir(dir string) {
	c.dir = dir
}

func (c *command) setEnv(env []string) {
	c.env = env
}

func (c *command) runOutput(arg ...string) ([]byte, error) {
	cmd := execabs.Command(c.exe, arg...)
	cmd.Dir = c.dir
	if c.env != nil {
		cmd.Env = c.env
		c.env = []string{}
	}
	return cmd.CombinedOutput()
}

func newCommand(cmd string) *command {
	return &command{exe: cmd}
}
