package repo

import "voltaserve/model"

type CoreUserRepo interface {
	Find(id string) (model.CoreUser, error)
	FindByEmail(email string) (model.CoreUser, error)
	FindAll() ([]model.CoreUser, error)
}
