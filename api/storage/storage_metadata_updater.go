package storage

import (
	"voltaserve/cache"
	"voltaserve/model"
	"voltaserve/repo"
	"voltaserve/search"
)

type storageMetadataUpdater struct {
	snapshotRepo repo.SnapshotRepo
	fileRepo     repo.FileRepo
	fileCache    *cache.FileCache
	fileSearch   *search.FileSearch
}

func newMetadataUpdater() *storageMetadataUpdater {
	return &storageMetadataUpdater{
		snapshotRepo: repo.NewSnapshotRepo(),
		fileRepo:     repo.NewFileRepo(),
		fileCache:    cache.NewFileCache(),
		fileSearch:   search.NewFileSearch(),
	}
}

func (mu *storageMetadataUpdater) update(snapshot model.CoreSnapshot, fileId string) error {
	if err := repo.NewSnapshotRepo().Save(snapshot); err != nil {
		return err
	}
	file, err := mu.fileRepo.Find(fileId)
	if err != nil {
		return err
	}
	if err = mu.fileCache.Set(file); err != nil {
		return err
	}
	if err = mu.fileSearch.Update([]model.CoreFile{file}); err != nil {
		return err
	}
	return nil
}
