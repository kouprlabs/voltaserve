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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/infra"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/repo"
	"github.com/kouprlabs/voltaserve/api/search"
)

func TestBleve(t *testing.T) {
	orgSearch := search.NewOrganizationSearch()
	opts := repo.NewOrganizationOptions{
		ID:   helper.NewID(),
		Name: "foo",
	}
	if err := orgSearch.Index([]model.Organization{repo.NewOrganizationWithOptions(opts)}); err != nil {
		t.Fatal(err)
	}
	hits, err := orgSearch.Query(opts.Name, infra.QueryOptions{Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	if assert.Len(t, hits, 1) {
		assert.Equal(t, hits[0].GetID(), opts.ID)
	}
}
