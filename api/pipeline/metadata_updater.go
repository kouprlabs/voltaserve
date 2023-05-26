package pipeline

import (
	"voltaserve/cache"
	"voltaserve/model"
	"voltaserve/repo"
	"voltaserve/search"
)

type metadataUpdater struct {
	snapshotRepo repo.SnapshotRepo
	fileRepo     repo.FileRepo
	fileCache    *cache.FileCache
	fileSearch   *search.FileSearch
}

func newMetadataUpdater() *metadataUpdater {
	return &metadataUpdater{
		snapshotRepo: repo.NewSnapshotRepo(),
		fileRepo:     repo.NewFileRepo(),
		fileCache:    cache.NewFileCache(),
		fileSearch:   search.NewFileSearch(),
	}
}

func (mu *metadataUpdater) update(snapshot model.Snapshot, fileId string) error {
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
	if err = mu.fileSearch.Update([]model.File{file}); err != nil {
		return err
	}
	return nil
}
