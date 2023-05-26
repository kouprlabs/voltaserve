package repo

import (
	"fmt"
	"voltaserve/config"
	"voltaserve/model"
)

type CoreUserRepo interface {
	Find(id string) (model.CoreUser, error)
	FindByEmail(email string) (model.CoreUser, error)
	FindAll() ([]model.CoreUser, error)
}

func NewUserRepo() CoreUserRepo {
	if config.GetConfig().DatabaseType == config.DATABASE_TYPE_POSTGRES {
		return NewPostgresUserRepo()
	}
	panic(fmt.Sprintf("database type %s repo not implemented", config.GetConfig().DatabaseType))
}
