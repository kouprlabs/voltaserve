package repo

import (
	"errors"
	"time"
	"voltaserve/errorpkg"
	"voltaserve/helpers"
	"voltaserve/infra"
	"voltaserve/model"

	"gorm.io/gorm"
)

type WorkspaceEntity struct {
	Id               string                   `json:"id," gorm:"column:id;size:36"`
	Name             string                   `json:"name" gorm:"column:name;size:255"`
	StorageCapacity  int64                    `json:"storageCapacity" gorm:"column:storage_capacity"`
	RootId           string                   `json:"rootId" gorm:"column:root_id;size:36"`
	OrganizationId   string                   `json:"organizationId" gorm:"column:organization_id;size:36"`
	UserPermissions  []*model.UserPermission  `json:"userPermissions" gorm:"-"`
	GroupPermissions []*model.GroupPermission `json:"groupPermissions" gorm:"-"`
	Bucket           string                   `json:"bucket" gorm:"column:bucket;size:255"`
	CreateTime       string                   `json:"createTime" gorm:"column:create_time"`
	UpdateTime       *string                  `json:"updateTime,omitempty" gorm:"column:update_time"`
}

func (WorkspaceEntity) TableName() string {
	return "workspace"
}

func (w *WorkspaceEntity) BeforeCreate(tx *gorm.DB) (err error) {
	w.CreateTime = time.Now().UTC().Format(time.RFC3339)
	return nil
}

func (w *WorkspaceEntity) BeforeSave(tx *gorm.DB) (err error) {
	timeNow := time.Now().UTC().Format(time.RFC3339)
	w.UpdateTime = &timeNow
	return nil
}

func (w WorkspaceEntity) GetId() string {
	return w.Id
}

func (w WorkspaceEntity) GetName() string {
	return w.Name
}

func (w WorkspaceEntity) GetStorageCapacity() int64 {
	return w.StorageCapacity
}

func (w WorkspaceEntity) GetRootId() string {
	return w.RootId
}

func (w WorkspaceEntity) GetOrganizationId() string {
	return w.OrganizationId
}

func (w WorkspaceEntity) GetUserPermissions() []model.UserPermissionModel {
	var res []model.UserPermissionModel
	for _, p := range w.UserPermissions {
		res = append(res, p)
	}
	return res
}

func (w WorkspaceEntity) GetGroupPermissions() []model.GroupPermissionModel {
	var res []model.GroupPermissionModel
	for _, p := range w.GroupPermissions {
		res = append(res, p)
	}
	return res
}

func (w WorkspaceEntity) GetBucket() string {
	return w.Bucket
}

func (w WorkspaceEntity) GetCreateTime() string {
	return w.CreateTime
}

func (w WorkspaceEntity) GetUpdateTime() *string {
	return w.UpdateTime
}

func (w *WorkspaceEntity) SetName(name string) {
	w.Name = name
}

func (w *WorkspaceEntity) SetUpdateTime(updateTime *string) {
	w.UpdateTime = updateTime
}

type WorkspaceInsertOptions struct {
	Id              string
	Name            string
	StorageCapacity int64
	Image           *string
	OrganizationId  string
	RootId          string
	Bucket          string
}

type WorkspaceRepo struct {
	db             *gorm.DB
	permissionRepo *PermissionRepo
}

func NewWorkspaceRepo() *WorkspaceRepo {
	return &WorkspaceRepo{
		db:             infra.GetDb(),
		permissionRepo: NewPermissionRepo(),
	}
}

func (repo *WorkspaceRepo) Insert(opts WorkspaceInsertOptions) (model.WorkspaceModel, error) {
	var id string
	if len(opts.Id) > 0 {
		id = opts.Id
	} else {
		id = helpers.NewId()
	}
	workspace := WorkspaceEntity{
		Id:              id,
		Name:            opts.Name,
		StorageCapacity: opts.StorageCapacity,
		RootId:          opts.RootId,
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
	if err := repo.populateModelFields([]*WorkspaceEntity{res}); err != nil {
		return nil, err
	}
	return res, nil
}

func (repo *WorkspaceRepo) findByName(name string) (*WorkspaceEntity, error) {
	var res = WorkspaceEntity{}
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

func (repo *WorkspaceRepo) FindByName(name string) (model.WorkspaceModel, error) {
	workspace, err := repo.findByName(name)
	if err != nil {
		return nil, err
	}
	if err := repo.populateModelFields([]*WorkspaceEntity{workspace}); err != nil {
		return nil, err
	}
	return workspace, err
}

func (repo *WorkspaceRepo) findByID(id string) (*WorkspaceEntity, error) {
	var res = WorkspaceEntity{}
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

func (repo *WorkspaceRepo) FindByID(id string) (model.WorkspaceModel, error) {
	workspace, err := repo.findByID(id)
	if err != nil {
		return nil, err
	}
	if err := repo.populateModelFields([]*WorkspaceEntity{workspace}); err != nil {
		return nil, err
	}
	return workspace, err
}

func (repo *WorkspaceRepo) UpdateName(id string, name string) (model.WorkspaceModel, error) {
	workspace, err := repo.findByID(id)
	if err != nil {
		return &WorkspaceEntity{}, err
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

func (repo *WorkspaceRepo) UpdateStorageCapacity(id string, storageCapacity int64) (model.WorkspaceModel, error) {
	workspace, err := repo.findByID(id)
	if err != nil {
		return &WorkspaceEntity{}, err
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

func (repo *WorkspaceRepo) Delete(id string) error {
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

func (repo *WorkspaceRepo) GetIds() ([]string, error) {
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

func (repo *WorkspaceRepo) GetIdsByOrganization(organizationId string) ([]string, error) {
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

func (repo *WorkspaceRepo) UpdateRootId(id string, rootNodeId string) error {
	db := repo.db.Exec("UPDATE workspace SET root_id = ? WHERE id = ?", rootNodeId, id)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *WorkspaceRepo) GrantUserPermission(id string, userId string, permission string) error {
	db := repo.db.Exec(
		"INSERT INTO userpermission (id, user_id, resource_id, permission) "+
			"VALUES (?, ?, ?, ?) ON CONFLICT (user_id, resource_id) DO UPDATE SET permission = ?",
		helpers.NewId(), userId, id, permission, permission)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *WorkspaceRepo) populateModelFields(workspaces []*WorkspaceEntity) error {
	for _, w := range workspaces {
		w.UserPermissions = make([]*model.UserPermission, 0)
		userPermissions, err := repo.permissionRepo.GetUserPermissions(w.Id)
		if err != nil {
			return err
		}
		for _, p := range userPermissions {
			w.UserPermissions = append(w.UserPermissions, &model.UserPermission{
				UserId: p.UserId,
				Value:  p.Permission,
			})
		}
		w.GroupPermissions = make([]*model.GroupPermission, 0)
		groupPermissions, err := repo.permissionRepo.GetGroupPermissions(w.Id)
		if err != nil {
			return err
		}
		for _, p := range groupPermissions {
			w.GroupPermissions = append(w.GroupPermissions, &model.GroupPermission{
				GroupId: p.GroupId,
				Value:   p.Permission,
			})
		}
	}
	return nil
}
