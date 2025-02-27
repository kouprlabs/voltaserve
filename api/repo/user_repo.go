// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package repo

import (
	"errors"

	"gorm.io/gorm"

	"github.com/kouprlabs/voltaserve/shared/errorpkg"
	"github.com/kouprlabs/voltaserve/shared/model"

	"github.com/kouprlabs/voltaserve/api/infra"
)

type userEntity struct {
	ID                     string  `gorm:"column:id"                       json:"id"`
	FullName               string  `gorm:"column:full_name"                json:"fullName"`
	Username               string  `gorm:"column:username"                 json:"username"`
	Email                  string  `gorm:"column:email"                    json:"email"`
	Picture                *string `gorm:"column:picture"                  json:"picture"`
	IsEmailConfirmed       bool    `gorm:"column:is_email_confirmed"       json:"isEmailConfirmed"`
	PasswordHash           string  `gorm:"column:password_hash"            json:"passwordHash"`
	RefreshTokenValue      *string `gorm:"column:refresh_token_value"      json:"refreshTokenValue"`
	RefreshTokenValidTo    *int64  `gorm:"column:refresh_token_valid_to"   json:"refreshTokenValidTo"`
	ResetPasswordToken     *string `gorm:"column:reset_password_token"     json:"resetPasswordToken"`
	EmailConfirmationToken *string `gorm:"column:email_confirmation_token" json:"emailConfirmationToken"`
	CreateTime             string  `gorm:"column:create_time"              json:"createTime"`
	UpdateTime             *string `gorm:"column:update_time"              json:"updateTime"`
}

func (userEntity) TableName() string {
	return "user"
}

func (u userEntity) GetID() string {
	return u.ID
}

func (u userEntity) GetFullName() string {
	return u.FullName
}

func (u userEntity) GetUsername() string {
	return u.Username
}

func (u userEntity) GetEmail() string {
	return u.Email
}

func (u userEntity) GetPicture() *string {
	return u.Picture
}

func (u userEntity) GetIsEmailConfirmed() bool {
	return u.IsEmailConfirmed
}

func (u userEntity) GetCreateTime() string {
	return u.CreateTime
}

func (u userEntity) GetUpdateTime() *string {
	return u.UpdateTime
}

func NewUserModel() model.User {
	return &userEntity{}
}

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		db: infra.NewPostgresManager().GetDBOrPanic(),
	}
}

func (repo *UserRepo) Find(id string) (model.User, error) {
	res := userEntity{}
	db := repo.db.Where("id = ?", id).First(&res)
	if db.Error != nil {
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			return nil, errorpkg.NewUserNotFoundError(db.Error)
		} else {
			return nil, errorpkg.NewInternalServerError(db.Error)
		}
	}
	return &res, nil
}

func (repo *UserRepo) FindOrNil(id string) model.User {
	res, err := repo.Find(id)
	if err != nil {
		return nil
	}
	return res
}

func (repo *UserRepo) FindByEmail(email string) (model.User, error) {
	res := userEntity{}
	db := repo.db.Where("LOWER(email) = LOWER(?)", email).First(&res)
	if db.Error != nil {
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			return nil, errorpkg.NewUserNotFoundError(db.Error)
		} else {
			return nil, errorpkg.NewInternalServerError(db.Error)
		}
	}
	return &res, nil
}

func (repo *UserRepo) FindAll() ([]model.User, error) {
	var entities []*userEntity
	db := repo.db.Raw(`SELECT * FROM "user"`).Scan(&entities)
	if db.Error != nil {
		return nil, db.Error
	}
	var res []model.User
	for _, u := range entities {
		res = append(res, u)
	}
	return res, nil
}

func (repo *UserRepo) Count() (int64, error) {
	var count int64
	db := repo.db.Model(&userEntity{}).Count(&count)
	if db.Error != nil {
		return -1, db.Error
	}
	return count, nil
}
