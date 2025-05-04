package utils

import (
	"os/exec"
)

type Command struct {
	Name  string
	Probe string
}

func (c *Command) Exists() bool {
	_, err := exec.LookPath(c.Name)
	if err != nil {
		return false
	}

	// check that the command is installed with the probe
	// this may be hacky, but my enviroment is wonky due to pyenv and `ffmpeg` returns that it exists because it is in a pyenv/shims directory even though it is not installed
	if c.Probe != "" {
		cmd := exec.Command(c.Name, c.Probe)
		return cmd.Run() == nil
	}

	return true
}
