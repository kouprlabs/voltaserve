package repo

import (
	"errors"
	"voltaserve/errorpkg"
	"voltaserve/infra"
	"voltaserve/model"

	"gorm.io/gorm"
)

type processEntity struct {
	ID              string  `json:"id" gorm:"column:id"`
	Description     string  `json:"description" gorm:"column:description"`
	Error           *string `json:"error" gorm:"column:error"`
	Percentage      *int    `json:"percentage" gorm:"column:percentage"`
	IsComplete      bool    `json:"isComplete" gorm:"column:is_complete"`
	IsIndeterminate bool    `json:"isIndeterminate" gorm:"column:is_indeterminate"`
}

func (*processEntity) TableName() string {
	return "process"
}

func (p *processEntity) GetID() string {
	return p.ID
}

func (p *processEntity) GetDescription() string {
	return p.Description
}

func (p *processEntity) GetError() *string {
	return p.Error
}

func (p *processEntity) GetPercentage() *int {
	return p.Percentage
}

func (p *processEntity) GetIsComplete() bool {
	return p.IsComplete
}

func (p *processEntity) GetIsIndeterminate() bool {
	return p.IsIndeterminate
}

func (p *processEntity) SetDescription(description string) {
	p.Description = description
}

func (p *processEntity) SetError(error *string) {
	p.Error = error
}

func (p *processEntity) SetPercentage(percentage *int) {
	p.Percentage = percentage
}

func (p *processEntity) SetIsComplete(isComplete bool) {
	p.IsComplete = isComplete
}

func (p *processEntity) SetIsIndeterminate(isIndeterminate bool) {
	p.IsIndeterminate = isIndeterminate
}

type ProcessRepo interface {
	Insert(opts ProcessInsertOptions) (model.Process, error)
	Find(id string) (model.Process, error)
	Save(org model.Process) error
	Delete(id string) error
}

func NewProcessRepo() ProcessRepo {
	return newProcessRepo()
}

func NewProcess() model.Process {
	return &processEntity{}
}

type processRepo struct {
	db *gorm.DB
}

func newProcessRepo() *processRepo {
	return &processRepo{
		db: infra.NewPostgresManager().GetDBOrPanic(),
	}
}

type ProcessInsertOptions struct {
	ID              string
	Name            string
	Description     string
	Error           *string
	Percentage      *int
	IsComplete      bool
	IsIndeterminate bool
}

func (repo *processRepo) Insert(opts ProcessInsertOptions) (model.Process, error) {
	org := processEntity{
		ID:              opts.ID,
		Description:     opts.Description,
		Error:           opts.Error,
		Percentage:      opts.Percentage,
		IsComplete:      opts.IsComplete,
		IsIndeterminate: opts.IsIndeterminate,
	}
	if db := repo.db.Create(&org); db.Error != nil {
		return nil, db.Error
	}
	res, err := repo.Find(opts.ID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (repo *processRepo) find(id string) (*processEntity, error) {
	var res = processEntity{}
	db := repo.db.Where("id = ?", id).First(&res)
	if db.Error != nil {
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			return nil, errorpkg.NewProcessNotFoundError(db.Error)
		} else {
			return nil, errorpkg.NewInternalServerError(db.Error)
		}
	}
	return &res, nil
}

func (repo *processRepo) Find(id string) (model.Process, error) {
	process, err := repo.find(id)
	if err != nil {
		return nil, err
	}
	return process, nil
}

func (repo *processRepo) Save(org model.Process) error {
	db := repo.db.Save(org)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *processRepo) Delete(id string) error {
	db := repo.db.Exec("DELETE FROM process WHERE id = ?", id)
	if db.Error != nil {
		return db.Error
	}
	return nil
}
