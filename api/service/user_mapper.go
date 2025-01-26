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
	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/model"
)

type userMapper struct{}

func newUserMapper() *userMapper {
	return &userMapper{}
}

func (mp *userMapper) mapOne(user model.User) *User {
	res := &User{
		ID:         user.GetID(),
		FullName:   user.GetFullName(),
		Email:      user.GetEmail(),
		Username:   user.GetUsername(),
		CreateTime: user.GetCreateTime(),
		UpdateTime: user.GetUpdateTime(),
	}
	if user.GetPicture() != nil {
		res.Picture = &Picture{
			Extension: helper.Base64ToExtension(*user.GetPicture()),
		}
	}
	return res
}

func (mp *userMapper) mapMany(users []model.User) ([]*User, error) {
	res := make([]*User, 0)
	for _, user := range users {
		res = append(res, mp.mapOne(user))
	}
	return res, nil
}
