// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

package repo

import (
	"errors"
	"voltaserve/errorpkg"
	"voltaserve/infra"
	"voltaserve/model"

	"gorm.io/gorm"
)

type UserRepo interface {
	Find(id string) (model.User, error)
	FindByEmail(email string) (model.User, error)
	FindAll() ([]model.User, error)
}

func NewUserRepo() UserRepo {
	return newUserRepo()
}

func NewUser() model.User {
	return &userEntity{}
}

type userEntity struct {
	ID                     string  `json:"id" gorm:"column:id"`
	FullName               string  `json:"fullName" gorm:"column:full_name"`
	Username               string  `json:"username" gorm:"column:username"`
	Email                  string  `json:"email" gorm:"column:email"`
	Picture                *string `json:"picture" gorm:"column:picture"`
	IsEmailConfirmed       bool    `json:"isEmailConfirmed" gorm:"column:is_email_confirmed"`
	PasswordHash           string  `json:"passwordHash" gorm:"column:password_hash"`
	RefreshTokenValue      *string `json:"refreshTokenValue" gorm:"column:refresh_token_value"`
	RefreshTokenValidTo    *int64  `json:"refreshTokenValidTo" gorm:"column:refresh_token_valid_to"`
	ResetPasswordToken     *string `json:"resetPasswordToken" gorm:"column:reset_password_token"`
	EmailConfirmationToken *string `json:"emailConfirmationToken" gorm:"column:email_confirmation_token"`
	CreateTime             string  `json:"createTime" gorm:"column:create_time"`
	UpdateTime             *string `json:"updateTime" gorm:"column:update_time"`
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

type userRepo struct {
	db *gorm.DB
}

func newUserRepo() *userRepo {
	return &userRepo{
		db: infra.NewPostgresManager().GetDBOrPanic(),
	}
}

func (repo *userRepo) Find(id string) (model.User, error) {
	var res = userEntity{}
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

func (repo *userRepo) FindByEmail(email string) (model.User, error) {
	var res = userEntity{}
	db := repo.db.Where("email = ?", email).First(&res)
	if db.Error != nil {
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			return nil, errorpkg.NewUserNotFoundError(db.Error)
		} else {
			return nil, errorpkg.NewInternalServerError(db.Error)
		}
	}
	return &res, nil
}

func (repo *userRepo) FindAll() ([]model.User, error) {
	var entities []*userEntity
	db := repo.db.Raw(`select * from "user"`).Scan(&entities)
	if db.Error != nil {
		return nil, db.Error
	}
	var res []model.User
	for _, u := range entities {
		res = append(res, u)
	}
	return res, nil
}
