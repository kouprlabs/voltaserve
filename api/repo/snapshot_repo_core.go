package repo

import (
	"fmt"
	"voltaserve/config"
	"voltaserve/model"
)

type CoreSnapshotRepo interface {
	Find(id string) (model.CoreSnapshot, error)
	Save(snapshot model.CoreSnapshot) error
	MapWithFile(id string, fileId string) error
	DeleteMappingsForFile(fileId string) error
	FindAllDangling() ([]model.CoreSnapshot, error)
	DeleteAllDangling() error
	GetLatestVersionForFile(fileId string) (int64, error)
}

func NewSnapshotRepo() CoreSnapshotRepo {
	if config.GetConfig().DatabaseType == config.DATABASE_TYPE_POSTGRES {
		return NewPostgresSnapshotRepo()
	}
	panic(fmt.Sprintf("database type %s repo not implemented", config.GetConfig().DatabaseType))
}
