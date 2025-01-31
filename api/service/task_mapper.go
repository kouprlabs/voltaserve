// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package service

import (
	"errors"

	"github.com/kouprlabs/voltaserve/api/cache"
	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/model"
)

type taskMapper struct {
	groupCache cache.TaskCache
}

func newTaskMapper() *taskMapper {
	return &taskMapper{
		groupCache: cache.NewTaskCache(),
	}
}

func (mp *taskMapper) mapOne(m model.Task) (*Task, error) {
	return &Task{
		ID:              m.GetID(),
		Name:            m.GetName(),
		Error:           m.GetError(),
		Percentage:      m.GetPercentage(),
		IsIndeterminate: m.GetIsIndeterminate(),
		UserID:          m.GetUserID(),
		Status:          m.GetStatus(),
		Payload:         m.GetPayload(),
		CreateTime:      m.GetCreateTime(),
		UpdateTime:      m.GetUpdateTime(),
	}, nil
}

func (mp *taskMapper) mapMany(tasks []model.Task) ([]*Task, error) {
	res := make([]*Task, 0)
	for _, task := range tasks {
		t, err := mp.mapOne(task)
		if err != nil {
			var e *errorpkg.ErrorResponse
			if errors.As(err, &e) && e.Code == errorpkg.NewTaskNotFoundError(nil).Code {
				continue
			} else {
				return nil, err
			}
		}
		res = append(res, t)
	}
	return res, nil
}
