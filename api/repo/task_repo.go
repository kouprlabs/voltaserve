package repo

import (
	"errors"
	"voltaserve/errorpkg"
	"voltaserve/infra"
	"voltaserve/model"

	"gorm.io/gorm"
)

type taskEntity struct {
	ID              string  `json:"id" gorm:"column:id"`
	Name            string  `json:"name" gorm:"column:name"`
	Error           *string `json:"error" gorm:"column:error"`
	Percentage      *int    `json:"percentage" gorm:"column:percentage"`
	IsComplete      bool    `json:"isComplete" gorm:"column:is_complete"`
	IsIndeterminate bool    `json:"isIndeterminate" gorm:"column:is_indeterminate"`
	UserID          string  `json:"userId" gorm:"column:user_id"`
	CreateTime      string  `json:"createTime" gorm:"column:create_time"`
	UpdateTime      *string `json:"updateTime,omitempty" gorm:"column:update_time"`
}

func (*taskEntity) TableName() string {
	return "task"
}

func (p *taskEntity) GetID() string {
	return p.ID
}

func (p *taskEntity) GetName() string {
	return p.Name
}

func (p *taskEntity) GetError() *string {
	return p.Error
}

func (p *taskEntity) GetPercentage() *int {
	return p.Percentage
}

func (p *taskEntity) GetIsComplete() bool {
	return p.IsComplete
}

func (p *taskEntity) GetIsIndeterminate() bool {
	return p.IsIndeterminate
}

func (p *taskEntity) GetUserID() string {
	return p.UserID
}

func (o *taskEntity) GetCreateTime() string {
	return o.CreateTime
}

func (o *taskEntity) GetUpdateTime() *string {
	return o.UpdateTime
}

func (p *taskEntity) SetName(name string) {
	p.Name = name
}

func (p *taskEntity) SetError(error *string) {
	p.Error = error
}

func (p *taskEntity) SetPercentage(percentage *int) {
	p.Percentage = percentage
}

func (p *taskEntity) SetIsComplete(isComplete bool) {
	p.IsComplete = isComplete
}

func (p *taskEntity) SetIsIndeterminate(isIndeterminate bool) {
	p.IsIndeterminate = isIndeterminate
}

func (p *taskEntity) SetUserID(userID string) {
	p.UserID = userID
}

type TaskRepo interface {
	Insert(opts TaskInsertOptions) (model.Task, error)
	Find(id string) (model.Task, error)
	GetIDs() ([]string, error)
	Save(org model.Task) error
	Delete(id string) error
}

func NewTaskRepo() TaskRepo {
	return newTaskRepo()
}

func NewTask() model.Task {
	return &taskEntity{}
}

type taskRepo struct {
	db *gorm.DB
}

func newTaskRepo() *taskRepo {
	return &taskRepo{
		db: infra.NewPostgresManager().GetDBOrPanic(),
	}
}

type TaskInsertOptions struct {
	ID              string
	Name            string
	Description     string
	Error           *string
	Percentage      *int
	IsComplete      bool
	IsIndeterminate bool
}

func (repo *taskRepo) Insert(opts TaskInsertOptions) (model.Task, error) {
	org := taskEntity{
		ID:              opts.ID,
		Name:            opts.Name,
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

func (repo *taskRepo) find(id string) (*taskEntity, error) {
	var res = taskEntity{}
	db := repo.db.Where("id = ?", id).First(&res)
	if db.Error != nil {
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			return nil, errorpkg.NewTaskNotFoundError(db.Error)
		} else {
			return nil, errorpkg.NewInternalServerError(db.Error)
		}
	}
	return &res, nil
}

func (repo *taskRepo) Find(id string) (model.Task, error) {
	res, err := repo.find(id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (repo *taskRepo) GetIDs() ([]string, error) {
	type Value struct {
		Result string
	}
	var values []Value
	db := repo.db.Raw("SELECT id result FROM task ORDER BY create_time DESC").Scan(&values)
	if db.Error != nil {
		return []string{}, db.Error
	}
	res := []string{}
	for _, v := range values {
		res = append(res, v.Result)
	}
	return res, nil
}

func (repo *taskRepo) Save(org model.Task) error {
	db := repo.db.Save(org)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *taskRepo) Delete(id string) error {
	db := repo.db.Exec("DELETE FROM task WHERE id = ?", id)
	if db.Error != nil {
		return db.Error
	}
	return nil
}
