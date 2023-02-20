package repo

import (
	"errors"
	"voltaserve/errorpkg"
	"voltaserve/infra"
	"voltaserve/model"

	"gorm.io/gorm"
)

type UserEntity struct {
	Id                     string  `json:"id"`
	FullName               string  `json:"fullName"`
	Username               string  `json:"username"`
	Email                  string  `json:"email"`
	Picture                *string `json:"picture"`
	IsEmailConfirmed       bool    `json:"isEmailConfirmed"`
	PasswordHash           string  `json:"passwordHash"`
	RefreshTokenValue      *string `json:"refreshTokenValue"`
	RefreshTokenValidTo    *int64  `json:"refreshTokenValidTo"`
	ResetPasswordToken     *string `json:"resetPasswordToken"`
	EmailConfirmationToken *string `json:"emailConfirmationToken"`
	CreateTime             string  `json:"createTime"`
	UpdateTime             *string `json:"updateTime"`
}

func (UserEntity) TableName() string {
	return "user"
}

func (u UserEntity) GetId() string {
	return u.Id
}

func (u UserEntity) GetFullName() string {
	return u.FullName
}

func (u UserEntity) GetUsername() string {
	return u.Username
}

func (u UserEntity) GetEmail() string {
	return u.Email
}

func (u UserEntity) GetPicture() *string {
	return u.Picture
}

func (u UserEntity) GetIsEmailConfirmed() bool {
	return u.IsEmailConfirmed
}

func (u UserEntity) GetCreateTime() string {
	return u.CreateTime
}

func (u UserEntity) GetUpdateTime() *string {
	return u.UpdateTime
}

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		db: infra.GetDb(),
	}
}

func (repo *UserRepo) Find(id string) (model.UserModel, error) {
	var res = UserEntity{}
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

func (repo *UserRepo) FindByEmail(email string) (model.UserModel, error) {
	var res = UserEntity{}
	db := repo.db.Where("email = ?", email).First(&res)
	if db.Error != nil {
		return nil, db.Error
	}
	return &res, nil
}

func (repo *UserRepo) FindAll() ([]model.UserModel, error) {
	var entities []*UserEntity
	db := repo.db.Raw(`select * from "user" u`).Scan(&entities)
	if db.Error != nil {
		return nil, db.Error
	}
	var res []model.UserModel
	for _, u := range entities {
		res = append(res, u)
	}
	return res, nil
}
