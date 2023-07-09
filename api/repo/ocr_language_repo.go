package repo

import (
	"voltaserve/infra"
	"voltaserve/model"

	"gorm.io/gorm"
)

type OCRLanguageRepo interface {
	FindAll() ([]model.OCRLanguage, error)
}

func NewOCRLanguageRepo() OCRLanguageRepo {
	return newOCRLanguageRepo()
}

type ocrLanguageEntity struct {
	ID        string `json:"id" gorm:"column:id"`
	ISO639Pt3 string `json:"iso639pt3" gorm:"column:iso639_3"`
	Name      string `json:"name" gorm:"column:name"`
}

func (ocrLanguageEntity) TableName() string {
	return "ocrlanguage"
}

func (ol ocrLanguageEntity) GetID() string {
	return ol.ID
}

func (ol ocrLanguageEntity) GetISO639Pt3() string {
	return ol.ISO639Pt3
}

func (ol ocrLanguageEntity) GetName() string {
	return ol.Name
}

type ocrLanguageRepo struct {
	db *gorm.DB
}

func newOCRLanguageRepo() *ocrLanguageRepo {
	return &ocrLanguageRepo{
		db: infra.GetDb(),
	}
}

func (repo *ocrLanguageRepo) FindAll() ([]model.OCRLanguage, error) {
	var entities []*ocrLanguageEntity
	db := repo.db.Raw(`select * from "ocrlanguage"`).Scan(&entities)
	if db.Error != nil {
		return nil, db.Error
	}
	var res []model.OCRLanguage
	for _, ol := range entities {
		res = append(res, ol)
	}
	return res, nil
}
