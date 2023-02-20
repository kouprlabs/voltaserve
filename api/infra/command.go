package infra

import (
	"os/exec"
	"time"
	"voltaserve/config"
)

type Command struct {
	config config.Config
}

func NewCommand() *Command {
	return &Command{config: config.GetConfig()}
}

func (r *Command) Exec(name string, arg ...string) error {
	timeout := time.Duration(r.config.Limits.ExternalCommandTimeoutSec) * time.Second
	cmd := exec.Command(name, arg...)
	if err := cmd.Start(); err != nil {
		return err
	}
	timer := time.AfterFunc(timeout, func() {
		_ = cmd.Process.Kill()
	})
	if err := cmd.Wait(); err != nil {
		return err
	}
	timer.Stop()
	return nil
}

func (r *Command) ReadOutput(name string, arg ...string) (string, error) {
	timeout := time.Duration(r.config.Limits.ExternalCommandTimeoutSec) * time.Second
	cmd := exec.Command(name, arg...)
	res, err := cmd.Output()
	if err != nil {
		return "", err
	}
	timer := time.AfterFunc(timeout, func() {
		_ = cmd.Process.Kill()
	})
	timer.Stop()
	return string(res), nil
}
