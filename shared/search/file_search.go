// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package search

import (
	"encoding/json"

	"github.com/minio/minio-go/v7"

	"github.com/kouprlabs/voltaserve/shared/config"
	"github.com/kouprlabs/voltaserve/shared/infra"
	"github.com/kouprlabs/voltaserve/shared/model"
	"github.com/kouprlabs/voltaserve/shared/repo"
)

type FileSearch struct {
	search       infra.SearchManager
	index        string
	s3           infra.S3Manager
	snapshotRepo *repo.SnapshotRepo
}

type fileEntity struct {
	ID          string  `json:"id"`
	WorkspaceID string  `json:"workspaceId"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	ParentID    *string `json:"parentId,omitempty"`
	Text        *string `json:"text,omitempty"`
	Summary     *string `json:"summary,omitempty"`
	SnapshotID  *string `json:"snapshotId,omitempty"`
	CreateTime  string  `json:"createTime"`
	UpdateTime  *string `json:"updateTime,omitempty"`
}

func (f fileEntity) GetID() string {
	return f.ID
}

func NewFileSearch(postgres config.PostgresConfig, search config.SearchConfig, s3 config.S3Config, environment config.EnvironmentConfig) *FileSearch {
	return &FileSearch{
		index:        infra.FileSearchIndex,
		search:       infra.NewSearchManager(search, environment),
		s3:           infra.NewS3Manager(s3, environment),
		snapshotRepo: repo.NewSnapshotRepo(postgres, environment),
	}
}

func (s *FileSearch) Index(files []model.File) (err error) {
	if len(files) == 0 {
		return nil
	}
	if err = s.populateTextAndSummaryFields(files); err != nil {
		return err
	}
	var models []infra.SearchModel
	for _, f := range files {
		models = append(models, s.mapEntity(f))
	}
	if err := s.search.Index(s.index, models); err != nil {
		return err
	}
	return nil
}

func (s *FileSearch) Update(files []model.File) (err error) {
	if len(files) == 0 {
		return nil
	}
	if err = s.populateTextAndSummaryFields(files); err != nil {
		return err
	}
	var models []infra.SearchModel
	for _, f := range files {
		models = append(models, s.mapEntity(f))
	}
	if err := s.search.Update(s.index, models); err != nil {
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

func (s *FileSearch) Query(query string, opts infra.SearchQueryOptions) ([]model.File, error) {
	hits, err := s.search.Query(s.index, query, opts)
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
		file := repo.NewFileModel()
		if err = json.Unmarshal(b, &file); err != nil {
			return nil, err
		}
		res = append(res, file)
	}
	return res, nil
}

func (s *FileSearch) populateTextAndSummaryFields(files []model.File) error {
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
			f.SetSummary(snapshot.GetSummary())
		}
	}
	return nil
}

func (s *FileSearch) mapEntity(file model.File) *fileEntity {
	return &fileEntity{
		ID:          file.GetID(),
		WorkspaceID: file.GetWorkspaceID(),
		Name:        file.GetName(),
		Type:        file.GetType(),
		ParentID:    file.GetParentID(),
		Text:        file.GetText(),
		Summary:     file.GetSummary(),
		SnapshotID:  file.GetSnapshotID(),
		CreateTime:  file.GetCreateTime(),
		UpdateTime:  file.GetUpdateTime(),
	}
}
