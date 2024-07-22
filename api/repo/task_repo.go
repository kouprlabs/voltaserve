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
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/infra"
	"github.com/kouprlabs/voltaserve/api/log"
	"github.com/kouprlabs/voltaserve/api/model"
)

type taskEntity struct {
	ID              string         `gorm:"column:id"               json:"id"`
	Name            string         `gorm:"column:name"             json:"name"`
	Error           *string        `gorm:"column:error"            json:"error,omitempty"`
	Percentage      *int           `gorm:"column:percentage"       json:"percentage,omitempty"`
	IsIndeterminate bool           `gorm:"column:is_indeterminate" json:"isIndeterminate"`
	UserID          string         `gorm:"column:user_id"          json:"userId"`
	Status          string         `gorm:"column:status"           json:"status"`
	Payload         datatypes.JSON `gorm:"column:payload"          json:"payload"`
	CreateTime      string         `gorm:"column:create_time"      json:"createTime"`
	UpdateTime      *string        `gorm:"column:update_time"      json:"updateTime,omitempty"`
}

func (*taskEntity) TableName() string {
	return "task"
}

func (o *taskEntity) BeforeCreate(*gorm.DB) (err error) {
	o.CreateTime = time.Now().UTC().Format(time.RFC3339)
	return nil
}

func (o *taskEntity) BeforeSave(*gorm.DB) (err error) {
	timeNow := time.Now().UTC().Format(time.RFC3339)
	o.UpdateTime = &timeNow
	return nil
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

func (p *taskEntity) GetIsIndeterminate() bool {
	return p.IsIndeterminate
}

func (p *taskEntity) GetUserID() string {
	return p.UserID
}

func (p *taskEntity) GetStatus() string {
	return p.Status
}

func (s *taskEntity) GetPayload() map[string]string {
	if s.Payload.String() == "" {
		return nil
	}
	res := map[string]string{}
	if err := json.Unmarshal([]byte(s.Payload.String()), &res); err != nil {
		log.GetLogger().Fatal(err)
		return nil
	}
	return res
}

func (o *taskEntity) GetCreateTime() string {
	return o.CreateTime
}

func (o *taskEntity) GetUpdateTime() *string {
	return o.UpdateTime
}

func (p *taskEntity) HasError() bool {
	return p.Error != nil
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

func (p *taskEntity) SetIsIndeterminate(isIndeterminate bool) {
	p.IsIndeterminate = isIndeterminate
}

func (p *taskEntity) SetUserID(userID string) {
	p.UserID = userID
}

func (p *taskEntity) SetStatus(status string) {
	p.Status = status
}

func (s *taskEntity) SetPayload(p map[string]string) {
	if p == nil {
		s.Payload = nil
	} else {
		b, err := json.Marshal(p)
		if err != nil {
			log.GetLogger().Fatal(err)
			return
		}
		if err := s.Payload.UnmarshalJSON(b); err != nil {
			log.GetLogger().Fatal(err)
		}
	}
}

type TaskRepo interface {
	Insert(opts TaskInsertOptions) (model.Task, error)
	Find(id string) (model.Task, error)
	Count() (int64, error)
	CountByEmail(email string) (int64, error)
	GetIDs(userID string) ([]string, error)
	Save(task model.Task) error
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
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	Error           *string           `json:"error,omitempty"`
	Percentage      *int              `json:"percentage,omitempty"`
	IsIndeterminate bool              `json:"isIndeterminate"`
	UserID          string            `json:"userId"`
	Status          string            `json:"status"`
	Payload         map[string]string `json:"payload,omitempty"`
}

func (repo *taskRepo) Insert(opts TaskInsertOptions) (model.Task, error) {
	task := taskEntity{
		ID:              opts.ID,
		Name:            opts.Name,
		Error:           opts.Error,
		Percentage:      opts.Percentage,
		IsIndeterminate: opts.IsIndeterminate,
		UserID:          opts.UserID,
		Status:          opts.Status,
	}
	if opts.Payload != nil {
		task.SetPayload(opts.Payload)
	}
	if db := repo.db.Create(&task); db.Error != nil {
		return nil, db.Error
	}
	res, err := repo.Find(opts.ID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (repo *taskRepo) find(id string) (*taskEntity, error) {
	res := taskEntity{}
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

func (repo *taskRepo) Count() (int64, error) {
	type Result struct {
		Result int64
	}
	var res Result
	db := repo.db.
		Raw("SELECT count(*) as result FROM task").
		Scan(&res)
	if db.Error != nil {
		return 0, db.Error
	}
	return res.Result, nil
}

func (repo *taskRepo) GetIDs(userID string) ([]string, error) {
	type Value struct {
		Result string
	}
	var values []Value
	db := repo.db.
		Raw("SELECT id result FROM task WHERE user_id = ? ORDER BY create_time DESC", userID).
		Scan(&values)
	if db.Error != nil {
		return []string{}, db.Error
	}
	res := []string{}
	for _, v := range values {
		res = append(res, v.Result)
	}
	return res, nil
}

func (repo *taskRepo) CountByEmail(userID string) (int64, error) {
	var count int64
	db := repo.db.
		Model(&taskEntity{}).
		Where("user_id = ?", userID).
		Count(&count)
	if db.Error != nil {
		return 0, db.Error
	}
	return count, nil
}

func (repo *taskRepo) Save(task model.Task) error {
	db := repo.db.Save(task)
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
