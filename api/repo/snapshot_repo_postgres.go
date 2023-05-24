package repo

import (
	"encoding/json"
	"errors"
	"log"
	"time"
	"voltaserve/errorpkg"
	"voltaserve/infra"
	"voltaserve/model"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type PostgresSnapshot struct {
	Id         string         `json:"id" gorm:"column:id;size:36"`
	Version    int64          `json:"version" gorm:"column:version"`
	Original   datatypes.JSON `json:"original,omitempty" gorm:"column:original"`
	Preview    datatypes.JSON `json:"preview,omitempty" gorm:"column:preview"`
	Text       datatypes.JSON `json:"text,omitempty" gorm:"column:text"`
	Ocr        datatypes.JSON `json:"ocr,omitempty" gorm:"column:ocr"`
	Thumbnail  datatypes.JSON `json:"thumbnail,omitempty" gorm:"column:thumbnail"`
	CreateTime string         `json:"createTime" gorm:"column:create_time"`
	UpdateTime *string        `json:"updateTime,omitempty" gorm:"column:update_time"`
}

func (PostgresSnapshot) TableName() string {
	return "snapshot"
}

func (s *PostgresSnapshot) BeforeCreate(tx *gorm.DB) (err error) {
	s.CreateTime = time.Now().UTC().Format(time.RFC3339)
	return nil
}

func (s *PostgresSnapshot) BeforeSave(tx *gorm.DB) (err error) {
	timeNow := time.Now().UTC().Format(time.RFC3339)
	s.UpdateTime = &timeNow
	return nil
}

func (s PostgresSnapshot) GetId() string {
	return s.Id
}

func (s PostgresSnapshot) GetVersion() int64 {
	return s.Version
}

func (s PostgresSnapshot) GetOriginal() *model.S3Object {
	if s.Original.String() == "" {
		return nil
	}
	var res = model.S3Object{}
	if err := json.Unmarshal([]byte(s.Original.String()), &res); err != nil {
		log.Fatal(err)
		return nil
	}
	return &res
}

func (s PostgresSnapshot) GetPreview() *model.S3Object {
	if s.Preview.String() == "" {
		return nil
	}
	var res = model.S3Object{}
	if err := json.Unmarshal([]byte(s.Preview.String()), &res); err != nil {
		log.Fatal(err)
		return nil
	}
	return &res
}

func (s PostgresSnapshot) GetText() *model.S3Object {
	if s.Text.String() == "" {
		return nil
	}
	var res = model.S3Object{}
	if err := json.Unmarshal([]byte(s.Text.String()), &res); err != nil {
		log.Fatal(err)
		return nil
	}
	return &res
}

func (s PostgresSnapshot) GetOcr() *model.S3Object {
	if s.Ocr.String() == "" {
		return nil
	}
	var res = model.S3Object{}
	if err := json.Unmarshal([]byte(s.Ocr.String()), &res); err != nil {
		log.Fatal(err)
		return nil
	}
	return &res
}

func (s PostgresSnapshot) GetThumbnail() *model.Thumbnail {
	if s.Thumbnail.String() == "" {
		return nil
	}
	var res = model.Thumbnail{}
	if err := json.Unmarshal([]byte(s.Thumbnail.String()), &res); err != nil {
		log.Fatal(err)
		return nil
	}
	return &res
}

func (s *PostgresSnapshot) SetOriginal(m *model.S3Object) {
	b, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
		return
	}
	if err := s.Original.UnmarshalJSON(b); err != nil {
		log.Fatal(err)
	}
}

func (s *PostgresSnapshot) SetPreview(m *model.S3Object) {
	b, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
		return
	}
	if err := s.Preview.UnmarshalJSON(b); err != nil {
		log.Fatal(err)
	}
}

func (s *PostgresSnapshot) SetText(m *model.S3Object) {
	b, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
		return
	}
	if err := s.Text.UnmarshalJSON(b); err != nil {
		log.Fatal(err)
	}
}

func (s *PostgresSnapshot) SetOcr(m *model.S3Object) {
	b, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
		return
	}
	if err := s.Ocr.UnmarshalJSON(b); err != nil {
		log.Fatal(err)
	}
}

func (s *PostgresSnapshot) SetThumbnail(m *model.Thumbnail) {
	b, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
		return
	}
	if err := s.Thumbnail.UnmarshalJSON(b); err != nil {
		log.Fatal(err)
	}
}

func (s PostgresSnapshot) HasOriginal() bool {
	return s.Original != nil
}

func (s PostgresSnapshot) HasPreview() bool {
	return s.Preview != nil
}

func (s PostgresSnapshot) HasText() bool {
	return s.Text != nil
}

func (s PostgresSnapshot) HasOcr() bool {
	return s.Ocr != nil
}

func (s PostgresSnapshot) HasThumbnail() bool {
	return s.Thumbnail != nil
}

func (s PostgresSnapshot) GetCreateTime() string {
	return s.CreateTime
}

func (s PostgresSnapshot) GetUpdateTime() *string {
	return s.UpdateTime
}

type PostgresSnapshotRepo struct {
	db *gorm.DB
}

func NewPostgresSnapshotRepo() *PostgresSnapshotRepo {
	return &PostgresSnapshotRepo{
		db: infra.GetDb(),
	}
}

func (repo *PostgresSnapshotRepo) find(id string) (*PostgresSnapshot, error) {
	var res PostgresSnapshot
	if db := repo.db.Where("id = ?", id).First(&res); db.Error != nil {
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			return nil, errorpkg.NewSnapshotNotFoundError(db.Error)
		} else {
			return nil, errorpkg.NewInternalServerError(db.Error)
		}
	}
	return &res, nil
}

func (repo *PostgresSnapshotRepo) Find(id string) (model.SnapshotModel, error) {
	res, err := repo.find(id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (repo *PostgresSnapshotRepo) Save(snapshot model.SnapshotModel) error {
	if db := repo.db.Save(snapshot); db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *PostgresSnapshotRepo) MapWithFile(id string, fileId string) error {
	if db := repo.db.Exec("INSERT INTO snapshot_file (snapshot_id, file_id) VALUES (?, ?)", id, fileId); db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *PostgresSnapshotRepo) DeleteMappingsForFile(fileId string) error {
	if db := repo.db.Exec("DELETE FROM snapshot_file WHERE file_id = ?", fileId); db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *PostgresSnapshotRepo) FindAllForFile(fileId string) ([]*PostgresSnapshot, error) {
	var res []*PostgresSnapshot
	db := repo.db.
		Raw("SELECT * FROM snapshot s LEFT JOIN snapshot_file sf ON s.id = sf.snapshot_id WHERE sf.file_id = ? ORDER BY s.version", fileId).
		Scan(&res)
	if db.Error != nil {
		return nil, db.Error
	}
	return res, nil
}

func (repo *PostgresSnapshotRepo) FindAllDangling() ([]model.SnapshotModel, error) {
	var snapshots []*PostgresSnapshot
	db := repo.db.Raw("SELECT * FROM snapshot s LEFT JOIN snapshot_file sf ON s.id = sf.snapshot_id WHERE sf.snapshot_id IS NULL").Scan(&snapshots)
	if db.Error != nil {
		return nil, db.Error
	}
	var res []model.SnapshotModel
	for _, s := range snapshots {
		res = append(res, s)
	}
	return res, nil
}

func (repo *PostgresSnapshotRepo) DeleteAllDangling() error {
	if db := repo.db.Exec("DELETE FROM snapshot WHERE id IN (SELECT s.id FROM (SELECT * FROM snapshot) s LEFT JOIN snapshot_file sf ON s.id = sf.snapshot_id WHERE sf.snapshot_id IS NULL)"); db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *PostgresSnapshotRepo) GetLatestVersionForFile(fileId string) (int64, error) {
	type Result struct {
		Result int64
	}
	var res Result
	if db := repo.db.
		Raw("SELECT coalesce(max(s.version), 0) + 1 result FROM snapshot s LEFT JOIN snapshot_file map ON s.id = map.snapshot_id WHERE map.file_id = ?", fileId).
		Scan(&res); db.Error != nil {
		return 0, db.Error
	}
	return res.Result, nil
}
