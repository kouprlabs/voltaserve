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
	return &postgresUser{}
}

type postgresUser struct {
	ID                     string  `json:"id"`
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

func (postgresUser) TableName() string {
	return "user"
}

func (u postgresUser) GetID() string {
	return u.ID
}

func (u postgresUser) GetFullName() string {
	return u.FullName
}

func (u postgresUser) GetUsername() string {
	return u.Username
}

func (u postgresUser) GetEmail() string {
	return u.Email
}

func (u postgresUser) GetPicture() *string {
	return u.Picture
}

func (u postgresUser) GetIsEmailConfirmed() bool {
	return u.IsEmailConfirmed
}

func (u postgresUser) GetCreateTime() string {
	return u.CreateTime
}

func (u postgresUser) GetUpdateTime() *string {
	return u.UpdateTime
}

type userRepo struct {
	db *gorm.DB
}

func newUserRepo() *userRepo {
	return &userRepo{
		db: infra.GetDb(),
	}
}

func (repo *userRepo) Find(id string) (model.User, error) {
	var res = postgresUser{}
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
	var res = postgresUser{}
	db := repo.db.Where("email = ?", email).First(&res)
	if db.Error != nil {
		return nil, db.Error
	}
	return &res, nil
}

func (repo *userRepo) FindAll() ([]model.User, error) {
	var entities []*postgresUser
	db := repo.db.Raw(`select * from "user" u`).Scan(&entities)
	if db.Error != nil {
		return nil, db.Error
	}
	var res []model.User
	for _, u := range entities {
		res = append(res, u)
	}
	return res, nil
}
