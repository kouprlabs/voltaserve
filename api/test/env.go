// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package test

import "os"

type Env struct{}

func NewEnv() *Env {
	return &Env{}
}

func (e *Env) Apply() error {
	if err := os.Setenv("TEST", "true"); err != nil {
		return err
	}
	if err := os.Setenv("LIMITS_FILE_PROCESSING_MB", "video:10000,*:1000"); err != nil {
		return err
	}
	if err := os.Setenv("DEFAULTS_WORKSPACE_STORAGE_CAPACITY_MB", "100000"); err != nil {
		return err
	}
	return nil
}
