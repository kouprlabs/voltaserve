// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package infra

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/kouprlabs/voltaserve/conversion/config"
	"github.com/kouprlabs/voltaserve/conversion/helper"
)

var commandMutex sync.Mutex

type Command struct {
	config *config.Config
}

func NewCommand() *Command {
	return &Command{config: config.GetConfig()}
}

func (r *Command) RunWithConsoleOutput(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (r *Command) Exec(name string, arg ...string) error {
	commandMutex.Lock()
	defer commandMutex.Unlock()

	timeout := time.Duration(r.config.Limits.ExternalCommandTimeoutSeconds) * time.Second
	cmd := exec.Command(name, arg...)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return err
	}
	timer := time.AfterFunc(timeout, func() {
		_ = cmd.Process.Kill()
	})
	if err := cmd.Wait(); err != nil {
		if stderr.Len() > 0 {
			return errors.New(stderr.String())
		} else {
			return err
		}
	}
	timer.Stop()
	return nil
}

func (r *Command) ReadOutput(name string, arg ...string) (*string, error) {
	commandMutex.Lock()
	defer commandMutex.Unlock()

	timeout := time.Duration(r.config.Limits.ExternalCommandTimeoutSeconds) * time.Second
	cmd := exec.Command(name, arg...)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	res, err := cmd.Output()
	if err != nil {
		if stderr.Len() > 0 {
			return nil, errors.New(stderr.String())
		} else {
			return nil, err
		}
	}
	timer := time.AfterFunc(timeout, func() {
		_ = cmd.Process.Kill()
	})
	timer.Stop()
	return helper.ToPtr(string(res)), nil
}
