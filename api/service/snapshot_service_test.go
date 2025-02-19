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
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/kouprlabs/voltaserve/api/cache"
	"github.com/kouprlabs/voltaserve/api/client/conversion_client"
	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/repo"
	"github.com/kouprlabs/voltaserve/api/service"
	"github.com/kouprlabs/voltaserve/api/test"
)

type SnapshotServiceSuite struct {
	suite.Suite
	file      *service.File
	workspace *service.Workspace
	org       *service.Organization
	users     []model.User
}

func TestSnapshotServiceSuite(t *testing.T) {
	suite.Run(t, new(SnapshotServiceSuite))
}

func (s *SnapshotServiceSuite) SetupTest() {
	var err error
	s.users, err = test.CreateUsers(1)
	if err != nil {
		s.Fail(err.Error())
		return
	}
	s.org, err = test.CreateOrganization(s.users[0].GetID())
	if err != nil {
		s.Fail(err.Error())
		return
	}
	s.workspace, err = test.CreateWorkspace(s.org.ID, s.users[0].GetID())
	if err != nil {
		s.Fail(err.Error())
		return
	}
	s.file, err = test.CreateFile(s.workspace.ID, s.workspace.RootID, s.users[0].GetID())
	if err != nil {
		s.Fail(err.Error())
		return
	}
}

func (s *SnapshotServiceSuite) TestList() {
	snapshots := []model.Snapshot{
		repo.NewSnapshotModelWithOptions(repo.SnapshotNewModelOptions{
			ID:         helper.NewID(),
			Version:    1,
			CreateTime: helper.NewTimeString(),
		}),
		repo.NewSnapshotModelWithOptions(repo.SnapshotNewModelOptions{
			ID:         helper.NewID(),
			Version:    2,
			CreateTime: helper.TimeToString(time.Now().Add(-time.Hour)),
		}),
		repo.NewSnapshotModelWithOptions(repo.SnapshotNewModelOptions{
			ID:         helper.NewID(),
			Version:    3,
			CreateTime: helper.TimeToString(time.Now().Add(-2 * time.Hour)),
		}),
	}
	for _, snapshot := range snapshots {
		err := repo.NewSnapshotRepo().Insert(snapshot)
		s.Require().NoError(err)
		err = cache.NewSnapshotCache().Set(snapshot)
		s.Require().NoError(err)
		err = repo.NewSnapshotRepo().MapWithFile(snapshot.GetID(), s.file.ID)
		s.Require().NoError(err)
	}

	list, err := service.NewSnapshotService().List(s.file.ID, service.SnapshotListOptions{
		Page: 1,
		Size: 10,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(1), list.Page)
	s.Equal(uint64(3), list.Size)
	s.Equal(uint64(3), list.TotalElements)
	s.Equal(uint64(1), list.TotalPages)
	s.Equal(snapshots[0].GetID(), list.Data[0].ID)
	s.Equal(snapshots[1].GetID(), list.Data[1].ID)
	s.Equal(snapshots[2].GetID(), list.Data[2].ID)
}

func (s *SnapshotServiceSuite) TestList_Paginate() {
	snapshots := []model.Snapshot{
		repo.NewSnapshotModelWithOptions(repo.SnapshotNewModelOptions{
			ID:         helper.NewID(),
			Version:    1,
			CreateTime: helper.NewTimeString(),
		}),
		repo.NewSnapshotModelWithOptions(repo.SnapshotNewModelOptions{
			ID:         helper.NewID(),
			Version:    2,
			CreateTime: helper.TimeToString(time.Now().Add(-time.Hour)),
		}),
		repo.NewSnapshotModelWithOptions(repo.SnapshotNewModelOptions{
			ID:         helper.NewID(),
			Version:    3,
			CreateTime: helper.TimeToString(time.Now().Add(-2 * time.Hour)),
		}),
	}
	for _, snapshot := range snapshots {
		err := repo.NewSnapshotRepo().Insert(snapshot)
		s.Require().NoError(err)
		err = cache.NewSnapshotCache().Set(snapshot)
		s.Require().NoError(err)
		err = repo.NewSnapshotRepo().MapWithFile(snapshot.GetID(), s.file.ID)
		s.Require().NoError(err)
	}

	list, err := service.NewSnapshotService().List(s.file.ID, service.SnapshotListOptions{
		Page: 1,
		Size: 2,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(1), list.Page)
	s.Equal(uint64(2), list.Size)
	s.Equal(uint64(3), list.TotalElements)
	s.Equal(uint64(2), list.TotalPages)
	s.Equal(snapshots[0].GetID(), list.Data[0].ID)
	s.Equal(snapshots[1].GetID(), list.Data[1].ID)

	list, err = service.NewSnapshotService().List(s.file.ID, service.SnapshotListOptions{
		Page: 2,
		Size: 2,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(2), list.Page)
	s.Equal(uint64(1), list.Size)
	s.Equal(uint64(3), list.TotalElements)
	s.Equal(uint64(2), list.TotalPages)
	s.Equal(snapshots[2].GetID(), list.Data[0].ID)
}

func (s *SnapshotServiceSuite) TestList_SortByVersionDescending() {
	snapshots := []model.Snapshot{
		repo.NewSnapshotModelWithOptions(repo.SnapshotNewModelOptions{
			ID:         helper.NewID(),
			Version:    1,
			CreateTime: helper.NewTimeString(),
		}),
		repo.NewSnapshotModelWithOptions(repo.SnapshotNewModelOptions{
			ID:         helper.NewID(),
			Version:    2,
			CreateTime: helper.TimeToString(time.Now().Add(-time.Hour)),
		}),
		repo.NewSnapshotModelWithOptions(repo.SnapshotNewModelOptions{
			ID:         helper.NewID(),
			Version:    3,
			CreateTime: helper.TimeToString(time.Now().Add(-2 * time.Hour)),
		}),
	}
	for _, snapshot := range snapshots {
		err := repo.NewSnapshotRepo().Insert(snapshot)
		s.Require().NoError(err)
		err = cache.NewSnapshotCache().Set(snapshot)
		s.Require().NoError(err)
		err = repo.NewSnapshotRepo().MapWithFile(snapshot.GetID(), s.file.ID)
		s.Require().NoError(err)
	}

	list, err := service.NewSnapshotService().List(s.file.ID, service.SnapshotListOptions{
		Page:      1,
		Size:      3,
		SortBy:    service.SnapshotSortByVersion,
		SortOrder: service.SnapshotSortOrderDesc,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(snapshots[2].GetID(), list.Data[0].ID)
	s.Equal(snapshots[1].GetID(), list.Data[1].ID)
	s.Equal(snapshots[0].GetID(), list.Data[2].ID)
}

func (s *SnapshotServiceSuite) TestProbe() {
	snapshots := []model.Snapshot{
		repo.NewSnapshotModelWithOptions(repo.SnapshotNewModelOptions{
			ID:         helper.NewID(),
			Version:    1,
			CreateTime: helper.NewTimeString(),
		}),
		repo.NewSnapshotModelWithOptions(repo.SnapshotNewModelOptions{
			ID:         helper.NewID(),
			Version:    2,
			CreateTime: helper.TimeToString(time.Now().Add(-time.Hour)),
		}),
		repo.NewSnapshotModelWithOptions(repo.SnapshotNewModelOptions{
			ID:         helper.NewID(),
			Version:    3,
			CreateTime: helper.TimeToString(time.Now().Add(-2 * time.Hour)),
		}),
	}
	for _, snapshot := range snapshots {
		err := repo.NewSnapshotRepo().Insert(snapshot)
		s.Require().NoError(err)
		err = cache.NewSnapshotCache().Set(snapshot)
		s.Require().NoError(err)
		err = repo.NewSnapshotRepo().MapWithFile(snapshot.GetID(), s.file.ID)
		s.Require().NoError(err)
	}

	probe, err := service.NewSnapshotService().Probe(s.file.ID, service.SnapshotListOptions{
		Page: 1,
		Size: 10,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(3), probe.TotalElements)
	s.Equal(uint64(1), probe.TotalPages)
}

func (s *SnapshotServiceSuite) TestActivate() {
	snapshot := repo.NewSnapshotModelWithOptions(repo.SnapshotNewModelOptions{
		ID:         helper.NewID(),
		Version:    1,
		CreateTime: helper.NewTimeString(),
	})
	err := repo.NewSnapshotRepo().Insert(snapshot)
	s.Require().NoError(err)
	err = cache.NewSnapshotCache().Set(snapshot)
	s.Require().NoError(err)
	err = repo.NewSnapshotRepo().MapWithFile(snapshot.GetID(), s.file.ID)
	s.Require().NoError(err)

	file, err := service.NewSnapshotService().Activate(snapshot.GetID(), s.users[0].GetID())
	s.Require().NoError(err)
	s.Require().Equal(snapshot.GetID(), file.Snapshot.ID)
}

func (s *SnapshotServiceSuite) TestDetach() {
	snapshot := repo.NewSnapshotModelWithOptions(repo.SnapshotNewModelOptions{
		ID:         helper.NewID(),
		Version:    1,
		CreateTime: helper.NewTimeString(),
	})
	err := repo.NewSnapshotRepo().Insert(snapshot)
	s.Require().NoError(err)
	err = cache.NewSnapshotCache().Set(snapshot)
	s.Require().NoError(err)
	err = repo.NewSnapshotRepo().MapWithFile(snapshot.GetID(), s.file.ID)
	s.Require().NoError(err)

	err = service.NewSnapshotService().Detach(snapshot.GetID(), s.users[0].GetID())
	s.Require().NoError(err)

	file, err := cache.NewFileCache().Get(s.file.ID)
	s.Require().NoError(err)
	s.Require().Nil(file.GetSnapshotID())
}

func (s *SnapshotServiceSuite) TestPatch() {
	snapshot := repo.NewSnapshotModelWithOptions(repo.SnapshotNewModelOptions{
		ID:         helper.NewID(),
		Version:    1,
		CreateTime: helper.NewTimeString(),
	})
	err := repo.NewSnapshotRepo().Insert(snapshot)
	s.Require().NoError(err)
	err = cache.NewSnapshotCache().Set(snapshot)
	s.Require().NoError(err)

	patched, err := service.NewSnapshotService().Patch(snapshot.GetID(), service.SnapshotPatchOptions{
		Options: conversion_client.PipelineRunOptions{SnapshotID: snapshot.GetID()},
		Fields:  []string{repo.SnapshotFieldStatus},
		Status:  helper.ToPtr(model.SnapshotStatusProcessing),
	})
	s.Require().NoError(err)
	s.Require().Equal(model.SnapshotStatusProcessing, patched.Status)
}
