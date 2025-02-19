// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package service_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/kouprlabs/voltaserve/api/test"
)

func TestMain(m *testing.M) {
	setup := test.NewSetup(m)
	if err := setup.Up(test.SetupOptions{
		Postgres: test.PostgresOptions{Port: 25432},
	}); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	code := m.Run()
	if err := setup.Down(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	os.Exit(code)
}
