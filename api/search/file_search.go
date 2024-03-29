package search

import (
	"encoding/json"
	"voltaserve/infra"
	"voltaserve/model"
	"voltaserve/repo"
)

type FileSearch struct {
	search *infra.SearchManager
	index  string
	s3     *infra.S3Manager
}

func NewFileSearch() *FileSearch {
	return &FileSearch{
		index:  infra.FileSearchIndex,
		search: infra.NewSearchManager(),
		s3:     infra.NewS3Manager(),
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
		if f.GetSnapshots() != nil &&
			len(f.GetSnapshots()) > 0 &&
			f.GetSnapshots()[0].HasText() {
			var text string
			text, err := s.s3.GetText(f.GetSnapshots()[0].GetText().Key, f.GetSnapshots()[0].GetText().Bucket)
			if err != nil {
				return err
			}
			f.SetText(&text)
		}
	}
	return nil
}
