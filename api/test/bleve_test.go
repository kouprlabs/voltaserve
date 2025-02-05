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

func TestBleve_Query(t *testing.T) {
	orgSearch := search.NewOrganizationSearch()
	values := []repo.NewOrganizationOptions{
		{ID: helper.NewID(), Name: "foo bar"},
		{ID: helper.NewID(), Name: "hello world"},
	}
	for _, v := range values {
		if err := orgSearch.Index([]model.Organization{repo.NewOrganizationWithOptions(v)}); err != nil {
			t.Fatal(err)
		}
	}
	hits, err := orgSearch.Query("foo", infra.QueryOptions{Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	if assert.Len(t, hits, 1) {
		assert.Equal(t, "foo bar", hits[0].GetName())
	}
	hitsAgain, err := orgSearch.Query("world", infra.QueryOptions{Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	if assert.Len(t, hitsAgain, 1) {
		assert.Equal(t, "hello world", hitsAgain[0].GetName())
	}
}

func TestBleve_Filter(t *testing.T) {
	orgSearch := search.NewOrganizationSearch()
	values := []repo.NewOrganizationOptions{
		{ID: "id_a", Name: "lorem"},
		{ID: "id_b", Name: "ipsum"},
		{ID: "id_c", Name: "lorem"},
	}
	for _, v := range values {
		if err := orgSearch.Index([]model.Organization{repo.NewOrganizationWithOptions(v)}); err != nil {
			t.Fatal(err)
		}
	}
	hits, err := orgSearch.Query("lorem", infra.QueryOptions{Limit: 10, Filter: "id=\"id_c\""})
	if err != nil {
		t.Fatal(err)
	}
	if assert.Len(t, hits, 1) {
		assert.Equal(t, "id_c", hits[0].GetID())
	}
}
