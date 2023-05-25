package repo

import "voltaserve/model"

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
	return NewPostgresSnapshotRepo()
}
