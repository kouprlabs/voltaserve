package repo

import "voltaserve/model"

type CoreUserRepo interface {
	Find(id string) (model.UserModel, error)
	FindByEmail(email string) (model.UserModel, error)
	FindAll() ([]model.UserModel, error)
}
