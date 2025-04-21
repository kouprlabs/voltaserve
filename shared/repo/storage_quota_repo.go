package repo

import (
	"errors"

	"gorm.io/gorm"

	"github.com/kouprlabs/voltaserve/shared/config"
	"github.com/kouprlabs/voltaserve/shared/errorpkg"
	"github.com/kouprlabs/voltaserve/shared/helper"
	"github.com/kouprlabs/voltaserve/shared/infra"
	"github.com/kouprlabs/voltaserve/shared/model"
)

type storageQuotaEntity struct {
	ID              string  `gorm:"column:id;size:36"       json:"id"`
	UserID          string  `gorm:"column:user_id"          json:"userId"`
	StorageCapacity int64   `gorm:"column:storage_capacity" json:"storageCapacity"`
	CreateTime      string  `gorm:"column:create_time"      json:"createTime"`
	UpdateTime      *string `gorm:"column:update_time"      json:"updateTime,omitempty"`
}

func (*storageQuotaEntity) TableName() string {
	return "storage_quota"
}

func (sq *storageQuotaEntity) BeforeCreate(*gorm.DB) (err error) {
	sq.CreateTime = helper.NewTimeString()
	return nil
}

func (sq *storageQuotaEntity) BeforeSave(*gorm.DB) (err error) {
	sq.UpdateTime = helper.ToPtr(helper.NewTimeString())
	return nil
}

func (sq *storageQuotaEntity) GetID() string {
	return sq.ID
}

func (sq *storageQuotaEntity) GetUserID() string {
	return sq.UserID
}

func (sq *storageQuotaEntity) GetStorageCapacity() int64 {
	return sq.StorageCapacity
}

func (sq *storageQuotaEntity) GetCreateTime() string {
	return sq.CreateTime
}

func (sq *storageQuotaEntity) GetUpdateTime() *string {
	return sq.UpdateTime
}

func (sq *storageQuotaEntity) SetID(id string) {
	sq.ID = id
}

func (sq *storageQuotaEntity) SetUserID(userId string) {
	sq.UserID = userId
}

func (sq *storageQuotaEntity) SetStorageCapacity(storageCapacity int64) {
	sq.StorageCapacity = storageCapacity
}

func (sq *storageQuotaEntity) SetCreateTime(createTime string) {
	sq.CreateTime = createTime
}

func (sq *storageQuotaEntity) SetUpdateTime(updateTime *string) {
	sq.UpdateTime = updateTime
}

func NewStorageQuotaModel() model.StorageQuota {
	return &storageQuotaEntity{}
}

type StorageQuotaRepo struct {
	db *gorm.DB
}

func NewStorageQuotaRepo(postgres config.PostgresConfig, environment config.EnvironmentConfig) *StorageQuotaRepo {
	return &StorageQuotaRepo{
		db: infra.NewPostgresManager(postgres, environment).GetDBOrPanic(),
	}
}

func (repo *StorageQuotaRepo) Insert(storageQuota model.StorageQuota) (model.StorageQuota, error) {
	if db := repo.db.Create(storageQuota); db.Error != nil {
		return nil, db.Error
	}
	res, err := repo.FindByUserID(storageQuota.GetUserID())
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (repo *StorageQuotaRepo) FindByUserID(userID string) (model.StorageQuota, error) {
	res, err := repo.findByUserID(userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (repo *StorageQuotaRepo) Save(storageQuota model.StorageQuota) error {
	db := repo.db.Save(storageQuota)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *StorageQuotaRepo) DeleteByUserID(userID string) error {
	db := repo.db.Exec(`DELETE FROM storage_quota WHERE user_id = ?`, userID)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *StorageQuotaRepo) findByUserID(userID string) (*storageQuotaEntity, error) {
	res := storageQuotaEntity{}
	db := repo.db.Where("user_id = ?", userID).First(&res)
	if db.Error != nil {
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			return nil, errorpkg.NewResourceNotFoundError(db.Error)
		} else {
			return nil, errorpkg.NewInternalServerError(db.Error)
		}
	}
	return &res, nil
}
