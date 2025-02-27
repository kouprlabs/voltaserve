// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package infra_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/kouprlabs/voltaserve/shared/helper"
	"github.com/kouprlabs/voltaserve/shared/infra"
	"github.com/kouprlabs/voltaserve/shared/model"

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
	values := []repo.OrganizationNewModelOptions{
		{ID: "org_a", Name: "foo bar"},
		{ID: "org_b", Name: "hello world"},
	}
	for _, v := range values {
		err := search.NewOrganizationSearch().Index([]model.Organization{repo.NewOrganizationModelWithOptions(v)})
		s.Require().NoError(err)
	}

	hits, err := search.NewOrganizationSearch().Query("foo", infra.SearchQueryOptions{Limit: 10})
	s.Require().NoError(err)
	if s.Len(hits, 1) {
		s.Equal("org_a", hits[0].GetID())
	}

	hits, err = search.NewOrganizationSearch().Query("world", infra.SearchQueryOptions{Limit: 10})
	s.Require().NoError(err)
	if s.Len(hits, 1) {
		s.Equal("org_b", hits[0].GetID())
	}
}

func (s *BleveSuite) TestFilter() {
	values := []repo.FileNewModelOptions{
		{
			ID:          "file_a",
			WorkspaceID: "workspace_a",
			Name:        "lorem_ipsum.txt",
			Type:        model.FileTypeFile,
			Text:        helper.ToPtr("red apple"),
		},
		{
			ID:          "file_b",
			WorkspaceID: "workspace_b",
			Name:        "lorem_ipsum.txt",
			Type:        model.FileTypeFile,
			Text:        helper.ToPtr("pink strawberry"),
		},
		{
			ID:          "file_c",
			WorkspaceID: "workspace_c",
			Name:        "dolor_sit_amet.pdf",
			Type:        model.FileTypeFile,
			Text:        helper.ToPtr("yellow pineapple"),
		},
	}
	for _, v := range values {
		err := search.NewFileSearch().Index([]model.File{repo.NewFileModelWithOptions(v)})
		s.Require().NoError(err)
	}

	hits, err := search.NewFileSearch().Query("strawberry", infra.SearchQueryOptions{
		Limit:  10,
		Filter: fmt.Sprintf("workspaceId=\"workspace_b\" AND type=\"%s\"", model.FileTypeFile),
	})
	s.Require().NoError(err)
	if s.Len(hits, 1) {
		s.Equal("file_b", hits[0].GetID())
	}
}
