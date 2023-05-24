package repo

import "voltaserve/model"

type CoreSnapshotRepo interface {
	Find(id string) (model.SnapshotModel, error)
	Save(snapshot model.SnapshotModel) error
	MapWithFile(id string, fileId string) error
	DeleteMappingsForFile(fileId string) error
	FindAllForFile(fileId string) ([]*PostgresSnapshot, error)
	FindAllDangling() ([]model.SnapshotModel, error)
	DeleteAllDangling() error
	GetLatestVersionForFile(fileId string) (int64, error)
}
