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

	"gorm.io/gorm"

	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/infra"
	"github.com/kouprlabs/voltaserve/api/model"
)

type UserRepo interface {
	Find(id string) (model.User, error)
	FindByEmail(email string) (model.User, error)
	Count() (int64, error)
}

func NewUserRepo() UserRepo {
	return newUserRepo()
}

func NewUser() model.User {
	return &userEntity{}
}

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

type userRepo struct {
	db *gorm.DB
}

func newUserRepo() *userRepo {
	return &userRepo{
		db: infra.NewPostgresManager().GetDBOrPanic(),
	}
}

func (repo *userRepo) Find(id string) (model.User, error) {
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

func (repo *userRepo) FindByEmail(email string) (model.User, error) {
	res := userEntity{}
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

func (repo *userRepo) Count() (int64, error) {
	type Result struct {
		Result int64
	}
	var res Result
	db := repo.db.
		Raw(`SELECT count(*) as result FROM "user"`).
		Scan(&res)
	if db.Error != nil {
		return 0, db.Error
	}
	return res.Result, nil
}
