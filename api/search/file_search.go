package search

import (
	"encoding/json"
	"voltaserve/infra"
	"voltaserve/model"
	"voltaserve/repo"

	"github.com/minio/minio-go/v7"
)

type FileSearch struct {
	search       *infra.SearchManager
	index        string
	s3           *infra.S3Manager
	snapshotRepo repo.SnapshotRepo
}

func NewFileSearch() *FileSearch {
	return &FileSearch{
		index:        infra.FileSearchIndex,
		search:       infra.NewSearchManager(),
		s3:           infra.NewS3Manager(),
		snapshotRepo: repo.NewSnapshotRepo(),
	}
}

func (s *FileSearch) Index(files []model.File) (err error) {
	if len(files) == 0 {
		return nil
	}
	if err = s.populateTextField(files); err != nil {
		return err
	}
	var res []infra.SearchModel
	for _, f := range files {
		res = append(res, f)
	}
	if err := s.search.Index(s.index, res); err != nil {
		return err
	}
	return nil
}

func (s *FileSearch) Update(files []model.File) (err error) {
	if len(files) == 0 {
		return nil
	}
	if err = s.populateTextField(files); err != nil {
		return err
	}
	var res []infra.SearchModel
	for _, f := range files {
		res = append(res, f)
	}
	if err := s.search.Update(s.index, res); err != nil {
		return err
	}
	return nil
}

func (s *FileSearch) Delete(ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	if err := s.search.Delete(s.index, ids); err != nil {
		return err
	}
	return nil
}

func (s *FileSearch) Query(query string) ([]model.File, error) {
	hits, err := s.search.Query(s.index, query)
	if err != nil {
		return nil, err
	}
	var res []model.File
	for _, v := range hits {
		var b []byte
		b, err = json.Marshal(v)
		if err != nil {
			return nil, err
		}
		file := repo.NewFile()
		if err = json.Unmarshal(b, &file); err != nil {
			return nil, err
		}
		res = append(res, file)
	}
	return res, nil
}

func (s *FileSearch) populateTextField(files []model.File) error {
	for _, f := range files {
		if f.GetType() == model.FileTypeFile && f.GetSnapshotID() != nil {
			snapshot, err := s.snapshotRepo.Find(*f.GetSnapshotID())
			if err != nil {
				return err
			}
			if snapshot.HasText() {
				text, err := s.s3.GetText(snapshot.GetText().Key, snapshot.GetText().Bucket, minio.GetObjectOptions{})
				if err != nil {
					return err
				}
				f.SetText(&text)
			}
		}
	}
	return nil
}
