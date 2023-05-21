package router

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"voltaserve/config"
	"voltaserve/core"
	"voltaserve/errorpkg"
	"voltaserve/helpers"
	"voltaserve/model"
	"voltaserve/storage"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type FileRouter struct {
	fileSvc      *core.FileService
	workspaceSvc *core.WorkspaceService
	storageSvc   *storage.StorageService
	config       config.Config
}

func NewFileRouter() *FileRouter {
	return &FileRouter{
		fileSvc:      core.NewFileService(),
		workspaceSvc: core.NewWorkspaceService(),
		storageSvc:   storage.NewStorageService(),
		config:       config.GetConfig(),
	}
}

func (r *FileRouter) AppendRoutes(g fiber.Router) {
	g.Post("/", r.Upload)
	g.Post("/create_folder", r.CreateFolder)
	g.Get("/list", r.ListByPath)
	g.Post("/search", r.Search)
	g.Post("/batch_delete", r.BatchDelete)
	g.Post("/batch_get", r.BatchGet)
	g.Get("/:id", r.GetById)
	g.Patch("/:id", r.Patch)
	g.Delete("/:id", r.Delete)
	g.Get("/:id/list", r.ListByID)
	g.Get("/:id/get_item_count", r.GetItemCount)
	g.Get("/:id/get_path", r.GetPath)
	g.Post("/:id/move", r.Move)
	g.Post("/:id/rename", r.Rename)
	g.Post("/:id/copy", r.Copy)
	g.Get("/:id/get_size", r.GetSize)
	g.Post("/grant_user_permission", r.GrantUserPermission)
	g.Post("/revoke_user_permission", r.RevokeUserPermission)
	g.Post("/grant_group_permission", r.GrantGroupPermission)
	g.Post("/revoke_group_permission", r.RevokeGroupPermission)
	g.Get("/:id/get_user_permissions", r.GetUserPermissions)
	g.Get("/:id/get_group_permissions", r.GetGroupPermissions)
}

// Upload godoc
// @Summary     Upload
// @Description Upload
// @Tags        Files
// @Id          files_upload
// @Accept      x-www-form-urlencoded
// @Produce     json
// @Param       workspace_id query    string true  "Workspace Id"
// @Param       parent_id    query    string false "Parent Id"
// @Success     200          {object} core.File
// @Failure     404          {object} errorpkg.ErrorResponse
// @Failure     400          {object} errorpkg.ErrorResponse
// @Failure     500          {object} errorpkg.ErrorResponse
// @Router      /files [post]
func (r *FileRouter) Upload(c *fiber.Ctx) error {
	userId := GetUserId(c)
	workspaceId := c.Query("workspace_id")
	if workspaceId == "" {
		return errorpkg.NewMissingQueryParamError("workspace_id")
	}
	parentId := c.Query("parent_id")
	if parentId == "" {
		workspace, err := r.workspaceSvc.Find(workspaceId, userId)
		if err != nil {
			return err
		}
		parentId = workspace.RootId
	}
	fh, err := c.FormFile("file")
	if err != nil {
		return err
	}
	ok, err := r.workspaceSvc.HasEnoughSpaceForByteSize(workspaceId, fh.Size)
	if err != nil {
		return err
	}
	if !ok {
		return errorpkg.NewStorageLimitExceededError()
	}
	file, err := r.fileSvc.Create(core.FileCreateOptions{
		Name:        fh.Filename,
		Type:        model.FileTypeFile,
		ParentId:    &parentId,
		WorkspaceId: workspaceId,
	}, userId)
	if err != nil {
		return err
	}
	path := filepath.FromSlash(os.TempDir() + "/" + helpers.NewId() + filepath.Ext(fh.Filename))
	if err := c.SaveFile(fh, path); err != nil {
		return err
	}
	defer os.Remove(path)
	file, err = r.storageSvc.Store(storage.StorageOptions{FileId: file.Id, FilePath: path}, userId)
	if err != nil {
		return err
	}
	return c.Status(http.StatusCreated).JSON(file)
}

// Patch godoc
// @Summary     Patch
// @Description Patch
// @Tags        Files
// @Id          files_patch
// @Accept      x-www-form-urlencoded
// @Produce     json
// @Param       id  path     string true "Id"
// @Success     200 {object} core.File
// @Failure     404 {object} errorpkg.ErrorResponse
// @Failure     400 {object} errorpkg.ErrorResponse
// @Failure     500 {object} errorpkg.ErrorResponse
// @Router      /files/{id} [patch]
func (r *FileRouter) Patch(c *fiber.Ctx) error {
	userId := GetUserId(c)
	files, err := r.fileSvc.Find([]string{c.Params("id")}, userId)
	if err != nil {
		return err
	}
	file := files[0]
	fh, err := c.FormFile("file")
	if err != nil {
		return err
	}
	ok, err := r.workspaceSvc.HasEnoughSpaceForByteSize(file.WorkspaceId, fh.Size)
	if err != nil {
		return err
	}
	if !ok {
		return errorpkg.NewStorageLimitExceededError()
	}
	path := filepath.FromSlash(os.TempDir() + "/" + helpers.NewId() + filepath.Ext(fh.Filename))
	if err := c.SaveFile(fh, path); err != nil {
		return err
	}
	defer os.Remove(path)
	file, err = r.storageSvc.Store(storage.StorageOptions{FileId: file.Id, FilePath: path}, userId)
	if err != nil {
		return err
	}
	return c.JSON(file)
}

// Create godoc
// @Summary     Create
// @Description Create
// @Tags        Files
// @Id          files_create_folder
// @Accept      json
// @Produce     json
// @Param       body body     core.FileCreateFolderOptions true "Body"
// @Success     200  {object} core.File
// @Failure     400  {object} errorpkg.ErrorResponse
// @Failure     500  {object} errorpkg.ErrorResponse
// @Router      /files/create_folder [post]
func (r *FileRouter) CreateFolder(c *fiber.Ctx) error {
	userId := GetUserId(c)
	req := new(core.FileCreateFolderOptions)
	if err := c.BodyParser(req); err != nil {
		return err
	}
	if err := validator.New().Struct(req); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	parentId := req.ParentId
	if parentId == nil {
		workspace, err := r.workspaceSvc.Find(req.WorkspaceId, userId)
		if err != nil {
			return err
		}
		parentId = &workspace.RootId
	}
	res, err := r.fileSvc.Create(core.FileCreateOptions{
		Name:        req.Name,
		Type:        model.FileTypeFolder,
		ParentId:    parentId,
		WorkspaceId: req.WorkspaceId,
	}, userId)
	if err != nil {
		return err
	}
	return c.Status(http.StatusCreated).JSON(res)
}

// Search godoc
// @Summary     Search
// @Description Search
// @Tags        Files
// @Id          files_search
// @Produce     json
// @Param       page query    string                 true "Page"
// @Param       size query    string                 true "Size"
// @Param       body body     core.FileSearchOptions true "Body"
// @Success     200  {object} core.FileSearchResult
// @Failure     500  {object} errorpkg.ErrorResponse
// @Router      /files/search [post]
func (r *FileRouter) Search(c *fiber.Ctx) error {
	userId := GetUserId(c)
	req := new(core.FileSearchOptions)
	if err := c.BodyParser(req); err != nil {
		return err
	}
	if err := validator.New().Struct(req); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	page, err := strconv.ParseUint(c.Params("page"), 10, 32)
	if err != nil {
		page = 1
	}
	size, err := strconv.ParseUint(c.Params("size"), 10, 32)
	if err != nil {
		size = 100
	}
	res, err := r.fileSvc.Search(*req, uint(page), uint(size), userId)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// GetById godoc
// @Summary     Get by Id
// @Description Get by Id
// @Tags        Files
// @Id          files_get_by_id
// @Produce     json
// @Param       id  path     string true "Id"
// @Success     200 {object} core.File
// @Failure     404 {object} errorpkg.ErrorResponse
// @Failure     500 {object} errorpkg.ErrorResponse
// @Router      /files/{id} [get]
func (r *FileRouter) GetById(c *fiber.Ctx) error {
	userId := GetUserId(c)
	res, err := r.fileSvc.Find([]string{c.Params("id")}, userId)
	if err != nil {
		return err
	}
	return c.JSON(res[0])
}

// ListByPath godoc
// @Summary     ListByPath
// @Description ListByPath
// @Tags        Files
// @Id          files_list_by_path
// @Produce     json
// @Param       path query    string true "Path"
// @Success     200  {array}  core.File
// @Failure     404  {object} errorpkg.ErrorResponse
// @Failure     500  {object} errorpkg.ErrorResponse
// @Router      /files/list [get]
func (r *FileRouter) ListByPath(c *fiber.Ctx) error {
	userId := GetUserId(c)
	if c.Query("path") == "" {
		return errorpkg.NewMissingQueryParamError("path")
	}
	res, err := r.fileSvc.ListByPath(c.Query("path"), userId)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// ListByID godoc
// @Summary     ListByID
// @Description ListByID
// @Tags        Files
// @Id          files_list_by_id
// @Produce     json
// @Param       id   path     string true  "Id"
// @Param       page query    string true  "Page"
// @Param       size query    string true  "Size"
// @Param       type query    string false "Type"
// @Success     200  {object} core.FileList
// @Failure     404  {object} errorpkg.ErrorResponse
// @Failure     500  {object} errorpkg.ErrorResponse
// @Router      /files/{id}/list [get]
func (r *FileRouter) ListByID(c *fiber.Ctx) error {
	if c.Query("page") == "" {
		return errorpkg.NewMissingQueryParamError("page")
	}
	if c.Query("size") == "" {
		return errorpkg.NewMissingQueryParamError("size")
	}
	fileType := c.Query("type")
	if fileType != model.FileTypeFile && fileType != model.FileTypeFolder && fileType != "" {
		return errorpkg.NewInvalidQueryParamError("type")
	}
	userId := GetUserId(c)
	page, err := strconv.ParseInt(c.Query("page"), 10, 32)
	if err != nil {
		page = 1
	}
	size, err := strconv.ParseInt(c.Query("size"), 10, 32)
	if err != nil {
		size = 100
	}
	res, err := r.fileSvc.ListByID(c.Params("id"), uint(page), uint(size), fileType, userId)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// GetPath godoc
// @Summary     Get path
// @Description Get path
// @Tags        Files
// @Id          files_get_path
// @Produce     json
// @Param       id  path     string true "Id"
// @Success     200 {array}  core.File
// @Failure     404 {object} errorpkg.ErrorResponse
// @Failure     500 {object} errorpkg.ErrorResponse
// @Router      /files/{id}/get_path [get]
func (r *FileRouter) GetPath(c *fiber.Ctx) error {
	userId := GetUserId(c)
	res, err := r.fileSvc.GetPath(c.Params("id"), userId)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// Copy godoc
// @Summary     Copy
// @Description Copy
// @Tags        Files
// @Id          files_copy
// @Produce     json
// @Param       id   path     string               true "Id"
// @Param       body body     core.FileCopyOptions true "Body"
// @Failure     404  {object} errorpkg.ErrorResponse
// @Failure     500  {object} errorpkg.ErrorResponse
// @Router      /files/{id}/copy [post]
func (r *FileRouter) Copy(c *fiber.Ctx) error {
	userId := GetUserId(c)
	req := new(core.FileCopyOptions)
	if err := c.BodyParser(req); err != nil {
		return err
	}
	if err := validator.New().Struct(req); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	if err := r.fileSvc.Copy(c.Params("id"), req.Ids, userId); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}

// Move godoc
// @Summary     Move
// @Description Move
// @Tags        Files
// @Id          files_move
// @Produce     json
// @Param       id   path     string               true "Id"
// @Param       body body     core.FileMoveOptions true "Body"
// @Failure     404  {object} errorpkg.ErrorResponse
// @Failure     500  {object} errorpkg.ErrorResponse
// @Router      /files/{id}/move [post]
func (r *FileRouter) Move(c *fiber.Ctx) error {
	userId := GetUserId(c)
	req := new(core.FileMoveOptions)
	if err := c.BodyParser(req); err != nil {
		return err
	}
	if err := validator.New().Struct(req); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	if _, err := r.fileSvc.Move(c.Params("id"), req.Ids, userId); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}

// Rename godoc
// @Summary     Rename
// @Description Rename
// @Tags        Files
// @Id          files_rename
// @Produce     json
// @Param       id   path     string                 true "Id"
// @Param       body body     core.FileRenameOptions true "Body"
// @Success     200  {object} core.File
// @Failure     404  {object} errorpkg.ErrorResponse
// @Failure     500  {object} errorpkg.ErrorResponse
// @Router      /files/{id}/rename [post]
func (r *FileRouter) Rename(c *fiber.Ctx) error {
	userId := GetUserId(c)
	req := new(core.FileRenameOptions)
	if err := c.BodyParser(req); err != nil {
		return err
	}
	if err := validator.New().Struct(req); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	res, err := r.fileSvc.Rename(c.Params("id"), req.Name, userId)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// Delete godoc
// @Summary     Delete
// @Description Delete
// @Tags        Files
// @Id          files_delete
// @Produce     json
// @Param       id  path     string true "Id"
// @Failure     404 {object} errorpkg.ErrorResponse
// @Failure     500 {object} errorpkg.ErrorResponse
// @Router      /files/{id} [delete]
func (r *FileRouter) Delete(c *fiber.Ctx) error {
	userId := GetUserId(c)
	_, err := r.fileSvc.Delete([]string{c.Params("id")}, userId)
	if err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}

// BatchGet godoc
// @Summary     Batch get
// @Description Batch get
// @Tags        Files
// @Id          files_batch_get
// @Produce     json
// @Param       body body     core.FileBatchGetOptions true "Body"
// @Success     200  {array}  core.File
// @Failure     500  {object} errorpkg.ErrorResponse
// @Router      /files/batch_get [post]
func (r *FileRouter) BatchGet(c *fiber.Ctx) error {
	userId := GetUserId(c)
	req := new(core.FileBatchGetOptions)
	if err := c.BodyParser(req); err != nil {
		return err
	}
	if err := validator.New().Struct(req); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	res, err := r.fileSvc.Find(req.Ids, userId)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// BatchDelete godoc
// @Summary     Batch delete
// @Description Batch delete
// @Tags        Files
// @Id          files_batch_delete
// @Produce     json
// @Param       body body     core.FileBatchDeleteOptions true "Body"
// @Success     200  {array}  string
// @Failure     500  {object} errorpkg.ErrorResponse
// @Router      /files/batch_delete [post]
func (r *FileRouter) BatchDelete(c *fiber.Ctx) error {
	userId := GetUserId(c)
	req := new(core.FileBatchDeleteOptions)
	if err := c.BodyParser(req); err != nil {
		return err
	}
	if err := validator.New().Struct(req); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	res, err := r.fileSvc.Delete(req.Ids, userId)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// GetSize godoc
// @Summary     Get size
// @Description Get size
// @Tags        Files
// @Id          files_get_size
// @Produce     json
// @Param       id  path     string true "Id"
// @Success     200 {object} int
// @Failure     404 {object} errorpkg.ErrorResponse
// @Failure     500 {object} errorpkg.ErrorResponse
// @Router      /files/{id}/get_size [get]
func (r *FileRouter) GetSize(c *fiber.Ctx) error {
	userId := GetUserId(c)
	id := c.Params("id")
	res, err := r.fileSvc.GetSize(id, userId)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// GetItemCount godoc
// @Summary     Get children count
// @Description Get children count
// @Tags        Files
// @Id          files_get_children_count
// @Produce     json
// @Param       id  path     string true "Id"
// @Success     200 {object} int
// @Failure     404 {object} errorpkg.ErrorResponse
// @Failure     500 {object} errorpkg.ErrorResponse
// @Router      /files/{id}/get_item_count [get]
func (r *FileRouter) GetItemCount(c *fiber.Ctx) error {
	userId := GetUserId(c)
	res, err := r.fileSvc.GetItemCount(c.Params("id"), userId)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// GrantUserPermission godoc
// @Summary     Grant user permission
// @Description Grant user permission
// @Tags        Files
// @Id          files_grant_user_permission
// @Produce     json
// @Param       id   path     string                              true "Id"
// @Param       body body     core.FileGrantUserPermissionOptions true "Body"
// @Failure     404  {object} errorpkg.ErrorResponse
// @Failure     500  {object} errorpkg.ErrorResponse
// @Router      /files/grant_user_permission [post]
func (r *FileRouter) GrantUserPermission(c *fiber.Ctx) error {
	userId := GetUserId(c)
	req := new(core.FileGrantUserPermissionOptions)
	if err := c.BodyParser(req); err != nil {
		return err
	}
	if err := validator.New().Struct(req); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	if err := r.fileSvc.GrantUserPermission(req.Ids, req.UserId, req.Permission, userId); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}

// RevokeUserPermission godoc
// @Summary     Revoke user permission
// @Description Revoke user permission
// @Tags        Files
// @Id          files_revoke_user_permission
// @Produce     json
// @Param       id   path     string                               true "Id"
// @Param       body body     core.FileRevokeUserPermissionOptions true "Body"
// @Failure     404  {object} errorpkg.ErrorResponse
// @Failure     500  {object} errorpkg.ErrorResponse
// @Router      /files/revoke_user_permission [post]
func (r *FileRouter) RevokeUserPermission(c *fiber.Ctx) error {
	userId := GetUserId(c)
	req := new(core.FileRevokeUserPermissionOptions)
	if err := c.BodyParser(req); err != nil {
		return err
	}
	if err := validator.New().Struct(req); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	if err := r.fileSvc.RevokeUserPermission(req.Ids, req.UserId, userId); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}

// GrantGroupPermission godoc
// @Summary     Grant group permission
// @Description Grant group permission
// @Tags        Files
// @Id          files_grant_group_permission
// @Produce     json
// @Param       id   path     string                               true "Id"
// @Param       body body     core.FileGrantGroupPermissionOptions true "Body"
// @Failure     404  {object} errorpkg.ErrorResponse
// @Failure     500  {object} errorpkg.ErrorResponse
// @Router      /files/grant_group_permission [post]
func (r *FileRouter) GrantGroupPermission(c *fiber.Ctx) error {
	userId := GetUserId(c)
	req := new(core.FileGrantGroupPermissionOptions)
	if err := c.BodyParser(req); err != nil {
		return err
	}
	if err := validator.New().Struct(req); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	if err := r.fileSvc.GrantGroupPermission(req.Ids, req.GroupId, req.Permission, userId); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}

// RevokeGroupPermission godoc
// @Summary     Revoke group permission
// @Description Revoke group permission
// @Tags        Files
// @Id          files_revoke_group_permission
// @Produce     json
// @Param       id   path     string                                true "Id"
// @Param       body body     core.FileRevokeGroupPermissionOptions true "Body"
// @Failure     404  {object} errorpkg.ErrorResponse
// @Failure     500  {object} errorpkg.ErrorResponse
// @Router      /files/{id}/revoke_group_permission [post]
func (r *FileRouter) RevokeGroupPermission(c *fiber.Ctx) error {
	userId := GetUserId(c)
	req := new(core.FileRevokeGroupPermissionOptions)
	if err := c.BodyParser(req); err != nil {
		return err
	}
	if err := validator.New().Struct(req); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	if err := r.fileSvc.RevokeGroupPermission(req.Ids, req.GroupId, userId); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}

// GetUserPermissions godoc
// @Summary     Get user permissions
// @Description Get user permissions
// @Tags        Files
// @Id          files_get_user_permissions
// @Produce     json
// @Param       id  path     string true "Id"
// @Success     200 {array}  core.UserPermission
// @Failure     404 {object} errorpkg.ErrorResponse
// @Failure     500 {object} errorpkg.ErrorResponse
// @Router      /files/{id}/get_user_permissions [get]
func (r *FileRouter) GetUserPermissions(c *fiber.Ctx) error {
	userId := GetUserId(c)
	res, err := r.fileSvc.GetUserPermissions(c.Params("id"), userId)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// GetGroupPermissions godoc
// @Summary     Get group permissions
// @Description Get group permissions
// @Tags        Files
// @Id          files_get_group_permissions
// @Produce     json
// @Param       id  path     string true "Id"
// @Success     200 {array}  core.GroupPermission
// @Failure     404 {object} errorpkg.ErrorResponse
// @Failure     500 {object} errorpkg.ErrorResponse
// @Router      /files/{id}/get_group_permissions [get]
func (r *FileRouter) GetGroupPermissions(c *fiber.Ctx) error {
	userId := GetUserId(c)
	res, err := r.fileSvc.GetGroupPermissions(c.Params("id"), userId)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

type FileDownloadRouter struct {
	fileSvc               *core.FileService
	accessTokenCookieName string
}

func NewFileDownloadRouter() *FileDownloadRouter {
	return &FileDownloadRouter{
		fileSvc:               core.NewFileService(),
		accessTokenCookieName: "voltaserve_access_token",
	}
}

func (r *FileDownloadRouter) AppendRoutes(g fiber.Router) {
	g.Get("/:id/original:ext", r.DownloadOriginal)
	g.Get("/:id/preview:ext", r.DownloadPreview)
}

// DownloadOriginal godoc
// @Summary     Download original
// @Description Download original
// @Tags        Files
// @Id          files_download_original
// @Produce     json
// @Param       id  path     string true "Id"
// @Failure     404 {object} errorpkg.ErrorResponse
// @Failure     500 {object} errorpkg.ErrorResponse
// @Router      /files/{id}/original{ext} [get]
func (r *FileDownloadRouter) DownloadOriginal(c *fiber.Ctx) error {
	accessToken := c.Cookies(r.accessTokenCookieName)
	if accessToken == "" {
		return errorpkg.NewFileNotFoundError(nil)
	}
	userId, err := r.getUserId(accessToken)
	if err != nil {
		return c.SendStatus(http.StatusNotFound)
	}
	buf, file, snapshot, err := r.fileSvc.DownloadOriginalBuffer(c.Params("id"), userId)
	if err != nil {
		return err
	}
	if filepath.Ext(snapshot.GetOriginal().Key) != c.Params("ext") {
		return errorpkg.NewS3ObjectNotFoundError(nil)
	}
	bytes := buf.Bytes()
	c.Set("Content-Type", storage.DetectMimeFromBytes(bytes))
	c.Set("Content-Disposition", fmt.Sprintf("filename=\"%s\"", file.GetName()))
	return c.Send(bytes)
}

// DownloadPreview godoc
// @Summary     Download preview
// @Description Download preview
// @Tags        Files
// @Id          files_download_preview
// @Produce     json
// @Param       id  path     string true "Id"
// @Failure     404 {object} errorpkg.ErrorResponse
// @Failure     500 {object} errorpkg.ErrorResponse
// @Router      /files/{id}/preview{ext} [get]
func (r *FileDownloadRouter) DownloadPreview(c *fiber.Ctx) error {
	accessToken := c.Cookies(r.accessTokenCookieName)
	if accessToken == "" {
		return errorpkg.NewFileNotFoundError(nil)
	}
	userId, err := r.getUserId(accessToken)
	if err != nil {
		return c.SendStatus(http.StatusNotFound)
	}
	buf, file, snapshot, err := r.fileSvc.DownloadPreviewBuffer(c.Params("id"), userId)
	if err != nil {
		return err
	}
	if filepath.Ext(snapshot.GetPreview().Key) != c.Params("ext") {
		return errorpkg.NewS3ObjectNotFoundError(nil)
	}
	bytes := buf.Bytes()
	c.Set("Content-Type", storage.DetectMimeFromBytes(bytes))
	c.Set("Content-Disposition", fmt.Sprintf("filename=\"%s\"", file.GetName()))
	return c.Send(bytes)
}

func (r *FileDownloadRouter) getUserId(accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.GetConfig().Security.JWTSigningKey), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["sub"].(string), nil
	} else {
		return "", errors.New("cannot find sub claim")
	}
}
