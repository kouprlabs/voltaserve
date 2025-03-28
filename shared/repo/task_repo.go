// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package repo

import (
	"encoding/json"
	"errors"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/kouprlabs/voltaserve/shared/config"
	"github.com/kouprlabs/voltaserve/shared/errorpkg"
	"github.com/kouprlabs/voltaserve/shared/helper"
	"github.com/kouprlabs/voltaserve/shared/infra"
	"github.com/kouprlabs/voltaserve/shared/logger"
	"github.com/kouprlabs/voltaserve/shared/model"
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

func (e *taskEntity) BeforeCreate(*gorm.DB) (err error) {
	e.CreateTime = helper.NewTimeString()
	return nil
}

func (e *taskEntity) BeforeSave(*gorm.DB) (err error) {
	e.UpdateTime = helper.ToPtr(helper.NewTimeString())
	return nil
}

func (e *taskEntity) GetID() string {
	return e.ID
}

func (e *taskEntity) GetName() string {
	return e.Name
}

func (e *taskEntity) GetError() *string {
	return e.Error
}

func (e *taskEntity) GetPercentage() *int {
	return e.Percentage
}

func (e *taskEntity) GetIsIndeterminate() bool {
	return e.IsIndeterminate
}

func (e *taskEntity) GetUserID() string {
	return e.UserID
}

func (e *taskEntity) GetStatus() string {
	return e.Status
}

func (e *taskEntity) GetPayload() map[string]string {
	if e.Payload.String() == "" {
		return nil
	}
	res := map[string]string{}
	if err := json.Unmarshal([]byte(e.Payload.String()), &res); err != nil {
		logger.GetLogger().Fatal(err)
		return nil
	}
	return res
}

func (e *taskEntity) GetCreateTime() string {
	return e.CreateTime
}

func (e *taskEntity) GetUpdateTime() *string {
	return e.UpdateTime
}

func (e *taskEntity) HasError() bool {
	return e.Error != nil
}

func (e *taskEntity) SetName(name string) {
	e.Name = name
}

func (e *taskEntity) SetError(error *string) {
	e.Error = error
}

func (e *taskEntity) SetPercentage(percentage *int) {
	e.Percentage = percentage
}

func (e *taskEntity) SetIsIndeterminate(isIndeterminate bool) {
	e.IsIndeterminate = isIndeterminate
}

func (e *taskEntity) SetID(id string) {
	e.ID = id
}

func (e *taskEntity) SetUserID(userID string) {
	e.UserID = userID
}

func (e *taskEntity) SetStatus(status string) {
	e.Status = status
}

func (e *taskEntity) SetPayload(p map[string]string) {
	if p == nil {
		e.Payload = nil
	} else {
		b, err := json.Marshal(p)
		if err != nil {
			logger.GetLogger().Fatal(err)
			return
		}
		if err := e.Payload.UnmarshalJSON(b); err != nil {
			logger.GetLogger().Fatal(err)
		}
	}
}

func (e *taskEntity) SetCreateTime(createTime string) {
	e.CreateTime = createTime
}

func (e *taskEntity) SetUpdateTime(updateTime *string) {
	e.UpdateTime = updateTime
}

func NewTaskModel() model.Task {
	return &taskEntity{}
}

type TaskNewModelOptions struct {
	ID              string
	Name            string
	Error           *string
	Percentage      *int
	IsIndeterminate bool
	UserID          string
	Status          string
	Payload         map[string]string
	CreateTime      string
	UpdateTime      *string
}

func NewTaskModelWithOptions(opts TaskNewModelOptions) model.Task {
	res := &taskEntity{
		ID:              opts.ID,
		Name:            opts.Name,
		Error:           opts.Error,
		Percentage:      opts.Percentage,
		IsIndeterminate: opts.IsIndeterminate,
		UserID:          opts.UserID,
		Status:          opts.Status,
		CreateTime:      opts.CreateTime,
		UpdateTime:      opts.UpdateTime,
	}
	res.SetPayload(opts.Payload)
	return res
}

type TaskRepo struct {
	db *gorm.DB
}

func NewTaskRepo(postgres config.PostgresConfig, environment config.EnvironmentConfig) *TaskRepo {
	return &TaskRepo{
		db: infra.NewPostgresManager(postgres, environment).GetDBOrPanic(),
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

const TaskPayloadObjectKey = "object"

func (repo *TaskRepo) Insert(opts TaskInsertOptions) (model.Task, error) {
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

func (repo *TaskRepo) Find(id string) (model.Task, error) {
	res, err := repo.find(id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (repo *TaskRepo) FindOrNil(id string) model.Task {
	res, err := repo.Find(id)
	if err != nil {
		return nil
	}
	return res
}

func (repo *TaskRepo) FindIDs(userID string) ([]string, error) {
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
	res := make([]string, 0)
	for _, v := range values {
		res = append(res, v.Result)
	}
	return res, nil
}

func (repo *TaskRepo) FindIDsByOwner(userID string) ([]string, error) {
	type IDResult struct {
		Result string
	}
	var ids []IDResult
	db := repo.db.Raw(`SELECT id result FROM task WHERE user_id = ?`, userID).Scan(&ids)
	if db.Error != nil {
		return nil, db.Error
	}
	res := make([]string, 0)
	for _, id := range ids {
		res = append(res, id.Result)
	}
	return res, nil
}

func (repo *TaskRepo) Count() (int64, error) {
	var count int64
	db := repo.db.Model(&taskEntity{}).Count(&count)
	if db.Error != nil {
		return -1, db.Error
	}
	return count, nil
}

func (repo *TaskRepo) CountByUserID(userID string) (int64, error) {
	var count int64
	db := repo.db.
		Model(&taskEntity{}).
		Where("user_id = ?", userID).
		Where("status = ? OR status = ?", model.TaskStatusWaiting, model.TaskStatusRunning).
		Count(&count)
	if db.Error != nil {
		return -1, db.Error
	}
	return count, nil
}

func (repo *TaskRepo) Save(task model.Task) error {
	db := repo.db.Save(task)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *TaskRepo) Delete(id string) error {
	db := repo.db.Exec("DELETE FROM task WHERE id = ?", id)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *TaskRepo) find(id string) (*taskEntity, error) {
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
