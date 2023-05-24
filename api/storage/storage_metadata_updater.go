package storage

import (
	"voltaserve/cache"
	"voltaserve/model"
	"voltaserve/repo"
	"voltaserve/search"
)

type storageMetadataUpdater struct {
	snapshotRepo repo.CoreSnapshotRepo
	fileRepo     repo.CoreFileRepo
	fileCache    *cache.FileCache
	fileSearch   *search.FileSearch
}

func newMetadataUpdater() *storageMetadataUpdater {
	return &storageMetadataUpdater{
		snapshotRepo: repo.NewPostgresSnapshotRepo(),
		fileRepo:     repo.NewPostgresFileRepo(),
		fileCache:    cache.NewFileCache(),
		fileSearch:   search.NewFileSearch(),
	}
}

func (mu *storageMetadataUpdater) update(snapshot model.SnapshotModel, fileId string) error {
	if err := repo.NewPostgresSnapshotRepo().Save(snapshot); err != nil {
		return err
	}
	file, err := mu.fileRepo.Find(fileId)
	if err != nil {
		return err
	}
	if err = mu.fileCache.Set(file); err != nil {
		return err
	}
	if err = mu.fileSearch.Update([]model.FileModel{file}); err != nil {
		return err
	}
	return nil
}
