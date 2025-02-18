// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/service"
	"github.com/kouprlabs/voltaserve/api/test/test_helper"
)

type TaskServiceSuite struct {
	suite.Suite
	users []model.User
}

func TestTaskServiceSuite(t *testing.T) {
	suite.Run(t, new(TaskServiceSuite))
}

func (s *TaskServiceSuite) SetupTest() {
	var err error
	s.users, err = test_helper.CreateUsers(2)
	if err != nil {
		s.Fail(err.Error())
		return
	}
}

func (s *TaskServiceSuite) TestCreate() {
	// Test creating a task with all fields
	task, err := service.NewTaskService().Create(service.TaskCreateOptions{
		Name:            "task A",
		Error:           nil,
		Percentage:      nil,
		IsIndeterminate: false,
		UserID:          s.users[0].GetID(),
		Status:          model.TaskStatusWaiting,
		Payload:         map[string]string{"key": "value"},
	})
	s.Require().NoError(err)
	s.NotNil(task)
	s.Equal("task A", task.Name)
	s.Equal(s.users[0].GetID(), task.UserID)
	s.Equal(model.TaskStatusWaiting, task.Status)
	s.Equal(map[string]string{"key": "value"}, task.Payload)

	// Test creating a task with minimal fields
	task, err = service.NewTaskService().Create(service.TaskCreateOptions{
		Name:   "task B",
		UserID: s.users[0].GetID(),
		Status: model.TaskStatusRunning,
	})
	s.Require().NoError(err)
	s.NotNil(task)
	s.Equal("task B", task.Name)
	s.Equal(s.users[0].GetID(), task.UserID)
	s.Equal(model.TaskStatusRunning, task.Status)
}

func (s *TaskServiceSuite) TestPatch() {
	// Create a task to patch
	task, err := service.NewTaskService().Create(service.TaskCreateOptions{
		Name:   "task",
		UserID: s.users[0].GetID(),
		Status: model.TaskStatusWaiting,
	})
	s.Require().NoError(err)

	// Patch the task's name and status
	task, err = service.NewTaskService().Patch(task.ID, service.TaskPatchOptions{
		Fields: []string{service.TaskFieldName, service.TaskFieldStatus},
		Name:   helper.ToPtr("task (edit)"),
		Status: helper.ToPtr(model.TaskStatusRunning),
	})
	s.Require().NoError(err)
	s.Equal("task (edit)", task.Name)
	s.Equal(model.TaskStatusRunning, task.Status)

	// Patch the task's error and percentage
	task, err = service.NewTaskService().Patch(task.ID, service.TaskPatchOptions{
		Fields:     []string{service.TaskFieldError, service.TaskFieldPercentage},
		Error:      helper.ToPtr("something went wrong"),
		Percentage: helper.ToPtr(50),
	})
	s.Require().NoError(err)
	s.Equal("something went wrong", *task.Error)
	s.Equal(50, *task.Percentage)

	// Patch the task's payload
	task, err = service.NewTaskService().Patch(task.ID, service.TaskPatchOptions{
		Fields:  []string{service.TaskFieldPayload},
		Payload: map[string]string{"newKey": "newValue"},
	})
	s.Require().NoError(err)
	s.Equal(map[string]string{"newKey": "newValue"}, task.Payload)
}

func (s *TaskServiceSuite) TestFind() {
	// Create a task to find
	task, err := service.NewTaskService().Create(service.TaskCreateOptions{
		Name:   "task",
		UserID: s.users[0].GetID(),
		Status: model.TaskStatusWaiting,
	})
	s.Require().NoError(err)

	// Find the task
	foundTask, err := service.NewTaskService().Find(task.ID, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(task.ID, foundTask.ID)
	s.Equal("task", foundTask.Name)

	// Try to find a task that doesn't belong to the user
	_, err = service.NewTaskService().Find(task.ID, s.users[1].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewTaskNotFoundError(nil).Error(), err.Error())
}

func (s *TaskServiceSuite) TestList() {
	// Create multiple tasks for listing
	_, err := service.NewTaskService().Create(service.TaskCreateOptions{
		Name:   "task A",
		UserID: s.users[0].GetID(),
		Status: model.TaskStatusWaiting,
	})
	s.Require().NoError(err)
	_, err = service.NewTaskService().Create(service.TaskCreateOptions{
		Name:   "task B",
		UserID: s.users[0].GetID(),
		Status: model.TaskStatusRunning,
	})
	s.Require().NoError(err)

	// List tasks with default options
	list, err := service.NewTaskService().List(service.TaskListOptions{Page: 1, Size: 10}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Len(list.Data, 2)
	s.Equal(uint64(2), list.TotalElements)

	// List tasks with sorting by name
	list, err = service.NewTaskService().List(service.TaskListOptions{
		Page:      1,
		Size:      10,
		SortBy:    service.TaskSortByName,
		SortOrder: service.TaskSortOrderDesc,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal("task B", list.Data[0].Name)
	s.Equal("task A", list.Data[1].Name)

	// List tasks with pagination
	list, err = service.NewTaskService().List(service.TaskListOptions{Page: 1, Size: 1}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Len(list.Data, 1)
	s.Equal(uint64(2), list.TotalElements)
}

func (s *TaskServiceSuite) TestProbe() {
	// Create multiple tasks for probing
	_, err := service.NewTaskService().Create(service.TaskCreateOptions{
		Name:   "task A",
		UserID: s.users[0].GetID(),
		Status: model.TaskStatusWaiting,
	})
	s.Require().NoError(err)

	_, err = service.NewTaskService().Create(service.TaskCreateOptions{
		Name:   "task B",
		UserID: s.users[0].GetID(),
		Status: model.TaskStatusRunning,
	})
	s.Require().NoError(err)

	// Probe tasks
	probe, err := service.NewTaskService().Probe(service.TaskListOptions{Page: 1, Size: 10}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(2), probe.TotalElements)
	s.Equal(uint64(1), probe.TotalPages)
}

func (s *TaskServiceSuite) TestDismiss() {
	// Create a task with an error to dismiss
	task, err := service.NewTaskService().Create(service.TaskCreateOptions{
		Name:   "task with error",
		UserID: s.users[0].GetID(),
		Status: model.TaskStatusError,
		Error:  helper.ToPtr("something went wrong"),
	})
	s.Require().NoError(err)

	// Dismiss the task
	err = service.NewTaskService().Dismiss(task.ID, s.users[0].GetID())
	s.Require().NoError(err)

	// Create another task with error
	task, err = service.NewTaskService().Create(service.TaskCreateOptions{
		Name:   "task with error",
		UserID: s.users[0].GetID(),
		Status: model.TaskStatusError,
		Error:  helper.ToPtr("something went wrong"),
	})
	s.Require().NoError(err)

	// Try to dismiss a task that doesn't belong to the user
	err = service.NewTaskService().Dismiss(task.ID, s.users[1].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewTaskBelongsToAnotherUserError(nil).Error(), err.Error())

	// Create another task with status running
	task, err = service.NewTaskService().Create(service.TaskCreateOptions{
		Name:   "task with status running",
		UserID: s.users[0].GetID(),
		Status: model.TaskStatusRunning,
	})
	s.Require().NoError(err)

	// Try to dismiss a task that is still running
	err = service.NewTaskService().Dismiss(task.ID, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewTaskIsRunningError(nil).Error(), err.Error())
}

func (s *TaskServiceSuite) TestDismissAll() {
	// Create multiple tasks with errors to dismiss
	_, err := service.NewTaskService().Create(service.TaskCreateOptions{
		Name:   "task A",
		UserID: s.users[0].GetID(),
		Status: model.TaskStatusError,
		Error:  helper.ToPtr("error A"),
	})
	s.Require().NoError(err)
	_, err = service.NewTaskService().Create(service.TaskCreateOptions{
		Name:   "task B",
		UserID: s.users[0].GetID(),
		Status: model.TaskStatusError,
		Error:  helper.ToPtr("error B"),
	})
	s.Require().NoError(err)

	// Dismiss all tasks
	dismissAllResult, err := service.NewTaskService().DismissAll(s.users[0].GetID())
	s.Require().NoError(err)
	s.Len(dismissAllResult.Succeeded, 2)
	s.Empty(dismissAllResult.Failed)

	// Verify that the tasks are dismissed
	list, err := service.NewTaskService().List(service.TaskListOptions{Page: 1, Size: 10}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Empty(list.Data)
}

func (s *TaskServiceSuite) TestDelete() {
	// Create a task to delete
	task, err := service.NewTaskService().Create(service.TaskCreateOptions{
		Name:   "task",
		UserID: s.users[0].GetID(),
		Status: model.TaskStatusWaiting,
	})
	s.Require().NoError(err)

	// Delete the task
	err = service.NewTaskService().Delete(task.ID)
	s.Require().NoError(err)

	// Verify that the task is deleted
	_, err = service.NewTaskService().Find(task.ID, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewTaskNotFoundError(nil).Error(), err.Error())
}

func (s *TaskServiceSuite) TestCount() {
	// Create multiple tasks for counting
	_, err := service.NewTaskService().Create(service.TaskCreateOptions{
		Name:   "task A",
		UserID: s.users[0].GetID(),
		Status: model.TaskStatusWaiting,
	})
	s.Require().NoError(err)
	_, err = service.NewTaskService().Create(service.TaskCreateOptions{
		Name:   "task B",
		UserID: s.users[0].GetID(),
		Status: model.TaskStatusRunning,
	})
	s.Require().NoError(err)

	// Count tasks
	count, err := service.NewTaskService().Count(s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(int64(2), *count)
}
