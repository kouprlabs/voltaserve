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
	"github.com/kouprlabs/voltaserve/api/cache"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/repo"
)

type invitationMapper struct {
	orgCache   *cache.OrganizationCache
	userRepo   repo.UserRepo
	userMapper *userMapper
	orgMapper  *organizationMapper
}

func newInvitationMapper() *invitationMapper {
	return &invitationMapper{
		orgCache:   cache.NewOrganizationCache(),
		userRepo:   repo.NewUserRepo(),
		userMapper: newUserMapper(),
		orgMapper:  newOrganizationMapper(),
	}
}

func (mp *invitationMapper) mapOne(m model.Invitation, userID string) (*Invitation, error) {
	owner, err := mp.userRepo.Find(m.GetOwnerID())
	if err != nil {
		return nil, err
	}
	org, err := mp.orgCache.Get(m.GetOrganizationID())
	if err != nil {
		return nil, err
	}
	o, err := mp.orgMapper.mapOne(org, userID)
	if err != nil {
		return nil, err
	}
	return &Invitation{
		ID:           m.GetID(),
		Owner:        mp.userMapper.mapOne(owner),
		Email:        m.GetEmail(),
		Organization: o,
		Status:       m.GetStatus(),
		CreateTime:   m.GetCreateTime(),
		UpdateTime:   m.GetUpdateTime(),
	}, nil
}

func (mp *invitationMapper) mapMany(invitations []model.Invitation, userID string) ([]*Invitation, error) {
	res := make([]*Invitation, 0)
	for _, invitation := range invitations {
		i, err := mp.mapOne(invitation, userID)
		if err != nil {
			return nil, err
		}
		res = append(res, i)
	}
	return res, nil
}
