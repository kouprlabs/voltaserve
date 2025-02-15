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
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/infra"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/repo"
	"github.com/kouprlabs/voltaserve/api/search"
)

type BleveSuite struct {
	suite.Suite
}

func TestBleveSuite(t *testing.T) {
	suite.Run(t, new(BleveSuite))
}

func (s *BleveSuite) TestQuery() {
	orgSearch := search.NewOrganizationSearch()
	values := []repo.NewOrganizationOptions{
		{ID: "org_a", Name: "foo bar"},
		{ID: "org_b", Name: "hello world"},
	}
	for _, v := range values {
		err := orgSearch.Index([]model.Organization{repo.NewOrganizationWithOptions(v)})
		s.Require().NoError(err)
	}

	hits, err := orgSearch.Query("foo", infra.QueryOptions{Limit: 10})
	s.Require().NoError(err)
	if s.Len(hits, 1) {
		s.Equal("org_a", hits[0].GetID())
	}

	hits, err = orgSearch.Query("world", infra.QueryOptions{Limit: 10})
	s.Require().NoError(err)
	if s.Len(hits, 1) {
		s.Equal("org_b", hits[0].GetID())
	}
}

func (s *BleveSuite) TestFilter() {
	fileSearch := search.NewFileSearch()
	values := []repo.NewFileOptions{
		{
			ID:          "file_a",
			WorkspaceID: "workspace_a",
			Name:        "lorem_ipsum.txt",
			Type:        model.FileTypeFile,
			Text:        helper.ToPtr("exercitation ullamco laboris"),
		},
		{
			ID:          "file_b",
			WorkspaceID: "workspace_b",
			Name:        "lorem_ipsum.txt",
			Type:        model.FileTypeFile,
			Text:        helper.ToPtr("exercitation ullamco laboris"),
		},
		{
			ID:          "file_c",
			WorkspaceID: "workspace_c",
			Name:        "dolor_sit_amet.pdf",
			Type:        model.FileTypeFile,
			Text:        helper.ToPtr("sed et class dis libero"),
		},
	}
	for _, v := range values {
		err := fileSearch.Index([]model.File{repo.NewFileWithOptions(v)})
		s.Require().NoError(err)
	}

	hits, err := fileSearch.Query("exercitation", infra.QueryOptions{
		Limit:  10,
		Filter: fmt.Sprintf("workspaceId=\"workspace_b\" AND type=\"%s\"", model.FileTypeFile),
	})
	s.Require().NoError(err)
	if s.Len(hits, 1) {
		s.Equal("file_b", hits[0].GetID())
	}
}
