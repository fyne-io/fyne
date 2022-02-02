package commands

import (
	"golang.org/x/sys/execabs"
)

type runner interface {
	RunOutput(arg ...string) ([]byte, error)
	DirSet(dir string)
	EnvSet(env []string)
}

type command struct {
	exe string
	dir *string
	env []string
}

func (c *command) DirSet(dir string) {
	c.dir = &dir
}

func (c *command) EnvSet(env []string) {
	c.env = env
}

func (c *command) RunOutput(arg ...string) ([]byte, error) {
	cmd := execabs.Command(c.exe, arg...)
	if c.dir != nil {
		cmd.Dir = *c.dir
		c.dir = nil
	}
	if c.env != nil {
		cmd.Env = c.env
		c.env = []string{}
	}
	return cmd.CombinedOutput()
}

func runCommandOutput(runner runner, args ...string) ([]byte, error) {
	return runner.RunOutput(args...)
}

func newCommand(cmd string) *command {
	return &command{exe: cmd}
}
