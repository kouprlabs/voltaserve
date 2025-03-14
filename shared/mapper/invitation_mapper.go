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
	"github.com/kouprlabs/voltaserve/shared/cache"
	"github.com/kouprlabs/voltaserve/shared/config"
	"github.com/kouprlabs/voltaserve/shared/dto"
	"github.com/kouprlabs/voltaserve/shared/model"
	"github.com/kouprlabs/voltaserve/shared/repo"
)

type InvitationMapper struct {
	orgCache   *cache.OrganizationCache
	userRepo   *repo.UserRepo
	userMapper *UserMapper
	orgMapper  *OrganizationMapper
}

func NewInvitationMapper(postgres config.PostgresConfig, redis config.RedisConfig, environment config.EnvironmentConfig) *InvitationMapper {
	return &InvitationMapper{
		orgCache:   cache.NewOrganizationCache(postgres, redis, environment),
		userRepo:   repo.NewUserRepo(postgres, environment),
		userMapper: NewUserMapper(),
		orgMapper:  NewOrganizationMapper(postgres, redis, environment),
	}
}

func (mp *InvitationMapper) Map(m model.Invitation, userID string) (*dto.Invitation, error) {
	owner, err := mp.userRepo.Find(m.GetOwnerID())
	if err != nil {
		return nil, err
	}
	org, err := mp.orgCache.Get(m.GetOrganizationID())
	if err != nil {
		return nil, err
	}
	o, err := mp.orgMapper.Map(org, userID)
	if err != nil {
		return nil, err
	}
	return &dto.Invitation{
		ID:           m.GetID(),
		Owner:        mp.userMapper.Map(owner),
		Email:        m.GetEmail(),
		Organization: o,
		Status:       m.GetStatus(),
		CreateTime:   m.GetCreateTime(),
		UpdateTime:   m.GetUpdateTime(),
	}, nil
}

func (mp *InvitationMapper) MapMany(invitations []model.Invitation, userID string) ([]*dto.Invitation, error) {
	res := make([]*dto.Invitation, 0)
	for _, invitation := range invitations {
		i, err := mp.Map(invitation, userID)
		if err != nil {
			return nil, err
		}
		res = append(res, i)
	}
	return res, nil
}
