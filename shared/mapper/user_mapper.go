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
	"github.com/kouprlabs/voltaserve/shared/dto"
	"github.com/kouprlabs/voltaserve/shared/helper"
	"github.com/kouprlabs/voltaserve/shared/model"
)

type UserMapper struct{}

func NewUserMapper() *UserMapper {
	return &UserMapper{}
}

func (mp *UserMapper) Map(user model.User) *dto.User {
	res := &dto.User{
		ID:         user.GetID(),
		FullName:   user.GetFullName(),
		Email:      user.GetEmail(),
		Username:   user.GetUsername(),
		CreateTime: user.GetCreateTime(),
		UpdateTime: user.GetUpdateTime(),
	}
	if user.GetPicture() != nil {
		res.Picture = &dto.Picture{
			Extension: helper.Base64ToExtension(*user.GetPicture()),
		}
	}
	return res
}

func (mp *UserMapper) MapMany(users []model.User) ([]*dto.User, error) {
	res := make([]*dto.User, 0)
	for _, user := range users {
		res = append(res, mp.Map(user))
	}
	return res, nil
}
