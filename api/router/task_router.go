package router

import (
	"net/http"
	"net/url"
	"strconv"
	"voltaserve/errorpkg"
	"voltaserve/service"

	"github.com/gofiber/fiber/v2"
)

type TaskRouter struct {
	taskSvc *service.TaskService
}

func NewTaskRouter() *TaskRouter {
	return &TaskRouter{
		taskSvc: service.NewTaskService(),
	}
}

func (r *TaskRouter) AppendRoutes(g fiber.Router) {
	g.Get("/", r.List)
	g.Get("/count", r.GetCount)
	g.Get("/:id", r.Get)
	g.Delete("/:id", r.Delete)
}

// Get godoc
//
//	@Summary		Get
//	@Description	Get
//	@Tags			Task
//	@Id				task_get
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
//	@Tags			Task
//	@Id				task_list
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
	var err error
	var page int64
	if c.Query("page") == "" {
		page = 1
	} else {
		page, err = strconv.ParseInt(c.Query("page"), 10, 64)
		if err != nil {
			page = 1
		}
	}
	var size int64
	if c.Query("size") == "" {
		size = OrganizationDefaultPageSize
	} else {
		size, err = strconv.ParseInt(c.Query("size"), 10, 64)
		if err != nil {
			return err
		}
	}
	sortBy := c.Query("sort_by")
	if !IsValidSortBy(sortBy) {
		return errorpkg.NewInvalidQueryParamError("sort_by")
	}
	sortOrder := c.Query("sort_order")
	if !IsValidSortOrder(sortOrder) {
		return errorpkg.NewInvalidQueryParamError("sort_order")
	}
	query, err := url.QueryUnescape(c.Query("query"))
	if err != nil {
		return errorpkg.NewInvalidQueryParamError("query")
	}
	res, err := r.taskSvc.List(service.TaskListOptions{
		Query:     query,
		Page:      uint(page),
		Size:      uint(size),
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}, GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// GetCount godoc
//
//	@Summary		Get Count
//	@Description	Get Count
//	@Tags			Task
//	@Id				task_get_count
//	@Produce		json
//	@Success		200	{object}	int
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/tasks/count [get]
func (r *TaskRouter) GetCount(c *fiber.Ctx) error {
	res, err := r.taskSvc.GetCount(GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// Delete godoc
//
//	@Summary		Delete
//	@Description	Delete
//	@Tags			Task
//	@Id				task_delete
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"ID"
//	@Success		200
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/tasks/{id} [delete]
func (r *TaskRouter) Delete(c *fiber.Ctx) error {
	userID := GetUserID(c)
	if err := r.taskSvc.Delete(c.Params("id"), userID); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}