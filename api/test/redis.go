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

import (
	"os"

	"github.com/alicebob/miniredis/v2"
)

type Redis struct {
	miniredis *miniredis.Miniredis
}

func NewRedis() *Redis {
	return &Redis{}
}

func (r *Redis) Start() error {
	var err error
	r.miniredis, err = miniredis.Run()
	if err != nil {
		return err
	}
	if err := os.Setenv("REDIS_ADDRESS", r.miniredis.Addr()); err != nil {
		return err
	}
	return nil
}
