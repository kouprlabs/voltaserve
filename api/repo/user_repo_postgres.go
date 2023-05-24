package repo

import (
	"errors"
	"voltaserve/errorpkg"
	"voltaserve/infra"
	"voltaserve/model"

	"gorm.io/gorm"
)

type PostgresUser struct {
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

func (PostgresUser) TableName() string {
	return "user"
}

func (u PostgresUser) GetId() string {
	return u.Id
}

func (u PostgresUser) GetFullName() string {
	return u.FullName
}

func (u PostgresUser) GetUsername() string {
	return u.Username
}

func (u PostgresUser) GetEmail() string {
	return u.Email
}

func (u PostgresUser) GetPicture() *string {
	return u.Picture
}

func (u PostgresUser) GetIsEmailConfirmed() bool {
	return u.IsEmailConfirmed
}

func (u PostgresUser) GetCreateTime() string {
	return u.CreateTime
}

func (u PostgresUser) GetUpdateTime() *string {
	return u.UpdateTime
}

type PostgresUserRepo struct {
	db *gorm.DB
}

func NewPostgresUserRepo() *PostgresUserRepo {
	return &PostgresUserRepo{
		db: infra.GetDb(),
	}
}

func (repo *PostgresUserRepo) Find(id string) (model.UserModel, error) {
	var res = PostgresUser{}
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

func (repo *PostgresUserRepo) FindByEmail(email string) (model.UserModel, error) {
	var res = PostgresUser{}
	db := repo.db.Where("email = ?", email).First(&res)
	if db.Error != nil {
		return nil, db.Error
	}
	return &res, nil
}

func (repo *PostgresUserRepo) FindAll() ([]model.UserModel, error) {
	var entities []*PostgresUser
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
