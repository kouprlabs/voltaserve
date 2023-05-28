package repo

import (
	"errors"
	"time"
	"voltaserve/errorpkg"
	"voltaserve/helper"
	"voltaserve/infra"
	"voltaserve/model"

	"gorm.io/gorm"
)

type WorkspaceInsertOptions struct {
	ID              string
	Name            string
	StorageCapacity int64
	Image           *string
	OrganizationId  string
	RootId          string
	Bucket          string
}

type WorkspaceRepo interface {
	Insert(opts WorkspaceInsertOptions) (model.Workspace, error)
	FindByName(name string) (model.Workspace, error)
	FindByID(id string) (model.Workspace, error)
	UpdateName(id string, name string) (model.Workspace, error)
	UpdateStorageCapacity(id string, storageCapacity int64) (model.Workspace, error)
	Delete(id string) error
	GetIDs() ([]string, error)
	GetIdsByOrganization(organizationId string) ([]string, error)
	UpdateRootID(id string, rootNodeId string) error
	GrantUserPermission(id string, userId string, permission string) error
}

func NewWorkspaceRepo() WorkspaceRepo {
	return newWorkspaceRepo()
}

func NewWorkspace() model.Workspace {
	return &workspaceEntity{}
}

type workspaceEntity struct {
	ID               string                  `json:"id," gorm:"column:id;size:36"`
	Name             string                  `json:"name" gorm:"column:name;size:255"`
	StorageCapacity  int64                   `json:"storageCapacity" gorm:"column:storage_capacity"`
	RootID           string                  `json:"rootId" gorm:"column:root_id;size:36"`
	OrganizationId   string                  `json:"organizationId" gorm:"column:organization_id;size:36"`
	UserPermissions  []*userPermissionValue  `json:"userPermissions" gorm:"-"`
	GroupPermissions []*groupPermissionValue `json:"groupPermissions" gorm:"-"`
	Bucket           string                  `json:"bucket" gorm:"column:bucket;size:255"`
	CreateTime       string                  `json:"createTime" gorm:"column:create_time"`
	UpdateTime       *string                 `json:"updateTime,omitempty" gorm:"column:update_time"`
}

func (workspaceEntity) TableName() string {
	return "workspace"
}

func (w *workspaceEntity) BeforeCreate(tx *gorm.DB) (err error) {
	w.CreateTime = time.Now().UTC().Format(time.RFC3339)
	return nil
}

func (w *workspaceEntity) BeforeSave(tx *gorm.DB) (err error) {
	timeNow := time.Now().UTC().Format(time.RFC3339)
	w.UpdateTime = &timeNow
	return nil
}

func (w workspaceEntity) GetID() string {
	return w.ID
}

func (w workspaceEntity) GetName() string {
	return w.Name
}

func (w workspaceEntity) GetStorageCapacity() int64 {
	return w.StorageCapacity
}

func (w workspaceEntity) GetRootID() string {
	return w.RootID
}

func (w workspaceEntity) GetOrganizationID() string {
	return w.OrganizationId
}

func (w workspaceEntity) GetUserPermissions() []model.CoreUserPermission {
	var res []model.CoreUserPermission
	for _, p := range w.UserPermissions {
		res = append(res, p)
	}
	return res
}

func (w workspaceEntity) GetGroupPermissions() []model.CoreGroupPermission {
	var res []model.CoreGroupPermission
	for _, p := range w.GroupPermissions {
		res = append(res, p)
	}
	return res
}

func (w workspaceEntity) GetBucket() string {
	return w.Bucket
}

func (w workspaceEntity) GetCreateTime() string {
	return w.CreateTime
}

func (w workspaceEntity) GetUpdateTime() *string {
	return w.UpdateTime
}

func (w *workspaceEntity) SetName(name string) {
	w.Name = name
}

func (w *workspaceEntity) SetUpdateTime(updateTime *string) {
	w.UpdateTime = updateTime
}

type workspaceRepo struct {
	db             *gorm.DB
	permissionRepo *permissionRepo
}

func newWorkspaceRepo() *workspaceRepo {
	return &workspaceRepo{
		db:             infra.GetDb(),
		permissionRepo: newPermissionRepo(),
	}
}

func (repo *workspaceRepo) Insert(opts WorkspaceInsertOptions) (model.Workspace, error) {
	var id string
	if len(opts.ID) > 0 {
		id = opts.ID
	} else {
		id = helper.NewId()
	}
	workspace := workspaceEntity{
		ID:              id,
		Name:            opts.Name,
		StorageCapacity: opts.StorageCapacity,
		RootID:          opts.RootId,
		OrganizationId:  opts.OrganizationId,
		Bucket:          opts.Bucket,
	}
	if db := repo.db.Save(&workspace); db.Error != nil {
		return nil, db.Error
	}
	res, err := repo.findByID(id)
	if err != nil {
		return nil, err
	}
	if err := repo.populateModelFields([]*workspaceEntity{res}); err != nil {
		return nil, err
	}
	return res, nil
}

func (repo *workspaceRepo) findByName(name string) (*workspaceEntity, error) {
	var res = workspaceEntity{}
	db := repo.db.Where("name = ?", name).First(&res)
	if db.Error != nil {
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			return nil, errorpkg.NewWorkspaceNotFoundError(db.Error)
		} else {
			return nil, errorpkg.NewInternalServerError(db.Error)
		}
	}
	return &res, nil
}

func (repo *workspaceRepo) FindByName(name string) (model.Workspace, error) {
	workspace, err := repo.findByName(name)
	if err != nil {
		return nil, err
	}
	if err := repo.populateModelFields([]*workspaceEntity{workspace}); err != nil {
		return nil, err
	}
	return workspace, err
}

func (repo *workspaceRepo) findByID(id string) (*workspaceEntity, error) {
	var res = workspaceEntity{}
	db := repo.db.Where("id = ?", id).First(&res)
	if db.Error != nil {
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			return nil, errorpkg.NewWorkspaceNotFoundError(db.Error)
		} else {
			return nil, errorpkg.NewInternalServerError(db.Error)
		}
	}
	return &res, nil
}

func (repo *workspaceRepo) FindByID(id string) (model.Workspace, error) {
	workspace, err := repo.findByID(id)
	if err != nil {
		return nil, err
	}
	if err := repo.populateModelFields([]*workspaceEntity{workspace}); err != nil {
		return nil, err
	}
	return workspace, err
}

func (repo *workspaceRepo) UpdateName(id string, name string) (model.Workspace, error) {
	workspace, err := repo.findByID(id)
	if err != nil {
		return &workspaceEntity{}, err
	}
	workspace.Name = name
	db := repo.db.Save(&workspace)
	if db.Error != nil {
		return nil, db.Error
	}
	res, err := repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (repo *workspaceRepo) UpdateStorageCapacity(id string, storageCapacity int64) (model.Workspace, error) {
	workspace, err := repo.findByID(id)
	if err != nil {
		return &workspaceEntity{}, err
	}
	workspace.StorageCapacity = storageCapacity
	db := repo.db.Save(&workspace)
	if db.Error != nil {
		return nil, db.Error
	}
	res, err := repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (repo *workspaceRepo) Delete(id string) error {
	db := repo.db.Exec("DELETE FROM workspace WHERE id = ?", id)
	if db.Error != nil {
		return db.Error
	}
	db = repo.db.Exec("DELETE FROM userpermission WHERE resource_id = ?", id)
	if db.Error != nil {
		return db.Error
	}
	db = repo.db.Exec("DELETE FROM grouppermission WHERE resource_id = ?", id)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *workspaceRepo) GetIDs() ([]string, error) {
	type IdResult struct {
		Result string
	}
	var ids []IdResult
	db := repo.db.Raw("SELECT id result FROM workspace ORDER BY create_time DESC").Scan(&ids)
	if db.Error != nil {
		return []string{}, db.Error
	}
	res := []string{}
	for _, id := range ids {
		res = append(res, id.Result)
	}
	return res, nil
}

func (repo *workspaceRepo) GetIdsByOrganization(organizationId string) ([]string, error) {
	type IdResult struct {
		Result string
	}
	var ids []IdResult
	db := repo.db.
		Raw("SELECT id result FROM workspace WHERE organization_id = ? ORDER BY create_time DESC", organizationId).
		Scan(&ids)
	if db.Error != nil {
		return nil, db.Error
	}
	res := []string{}
	for _, id := range ids {
		res = append(res, id.Result)
	}
	return res, nil
}

func (repo *workspaceRepo) UpdateRootID(id string, rootNodeId string) error {
	db := repo.db.Exec("UPDATE workspace SET root_id = ? WHERE id = ?", rootNodeId, id)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *workspaceRepo) GrantUserPermission(id string, userId string, permission string) error {
	db := repo.db.Exec(
		"INSERT INTO userpermission (id, user_id, resource_id, permission) "+
			"VALUES (?, ?, ?, ?) ON CONFLICT (user_id, resource_id) DO UPDATE SET permission = ?",
		helper.NewId(), userId, id, permission, permission)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *workspaceRepo) populateModelFields(workspaces []*workspaceEntity) error {
	for _, w := range workspaces {
		w.UserPermissions = make([]*userPermissionValue, 0)
		userPermissions, err := repo.permissionRepo.GetUserPermissions(w.ID)
		if err != nil {
			return err
		}
		for _, p := range userPermissions {
			w.UserPermissions = append(w.UserPermissions, &userPermissionValue{
				UserId: p.UserID,
				Value:  p.Permission,
			})
		}
		w.GroupPermissions = make([]*groupPermissionValue, 0)
		groupPermissions, err := repo.permissionRepo.GetGroupPermissions(w.ID)
		if err != nil {
			return err
		}
		for _, p := range groupPermissions {
			w.GroupPermissions = append(w.GroupPermissions, &groupPermissionValue{
				GroupID: p.GroupID,
				Value:   p.Permission,
			})
		}
	}
	return nil
}
