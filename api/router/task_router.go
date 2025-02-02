// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package router

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/kouprlabs/voltaserve/api/config"
	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/service"
)

type TaskRouter struct {
	taskSvc service.TaskService
	config  *config.Config
}

func NewTaskRouter() *TaskRouter {
	return &TaskRouter{
		taskSvc: service.NewTaskService(),
		config:  config.GetConfig(),
	}
}

func (r *TaskRouter) AppendRoutes(g fiber.Router) {
	g.Get("/", r.List)
	g.Get("/probe", r.Probe)
	g.Get("/count", r.Count)
	g.Post("/dismiss", r.DismissAll)
	g.Get("/:id", r.Get)
	g.Post("/:id/dismiss", r.Dismiss)
}

func (r *TaskRouter) AppendNonJWTRoutes(g fiber.Router) {
	g.Post("/", r.Create)
	g.Delete("/:id", r.Delete)
	g.Patch("/:id", r.Patch)
}

// Get godoc
//
//	@Summary		Read
//	@Description	Read
//	@Tags			Tasks
//	@Id				tasks_get
//	@Produce		json
//	@Param			id	path		string	true	"ID"
//	@Success		200	{object}	service.Task
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/tasks/{id} [get]
func (r *TaskRouter) Get(c *fiber.Ctx) error {
	userID := GetUserID(c)
	res, err := r.taskSvc.Find(c.Params("id"), userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// List godoc
//
//	@Summary		List
//	@Description	List
//	@Tags			Tasks
//	@Id				tasks_list
//	@Produce		json
//	@Param			query		query		string	false	"Query"
//	@Param			page		query		string	false	"Page"
//	@Param			size		query		string	false	"Size"
//	@Param			sort_by		query		string	false	"Sort By"
//	@Param			sort_order	query		string	false	"Sort Order"
//	@Success		200			{object}	service.TaskList
//	@Failure		404			{object}	errorpkg.ErrorResponse
//	@Failure		500			{object}	errorpkg.ErrorResponse
//	@Router			/tasks [get]
func (r *TaskRouter) List(c *fiber.Ctx) error {
	opts, err := r.parseListQueryParams(c)
	if err != nil {
		return err
	}
	res, err := r.taskSvc.List(*opts, GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// Probe godoc
//
//	@Summary		Probe
//	@Description	Probe
//	@Tags			Tasks
//	@Id				tasks_probe
//	@Produce		json
//	@Param			size	query		string	false	"Size"
//	@Success		200		{object}	service.TaskProbe
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/tasks/probe [get]
func (r *TaskRouter) Probe(c *fiber.Ctx) error {
	opts, err := r.parseListQueryParams(c)
	if err != nil {
		return err
	}
	res, err := r.taskSvc.Probe(*opts, GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

func (r *TaskRouter) parseListQueryParams(c *fiber.Ctx) (*service.TaskListOptions, error) {
	var err error
	var page uint64
	if c.Query("page") == "" {
		page = 1
	} else {
		page, err = strconv.ParseUint(c.Query("page"), 10, 64)
		if err != nil {
			return nil, errorpkg.NewInvalidQueryParamError("page")
		}
	}
	var size uint64
	if c.Query("size") == "" {
		size = OrganizationDefaultPageSize
	} else {
		size, err = strconv.ParseUint(c.Query("size"), 10, 64)
		if err != nil {
			return nil, errorpkg.NewInvalidQueryParamError("size")
		}
	}
	if size == 0 {
		return nil, errorpkg.NewInvalidQueryParamError("size")
	}
	sortBy := c.Query("sort_by")
	if !IsValidSortBy(sortBy) {
		return nil, errorpkg.NewInvalidQueryParamError("sort_by")
	}
	sortOrder := c.Query("sort_order")
	if !IsValidSortOrder(sortOrder) {
		return nil, errorpkg.NewInvalidQueryParamError("sort_order")
	}
	query, err := url.QueryUnescape(c.Query("query"))
	if err != nil {
		return nil, errorpkg.NewInvalidQueryParamError("query")
	}
	return &service.TaskListOptions{
		Query:     query,
		Page:      page,
		Size:      size,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}, nil
}

// Count godoc
//
//	@Summary		Count
//	@Description	Count
//	@Tags			Tasks
//	@Id				tasks_count
//	@Produce		json
//	@Success		200	{integer}	int
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/tasks/count [get]
func (r *TaskRouter) Count(c *fiber.Ctx) error {
	res, err := r.taskSvc.Count(GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// Dismiss godoc
//
//	@Summary		Dismiss
//	@Description	Dismiss
//	@Tags			Tasks
//	@Id				tasks_dismiss
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"ID"
//	@Success		200
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/tasks/{id}/dismiss [post]
func (r *TaskRouter) Dismiss(c *fiber.Ctx) error {
	userID := GetUserID(c)
	if err := r.taskSvc.Dismiss(c.Params("id"), userID); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}

// DismissAll godoc
//
//	@Summary		Dismiss All
//	@Description	Dismiss All
//	@Tags			Tasks
//	@Id				tasks_dismiss_all
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	service.TaskDismissAllResult
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/tasks/dismiss [post]
func (r *TaskRouter) DismissAll(c *fiber.Ctx) error {
	userID := GetUserID(c)
	res, err := r.taskSvc.DismissAll(userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// Create godoc
//
//	@Summary		Create
//	@Description	Create
//	@Tags			Tasks
//	@Id				tasks_create
//	@Produce		json
//	@Param			api_key	query	string						true	"API Key"
//	@Param			id		path	string						true	"ID"
//	@Param			body	body	service.TaskCreateOptions	true	"Body"
//	@Success		204
//	@Failure		401	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/tasks [post]
func (r *TaskRouter) Create(c *fiber.Ctx) error {
	apiKey := c.Query("api_key")
	if apiKey == "" {
		return errorpkg.NewMissingQueryParamError("api_key")
	}
	if apiKey != r.config.Security.APIKey {
		return errorpkg.NewInvalidAPIKeyError()
	}
	opts := new(service.TaskCreateOptions)
	if err := c.BodyParser(opts); err != nil {
		return err
	}
	if err := validator.New().Struct(opts); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	task, err := r.taskSvc.Create(*opts)
	if err != nil {
		return err
	}
	return c.JSON(task)
}

// Delete godoc
//
//	@Summary		Delete
//	@Description	Delete
//	@Tags			Tasks
//	@Id				tasks_delete
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"ID"
//	@Success		200
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/tasks/{id} [delete]
func (r *TaskRouter) Delete(c *fiber.Ctx) error {
	apiKey := c.Query("api_key")
	if apiKey == "" {
		return errorpkg.NewMissingQueryParamError("api_key")
	}
	if apiKey != r.config.Security.APIKey {
		return errorpkg.NewInvalidAPIKeyError()
	}
	if err := r.taskSvc.Delete(c.Params("id")); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}

// Patch godoc
//
//	@Summary		Patch
//	@Description	Patch
//	@Tags			Tasks
//	@Id				tasks_patch
//	@Produce		json
//	@Param			api_key	query		string						true	"API Key"
//	@Param			id		path		string						true	"ID"
//	@Param			body	body		service.TaskPatchOptions	true	"Body"
//	@Success		200		{object}	service.Task
//	@Failure		401		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/tasks/{id} [patch]
func (r *TaskRouter) Patch(c *fiber.Ctx) error {
	apiKey := c.Query("api_key")
	if apiKey == "" {
		return errorpkg.NewMissingQueryParamError("api_key")
	}
	if apiKey != r.config.Security.APIKey {
		return errorpkg.NewInvalidAPIKeyError()
	}
	opts := new(service.TaskPatchOptions)
	if err := c.BodyParser(opts); err != nil {
		return err
	}
	if err := validator.New().Struct(opts); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	task, err := r.taskSvc.Patch(c.Params("id"), *opts)
	if err != nil {
		return err
	}
	return c.JSON(task)
}
