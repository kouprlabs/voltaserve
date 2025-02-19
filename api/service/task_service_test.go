// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package service_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/service"
	"github.com/kouprlabs/voltaserve/api/test"
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
	s.users, err = test.CreateUsers(2)
	if err != nil {
		s.Fail(err.Error())
		return
	}
}

func (s *TaskServiceSuite) TestCreate() {
	task, err := service.NewTaskService().Create(service.TaskCreateOptions{
		Name:    "task A",
		UserID:  s.users[0].GetID(),
		Status:  model.TaskStatusWaiting,
		Payload: map[string]string{"key": "value"},
	})
	s.Require().NoError(err)
	s.NotNil(task)
	s.Equal("task A", task.Name)
	s.Equal(s.users[0].GetID(), task.UserID)
	s.Equal(model.TaskStatusWaiting, task.Status)
	s.Equal(map[string]string{"key": "value"}, task.Payload)
}

func (s *TaskServiceSuite) TestPatch_Name() {
	task, err := service.NewTaskService().Create(service.TaskCreateOptions{
		Name:   "task",
		UserID: s.users[0].GetID(),
		Status: model.TaskStatusWaiting,
	})
	s.Require().NoError(err)

	task, err = service.NewTaskService().Patch(task.ID, service.TaskPatchOptions{
		Fields: []string{service.TaskFieldName},
		Name:   helper.ToPtr("task (edit)"),
	})
	s.Require().NoError(err)
	s.Equal("task (edit)", task.Name)
}

func (s *TaskServiceSuite) TestPatch_Status() {
	task, err := service.NewTaskService().Create(service.TaskCreateOptions{
		Name:   "task",
		UserID: s.users[0].GetID(),
		Status: model.TaskStatusWaiting,
	})
	s.Require().NoError(err)

	task, err = service.NewTaskService().Patch(task.ID, service.TaskPatchOptions{
		Fields: []string{service.TaskFieldStatus},
		Status: helper.ToPtr(model.TaskStatusRunning),
	})
	s.Require().NoError(err)
	s.Equal(model.TaskStatusRunning, task.Status)
}

func (s *TaskServiceSuite) TestPatch_Error() {
	task, err := service.NewTaskService().Create(service.TaskCreateOptions{
		Name:   "task",
		UserID: s.users[0].GetID(),
		Status: model.TaskStatusWaiting,
	})
	s.Require().NoError(err)

	task, err = service.NewTaskService().Patch(task.ID, service.TaskPatchOptions{
		Fields: []string{service.TaskFieldError},
		Error:  helper.ToPtr("something went wrong"),
	})
	s.Require().NoError(err)
	s.Equal("something went wrong", *task.Error)
}

func (s *TaskServiceSuite) TestPatch_Percentage() {
	task, err := service.NewTaskService().Create(service.TaskCreateOptions{
		Name:   "task",
		UserID: s.users[0].GetID(),
		Status: model.TaskStatusWaiting,
	})
	s.Require().NoError(err)

	task, err = service.NewTaskService().Patch(task.ID, service.TaskPatchOptions{
		Fields:     []string{service.TaskFieldPercentage},
		Percentage: helper.ToPtr(50),
	})
	s.Require().NoError(err)
	s.Equal(50, *task.Percentage)
}

func (s *TaskServiceSuite) TestPatch_Payload() {
	task, err := service.NewTaskService().Create(service.TaskCreateOptions{
		Name:   "task",
		UserID: s.users[0].GetID(),
		Status: model.TaskStatusWaiting,
	})
	s.Require().NoError(err)

	task, err = service.NewTaskService().Patch(task.ID, service.TaskPatchOptions{
		Fields:  []string{service.TaskFieldPayload},
		Payload: map[string]string{"key": "value"},
	})
	s.Require().NoError(err)
	s.Equal(map[string]string{"key": "value"}, task.Payload)
}

func (s *TaskServiceSuite) TestFind() {
	task, err := service.NewTaskService().Create(service.TaskCreateOptions{
		Name:   "task",
		UserID: s.users[0].GetID(),
		Status: model.TaskStatusWaiting,
	})
	s.Require().NoError(err)

	found, err := service.NewTaskService().Find(task.ID, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(task.ID, found.ID)
	s.Equal("task", found.Name)
}

func (s *TaskServiceSuite) TestFind_UnauthorizedUser() {
	task, err := service.NewTaskService().Create(service.TaskCreateOptions{
		Name:   "task",
		UserID: s.users[0].GetID(),
		Status: model.TaskStatusWaiting,
	})
	s.Require().NoError(err)

	_, err = service.NewTaskService().Find(task.ID, s.users[1].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewTaskNotFoundError(nil).Error(), err.Error())
}

func (s *TaskServiceSuite) TestList() {
	statuses := []string{model.TaskStatusWaiting, model.TaskStatusRunning, model.TaskStatusWaiting}
	for i, name := range []string{"task A", "task B", "task C"} {
		_, err := service.NewTaskService().Create(service.TaskCreateOptions{
			Name:   name,
			UserID: s.users[0].GetID(),
			Status: statuses[i],
		})
		s.Require().NoError(err)
	}

	list, err := service.NewTaskService().List(service.TaskListOptions{
		Page: 1,
		Size: 10,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(1), list.Page)
	s.Equal(uint64(3), list.Size)
	s.Equal(uint64(3), list.TotalElements)
	s.Equal(uint64(1), list.TotalPages)
	s.Equal("task A", list.Data[0].Name)
	s.Equal("task B", list.Data[1].Name)
	s.Equal("task C", list.Data[2].Name)
}

func (s *TaskServiceSuite) TestList_Paginate() {
	statuses := []string{model.TaskStatusWaiting, model.TaskStatusRunning, model.TaskStatusWaiting}
	for i, name := range []string{"task A", "task B", "task C"} {
		_, err := service.NewTaskService().Create(service.TaskCreateOptions{
			Name:   name,
			UserID: s.users[0].GetID(),
			Status: statuses[i],
		})
		s.Require().NoError(err)
	}

	list, err := service.NewTaskService().List(service.TaskListOptions{
		Page: 1,
		Size: 2,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(1), list.Page)
	s.Equal(uint64(2), list.Size)
	s.Equal(uint64(3), list.TotalElements)
	s.Equal(uint64(2), list.TotalPages)
	s.Equal("task A", list.Data[0].Name)
	s.Equal("task B", list.Data[1].Name)

	list, err = service.NewTaskService().List(service.TaskListOptions{
		Page: 2,
		Size: 2,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(2), list.Page)
	s.Equal(uint64(1), list.Size)
	s.Equal(uint64(3), list.TotalElements)
	s.Equal(uint64(2), list.TotalPages)
	s.Equal("task C", list.Data[0].Name)
}

func (s *TaskServiceSuite) TestList_SortByStatusDescending() {
	statuses := []string{model.TaskStatusWaiting, model.TaskStatusRunning, model.TaskStatusWaiting}
	for i, name := range []string{"task A", "task B", "task C"} {
		_, err := service.NewTaskService().Create(service.TaskCreateOptions{
			Name:   name,
			UserID: s.users[0].GetID(),
			Status: statuses[i],
		})
		s.Require().NoError(err)
	}

	list, err := service.NewTaskService().List(service.TaskListOptions{
		Page:      1,
		Size:      3,
		SortBy:    service.TaskSortByStatus,
		SortOrder: service.TaskSortOrderDesc,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal("task B", list.Data[0].Name)
	s.Equal("task A", list.Data[1].Name)
	s.Equal("task C", list.Data[2].Name)
}

func (s *TaskServiceSuite) TestList_Query() {
	statuses := []string{model.TaskStatusWaiting, model.TaskStatusRunning, model.TaskStatusWaiting}
	for i, name := range []string{"foo bar", "hello world", "lorem ipsum"} {
		_, err := service.NewTaskService().Create(service.TaskCreateOptions{
			Name:   name,
			UserID: s.users[0].GetID(),
			Status: statuses[i],
		})
		s.Require().NoError(err)
	}

	list, err := service.NewTaskService().List(service.TaskListOptions{
		Query: "world",
		Page:  1,
		Size:  10,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(1), list.Page)
	s.Equal(uint64(1), list.Size)
	s.Equal(uint64(1), list.TotalElements)
	s.Equal(uint64(1), list.TotalPages)
	s.Equal("hello world", list.Data[0].Name)
}

func (s *TaskServiceSuite) TestProbe() {
	statuses := []string{model.TaskStatusWaiting, model.TaskStatusRunning, model.TaskStatusWaiting}
	for i, name := range []string{"task A", "task B", "task C"} {
		_, err := service.NewTaskService().Create(service.TaskCreateOptions{
			Name:   name,
			UserID: s.users[0].GetID(),
			Status: statuses[i],
		})
		s.Require().NoError(err)
	}

	probe, err := service.NewTaskService().Probe(service.TaskListOptions{
		Page: 1,
		Size: 10,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(3), probe.TotalElements)
	s.Equal(uint64(1), probe.TotalPages)
}

func (s *TaskServiceSuite) TestDismiss() {
	task, err := service.NewTaskService().Create(service.TaskCreateOptions{
		Name:   "task with error",
		UserID: s.users[0].GetID(),
		Status: model.TaskStatusError,
		Error:  helper.ToPtr("something went wrong"),
	})
	s.Require().NoError(err)

	err = service.NewTaskService().Dismiss(task.ID, s.users[0].GetID())
	s.Require().NoError(err)
}

func (s *TaskServiceSuite) TestDismiss_UnauthorizedUser() {
	task, err := service.NewTaskService().Create(service.TaskCreateOptions{
		Name:   "task with error",
		UserID: s.users[0].GetID(),
		Status: model.TaskStatusError,
		Error:  helper.ToPtr("something went wrong"),
	})
	s.Require().NoError(err)

	err = service.NewTaskService().Dismiss(task.ID, s.users[1].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewTaskBelongsToAnotherUserError(nil).Error(), err.Error())
}

func (s *TaskServiceSuite) TestDismiss_StatusRunning() {
	task, err := service.NewTaskService().Create(service.TaskCreateOptions{
		Name:   "task with status running",
		UserID: s.users[0].GetID(),
		Status: model.TaskStatusRunning,
	})
	s.Require().NoError(err)

	err = service.NewTaskService().Dismiss(task.ID, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewTaskIsRunningError(nil).Error(), err.Error())
}

func (s *TaskServiceSuite) TestDismissAll() {
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

	dismissAllResult, err := service.NewTaskService().DismissAll(s.users[0].GetID())
	s.Require().NoError(err)
	s.Len(dismissAllResult.Succeeded, 2)
	s.Empty(dismissAllResult.Failed)

	list, err := service.NewTaskService().List(service.TaskListOptions{
		Page: 1,
		Size: 10,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Empty(list.Data)
}

func (s *TaskServiceSuite) TestDelete() {
	task, err := service.NewTaskService().Create(service.TaskCreateOptions{
		Name:   "task",
		UserID: s.users[0].GetID(),
		Status: model.TaskStatusWaiting,
	})
	s.Require().NoError(err)

	err = service.NewTaskService().Delete(task.ID)
	s.Require().NoError(err)

	_, err = service.NewTaskService().Find(task.ID, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewTaskNotFoundError(nil).Error(), err.Error())
}

func (s *TaskServiceSuite) TestCount() {
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

	count, err := service.NewTaskService().Count(s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(int64(2), *count)
}
