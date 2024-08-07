// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

package repo

import (
	"encoding/json"
	"errors"
	"slices"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/infra"
	"github.com/kouprlabs/voltaserve/api/log"
	"github.com/kouprlabs/voltaserve/api/model"
)

type SnapshotRepo interface {
	Find(id string) (model.Snapshot, error)
	FindByVersion(version int64) (model.Snapshot, error)
	FindAllForFile(fileID string) ([]model.Snapshot, error)
	FindAllDangling() ([]model.Snapshot, error)
	FindAllPrevious(fileID string, version int64) ([]model.Snapshot, error)
	GetIDsByFile(fileID string) ([]string, error)
	Insert(snapshot model.Snapshot) error
	Save(snapshot model.Snapshot) error
	Delete(id string) error
	Update(id string, opts SnapshotUpdateOptions) error
	MapWithFile(id string, fileID string) error
	BulkMapWithFile(entities []*SnapshotFileEntity, chunkSize int) error
	DeleteMappingsForFile(fileID string) error
	DeleteMappingsForTree(fileID string) error
	DeleteAllDangling() error
	GetLatestVersionForFile(fileID string) (int64, error)
	CountAssociations(id string) (int, error)
	Attach(sourceFileID string, targetFileID string) error
	Detach(id string, fileID string) error
}

func NewSnapshotRepo() SnapshotRepo {
	return newSnapshotRepo()
}

func NewSnapshot() model.Snapshot {
	return &snapshotEntity{}
}

type snapshotEntity struct {
	ID           string         `gorm:"column:id;size:36"   json:"id"`
	Version      int64          `gorm:"column:version"      json:"version"`
	Original     datatypes.JSON `gorm:"column:original"     json:"original,omitempty"`
	Preview      datatypes.JSON `gorm:"column:preview"      json:"preview,omitempty"`
	Text         datatypes.JSON `gorm:"column:text"         json:"text,omitempty"`
	OCR          datatypes.JSON `gorm:"column:ocr"          json:"ocr,omitempty"`
	Entities     datatypes.JSON `gorm:"column:entities"     json:"entities,omitempty"`
	Mosaic       datatypes.JSON `gorm:"column:mosaic"       json:"mosaic,omitempty"`
	Segmentation datatypes.JSON `gorm:"column:segmentation" json:"segmentation,omitempty"`
	Thumbnail    datatypes.JSON `gorm:"column:thumbnail"    json:"thumbnail,omitempty"`
	Status       string         `gorm:"column,status"       json:"status,omitempty"`
	Language     *string        `gorm:"column:language"     json:"language,omitempty"`
	TaskID       *string        `gorm:"column:task_id"      json:"taskId,omitempty"`
	CreateTime   string         `gorm:"column:create_time"  json:"createTime"`
	UpdateTime   *string        `gorm:"column:update_time"  json:"updateTime,omitempty"`
}

func (*snapshotEntity) TableName() string {
	return "snapshot"
}

func (s *snapshotEntity) BeforeCreate(*gorm.DB) (err error) {
	s.CreateTime = helper.NewTimestamp()
	return nil
}

func (s *snapshotEntity) BeforeSave(*gorm.DB) (err error) {
	timeNow := helper.NewTimestamp()
	s.UpdateTime = &timeNow
	return nil
}

func (s *snapshotEntity) GetID() string {
	return s.ID
}

func (s *snapshotEntity) GetVersion() int64 {
	return s.Version
}

func (s *snapshotEntity) GetOriginal() *model.S3Object {
	if s.Original.String() == "" {
		return nil
	}
	res := model.S3Object{}
	if err := json.Unmarshal([]byte(s.Original.String()), &res); err != nil {
		log.GetLogger().Fatal(err)
		return nil
	}
	return &res
}

func (s *snapshotEntity) GetPreview() *model.S3Object {
	if s.Preview.String() == "" {
		return nil
	}
	res := model.S3Object{}
	if err := json.Unmarshal([]byte(s.Preview.String()), &res); err != nil {
		log.GetLogger().Fatal(err)
		return nil
	}
	return &res
}

func (s *snapshotEntity) GetText() *model.S3Object {
	if s.Text.String() == "" {
		return nil
	}
	res := model.S3Object{}
	if err := json.Unmarshal([]byte(s.Text.String()), &res); err != nil {
		log.GetLogger().Fatal(err)
		return nil
	}
	return &res
}

func (s *snapshotEntity) GetOCR() *model.S3Object {
	if s.OCR.String() == "" {
		return nil
	}
	res := model.S3Object{}
	if err := json.Unmarshal([]byte(s.OCR.String()), &res); err != nil {
		log.GetLogger().Fatal(err)
		return nil
	}
	return &res
}

func (s *snapshotEntity) GetEntities() *model.S3Object {
	if s.Entities.String() == "" {
		return nil
	}
	res := model.S3Object{}
	if err := json.Unmarshal([]byte(s.Entities.String()), &res); err != nil {
		log.GetLogger().Fatal(err)
		return nil
	}
	return &res
}

func (s *snapshotEntity) GetMosaic() *model.S3Object {
	if s.Mosaic.String() == "" {
		return nil
	}
	res := model.S3Object{}
	if err := json.Unmarshal([]byte(s.Mosaic.String()), &res); err != nil {
		log.GetLogger().Fatal(err)
		return nil
	}
	return &res
}

func (s *snapshotEntity) GetSegmentation() *model.S3Object {
	if s.Segmentation.String() == "" {
		return nil
	}
	res := model.S3Object{}
	if err := json.Unmarshal([]byte(s.Segmentation.String()), &res); err != nil {
		log.GetLogger().Fatal(err)
		return nil
	}
	return &res
}

func (s *snapshotEntity) GetThumbnail() *model.S3Object {
	if s.Thumbnail.String() == "" {
		return nil
	}
	res := model.S3Object{}
	if err := json.Unmarshal([]byte(s.Thumbnail.String()), &res); err != nil {
		log.GetLogger().Fatal(err)
		return nil
	}
	return &res
}

func (s *snapshotEntity) GetStatus() string {
	return s.Status
}

func (s *snapshotEntity) GetLanguage() *string {
	return s.Language
}

func (s *snapshotEntity) GetTaskID() *string {
	return s.TaskID
}

func (s *snapshotEntity) SetID(id string) {
	s.ID = id
}

func (s *snapshotEntity) SetVersion(version int64) {
	s.Version = version
}

func (s *snapshotEntity) SetOriginal(m *model.S3Object) {
	if m == nil {
		s.Original = nil
	} else {
		b, err := json.Marshal(m)
		if err != nil {
			log.GetLogger().Fatal(err)
			return
		}
		if err := s.Original.UnmarshalJSON(b); err != nil {
			log.GetLogger().Fatal(err)
		}
	}
}

func (s *snapshotEntity) SetPreview(m *model.S3Object) {
	if m == nil {
		s.Preview = nil
	} else {
		b, err := json.Marshal(m)
		if err != nil {
			log.GetLogger().Fatal(err)
			return
		}
		if err := s.Preview.UnmarshalJSON(b); err != nil {
			log.GetLogger().Fatal(err)
		}
	}
}

func (s *snapshotEntity) SetText(m *model.S3Object) {
	if m == nil {
		s.Text = nil
	} else {
		b, err := json.Marshal(m)
		if err != nil {
			log.GetLogger().Fatal(err)
			return
		}
		if err := s.Text.UnmarshalJSON(b); err != nil {
			log.GetLogger().Fatal(err)
		}
	}
}

func (s *snapshotEntity) SetOCR(m *model.S3Object) {
	if m == nil {
		s.OCR = nil
	} else {
		b, err := json.Marshal(m)
		if err != nil {
			log.GetLogger().Fatal(err)
			return
		}
		if err := s.OCR.UnmarshalJSON(b); err != nil {
			log.GetLogger().Fatal(err)
		}
	}
}

func (s *snapshotEntity) SetEntities(m *model.S3Object) {
	if m == nil {
		s.Entities = nil
	} else {
		b, err := json.Marshal(m)
		if err != nil {
			log.GetLogger().Fatal(err)
			return
		}
		if err := s.Entities.UnmarshalJSON(b); err != nil {
			log.GetLogger().Fatal(err)
		}
	}
}

func (s *snapshotEntity) SetMosaic(m *model.S3Object) {
	if m == nil {
		s.Mosaic = nil
	} else {
		b, err := json.Marshal(m)
		if err != nil {
			log.GetLogger().Fatal(err)
			return
		}
		if err := s.Mosaic.UnmarshalJSON(b); err != nil {
			log.GetLogger().Fatal(err)
		}
	}
}

func (s *snapshotEntity) SetSegmentation(m *model.S3Object) {
	if m == nil {
		s.Segmentation = nil
	} else {
		b, err := json.Marshal(m)
		if err != nil {
			log.GetLogger().Fatal(err)
			return
		}
		if err := s.Segmentation.UnmarshalJSON(b); err != nil {
			log.GetLogger().Fatal(err)
		}
	}
}

func (s *snapshotEntity) SetThumbnail(m *model.S3Object) {
	if m == nil {
		s.Thumbnail = nil
	} else {
		b, err := json.Marshal(m)
		if err != nil {
			log.GetLogger().Fatal(err)
			return
		}
		if err := s.Thumbnail.UnmarshalJSON(b); err != nil {
			log.GetLogger().Fatal(err)
		}
	}
}

func (s *snapshotEntity) SetStatus(status string) {
	s.Status = status
}

func (s *snapshotEntity) SetLanguage(language string) {
	s.Language = &language
}

func (s *snapshotEntity) SetTaskID(taskID *string) {
	s.TaskID = taskID
}

func (s *snapshotEntity) HasOriginal() bool {
	return s.Original != nil
}

func (s *snapshotEntity) HasPreview() bool {
	return s.Preview != nil
}

func (s *snapshotEntity) HasText() bool {
	return s.Text != nil
}

func (s *snapshotEntity) HasOCR() bool {
	return s.OCR != nil
}

func (s *snapshotEntity) HasEntities() bool {
	return s.Entities != nil
}

func (s *snapshotEntity) HasMosaic() bool {
	return s.Mosaic != nil
}

func (s *snapshotEntity) HasSegmentation() bool {
	return s.Segmentation != nil
}

func (s *snapshotEntity) HasThumbnail() bool {
	return s.Thumbnail != nil
}

func (s *snapshotEntity) GetCreateTime() string {
	return s.CreateTime
}

func (s *snapshotEntity) GetUpdateTime() *string {
	return s.UpdateTime
}

type SnapshotFileEntity struct {
	SnapshotID string `gorm:"column:snapshot_id"`
	FileID     string `gorm:"column:file_id"`
	CreateTime string `gorm:"column:create_time"`
}

func (*SnapshotFileEntity) TableName() string {
	return "snapshot_file"
}

func (s *SnapshotFileEntity) BeforeCreate(*gorm.DB) (err error) {
	s.CreateTime = helper.NewTimestamp()
	return nil
}

type snapshotRepo struct {
	db *gorm.DB
}

func newSnapshotRepo() *snapshotRepo {
	return &snapshotRepo{
		db: infra.NewPostgresManager().GetDBOrPanic(),
	}
}

func (repo *snapshotRepo) find(id string) (*snapshotEntity, error) {
	var res snapshotEntity
	if db := repo.db.Where("id = ?", id).First(&res); db.Error != nil {
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			return nil, errorpkg.NewSnapshotNotFoundError(db.Error)
		} else {
			return nil, errorpkg.NewInternalServerError(db.Error)
		}
	}
	return &res, nil
}

func (repo *snapshotRepo) Find(id string) (model.Snapshot, error) {
	res, err := repo.find(id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (repo *snapshotRepo) FindByVersion(version int64) (model.Snapshot, error) {
	res := snapshotEntity{}
	db := repo.db.Where("version = ?", version).First(&res)
	if db.Error != nil {
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			return nil, errorpkg.NewSnapshotNotFoundError(db.Error)
		} else {
			return nil, errorpkg.NewInternalServerError(db.Error)
		}
	}
	return &res, nil
}

func (repo *snapshotRepo) Insert(snapshot model.Snapshot) error {
	if db := repo.db.Create(snapshot); db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *snapshotRepo) Save(snapshot model.Snapshot) error {
	if db := repo.db.Save(snapshot); db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *snapshotRepo) Delete(id string) error {
	snapshot, err := repo.find(id)
	if err != nil {
		return err
	}
	if db := repo.db.Delete(snapshot); db.Error != nil {
		return db.Error
	}
	return nil
}

type SnapshotUpdateOptions struct {
	Fields       []string `json:"fields"`
	Original     *model.S3Object
	Preview      *model.S3Object
	Text         *model.S3Object
	OCR          *model.S3Object
	Entities     *model.S3Object
	Mosaic       *model.S3Object
	Segmentation *model.S3Object
	Thumbnail    *model.S3Object
	Status       *string
	Language     *string
	TaskID       *string
}

const (
	SnapshotFieldOriginal     = "original"
	SnapshotFieldPreview      = "preview"
	SnapshotFieldText         = "text"
	SnapshotFieldOCR          = "ocr"
	SnapshotFieldEntities     = "entities"
	SnapshotFieldMosaic       = "mosaic"
	SnapshotFieldSegmentation = "segmentation"
	SnapshotFieldThumbnail    = "thumbnail"
	SnapshotFieldStatus       = "status"
	SnapshotFieldLanguage     = "language"
	SnapshotFieldTaskID       = "taskId"
)

func (repo *snapshotRepo) Update(id string, opts SnapshotUpdateOptions) error {
	snapshot, err := repo.find(id)
	if err != nil {
		return err
	}
	if slices.Contains(opts.Fields, SnapshotFieldOriginal) {
		snapshot.SetOriginal(opts.Original)
	}
	if slices.Contains(opts.Fields, SnapshotFieldPreview) {
		snapshot.SetPreview(opts.Preview)
	}
	if slices.Contains(opts.Fields, SnapshotFieldText) {
		snapshot.SetText(opts.Text)
	}
	if slices.Contains(opts.Fields, SnapshotFieldOCR) {
		snapshot.SetOCR(opts.OCR)
	}
	if slices.Contains(opts.Fields, SnapshotFieldEntities) {
		snapshot.SetEntities(opts.Entities)
	}
	if slices.Contains(opts.Fields, SnapshotFieldMosaic) {
		snapshot.SetMosaic(opts.Mosaic)
	}
	if slices.Contains(opts.Fields, SnapshotFieldSegmentation) {
		snapshot.SetSegmentation(opts.Segmentation)
	}
	if slices.Contains(opts.Fields, SnapshotFieldThumbnail) {
		snapshot.SetThumbnail(opts.Thumbnail)
	}
	if slices.Contains(opts.Fields, SnapshotFieldStatus) {
		snapshot.SetStatus(*opts.Status)
	}
	if slices.Contains(opts.Fields, SnapshotFieldLanguage) {
		snapshot.SetLanguage(*opts.Language)
	}
	if slices.Contains(opts.Fields, SnapshotFieldTaskID) {
		snapshot.SetTaskID(opts.TaskID)
	}
	if db := repo.db.Save(&snapshot); db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *snapshotRepo) MapWithFile(id string, fileID string) error {
	if db := repo.db.Exec("INSERT INTO snapshot_file (snapshot_id, file_id, create_time) VALUES (?, ?, ?)", id, fileID, helper.NewTimestamp()); db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *snapshotRepo) BulkMapWithFile(entities []*SnapshotFileEntity, chunkSize int) error {
	if db := repo.db.CreateInBatches(entities, chunkSize); db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *snapshotRepo) DeleteMappingsForFile(fileID string) error {
	if db := repo.db.Exec("DELETE FROM snapshot_file WHERE file_id = ?", fileID); db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *snapshotRepo) findAllForFile(fileID string) ([]*snapshotEntity, error) {
	var res []*snapshotEntity
	db := repo.db.
		Raw(`SELECT * FROM snapshot s
             LEFT JOIN snapshot_file sf ON s.id = sf.snapshot_id
             WHERE sf.file_id = ? ORDER BY s.version`,
			fileID).
		Scan(&res)
	if db.Error != nil {
		return nil, db.Error
	}
	return res, nil
}

func (repo *snapshotRepo) DeleteMappingsForTree(fileID string) error {
	db := repo.db.
		Exec(`WITH RECURSIVE rec (id, parent_id, create_time) AS
              (SELECT f.id, f.parent_id, f.create_time FROM file f WHERE f.parent_id = ?
              UNION SELECT f.id, f.parent_id, f.create_time FROM rec, file f WHERE f.parent_id = rec.id)
              DELETE FROM snapshot_file WHERE file_id in (SELECT id FROM rec);`,
			fileID)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *snapshotRepo) FindAllForFile(fileID string) ([]model.Snapshot, error) {
	snapshots, err := repo.findAllForFile(fileID)
	if err != nil {
		return nil, err
	}
	var res []model.Snapshot
	for _, s := range snapshots {
		res = append(res, s)
	}
	return res, nil
}

func (repo *snapshotRepo) FindAllDangling() ([]model.Snapshot, error) {
	var entities []*snapshotEntity
	db := repo.db.
		Raw(`SELECT * FROM snapshot s
             LEFT JOIN snapshot_file sf ON s.id = sf.snapshot_id
             WHERE sf.snapshot_id IS NULL`).
		Scan(&entities)
	if db.Error != nil {
		return nil, db.Error
	}
	var res []model.Snapshot
	for _, s := range entities {
		res = append(res, s)
	}
	return res, nil
}

func (repo *snapshotRepo) FindAllPrevious(fileID string, version int64) ([]model.Snapshot, error) {
	var entities []*snapshotEntity
	db := repo.db.
		Raw(`SELECT * FROM snapshot s
             LEFT JOIN snapshot_file sf ON s.id = sf.snapshot_id
             WHERE sf.file_id = ? AND s.version < ?
             ORDER BY s.version DESC`,
			fileID, version).
		Scan(&entities)
	if db.Error != nil {
		return nil, db.Error
	}
	var res []model.Snapshot
	for _, s := range entities {
		res = append(res, s)
	}
	return res, nil
}

func (repo *snapshotRepo) GetIDsByFile(fileID string) ([]string, error) {
	type Value struct {
		Result string
	}
	var values []Value
	db := repo.db.
		Raw("SELECT snapshot_id result FROM snapshot_file WHERE file_id = ?", fileID).
		Scan(&values)
	if db.Error != nil {
		return nil, db.Error
	}
	res := []string{}
	for _, v := range values {
		res = append(res, v.Result)
	}
	return res, nil
}

func (repo *snapshotRepo) DeleteAllDangling() error {
	if db := repo.db.
		Exec(`DELETE FROM snapshot
              WHERE id IN (SELECT s.id FROM (SELECT * FROM snapshot) s 
              LEFT JOIN snapshot_file sf ON s.id = sf.snapshot_id WHERE sf.snapshot_id IS NULL)`); db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *snapshotRepo) GetLatestVersionForFile(fileID string) (int64, error) {
	type Result struct {
		Result int64
	}
	var res Result
	if db := repo.db.
		Raw(`SELECT coalesce(max(s.version), 0) result 
             FROM snapshot s LEFT JOIN snapshot_file map ON s.id = map.snapshot_id
             WHERE map.file_id = ?`,
			fileID).
		Scan(&res); db.Error != nil {
		return 0, db.Error
	}
	return res.Result, nil
}

func (repo *snapshotRepo) CountAssociations(id string) (int, error) {
	type Result struct {
		Count int
	}
	var res Result
	if db := repo.db.
		Raw("SELECT COUNT(*) count FROM snapshot_file WHERE snapshot_id = ?", id).
		Scan(&res); db.Error != nil {
		return 0, db.Error
	}
	return res.Count, nil
}

func (repo *snapshotRepo) Attach(sourceFileID string, targetFileID string) error {
	if db := repo.db.
		Exec(`INSERT INTO snapshot_file (snapshot_id, file_id, create_time) SELECT s.id, ?, ?
              FROM snapshot s LEFT JOIN snapshot_file map ON s.id = map.snapshot_id
              WHERE map.file_id = ? ORDER BY s.version DESC LIMIT 1`,
			targetFileID, helper.NewTimestamp(), sourceFileID); db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *snapshotRepo) Detach(id string, fileID string) error {
	if db := repo.db.Exec("DELETE FROM snapshot_file WHERE snapshot_id = ? AND file_id = ?", id, fileID); db.Error != nil {
		return db.Error
	}
	return nil
}
