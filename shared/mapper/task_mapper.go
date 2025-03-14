// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package mapper

import (
	"errors"

	"github.com/kouprlabs/voltaserve/shared/cache"
	"github.com/kouprlabs/voltaserve/shared/config"
	"github.com/kouprlabs/voltaserve/shared/dto"
	"github.com/kouprlabs/voltaserve/shared/errorpkg"
	"github.com/kouprlabs/voltaserve/shared/model"
)

type TaskMapper struct {
	groupCache *cache.TaskCache
}

func NewTaskMapper(postgres config.PostgresConfig, redis config.RedisConfig, environment config.EnvironmentConfig) *TaskMapper {
	return &TaskMapper{
		groupCache: cache.NewTaskCache(postgres, redis, environment),
	}
}

func (mp *TaskMapper) Map(m model.Task) (*dto.Task, error) {
	return &dto.Task{
		ID:              m.GetID(),
		Name:            m.GetName(),
		Error:           m.GetError(),
		Percentage:      m.GetPercentage(),
		IsIndeterminate: m.GetIsIndeterminate(),
		UserID:          m.GetUserID(),
		Status:          m.GetStatus(),
		IsDismissible:   m.GetStatus() == model.TaskStatusSuccess || m.GetStatus() == model.TaskStatusError,
		Payload:         m.GetPayload(),
		CreateTime:      m.GetCreateTime(),
		UpdateTime:      m.GetUpdateTime(),
	}, nil
}

func (mp *TaskMapper) MapMany(tasks []model.Task) ([]*dto.Task, error) {
	res := make([]*dto.Task, 0)
	for _, task := range tasks {
		t, err := mp.Map(task)
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
