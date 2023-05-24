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

type PostgresWorkspace struct {
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

func (PostgresWorkspace) TableName() string {
	return "workspace"
}

func (w *PostgresWorkspace) BeforeCreate(tx *gorm.DB) (err error) {
	w.CreateTime = time.Now().UTC().Format(time.RFC3339)
	return nil
}

func (w *PostgresWorkspace) BeforeSave(tx *gorm.DB) (err error) {
	timeNow := time.Now().UTC().Format(time.RFC3339)
	w.UpdateTime = &timeNow
	return nil
}

func (w PostgresWorkspace) GetId() string {
	return w.Id
}

func (w PostgresWorkspace) GetName() string {
	return w.Name
}

func (w PostgresWorkspace) GetStorageCapacity() int64 {
	return w.StorageCapacity
}

func (w PostgresWorkspace) GetRootId() string {
	return w.RootId
}

func (w PostgresWorkspace) GetOrganizationId() string {
	return w.OrganizationId
}

func (w PostgresWorkspace) GetUserPermissions() []model.UserPermissionModel {
	var res []model.UserPermissionModel
	for _, p := range w.UserPermissions {
		res = append(res, p)
	}
	return res
}

func (w PostgresWorkspace) GetGroupPermissions() []model.GroupPermissionModel {
	var res []model.GroupPermissionModel
	for _, p := range w.GroupPermissions {
		res = append(res, p)
	}
	return res
}

func (w PostgresWorkspace) GetBucket() string {
	return w.Bucket
}

func (w PostgresWorkspace) GetCreateTime() string {
	return w.CreateTime
}

func (w PostgresWorkspace) GetUpdateTime() *string {
	return w.UpdateTime
}

func (w *PostgresWorkspace) SetName(name string) {
	w.Name = name
}

func (w *PostgresWorkspace) SetUpdateTime(updateTime *string) {
	w.UpdateTime = updateTime
}

type PostgresWorkspaceRepo struct {
	db             *gorm.DB
	permissionRepo *PostgresPermissionRepo
}

func NewPostgresWorkspaceRepo() *PostgresWorkspaceRepo {
	return &PostgresWorkspaceRepo{
		db:             infra.GetDb(),
		permissionRepo: NewPostgresPermissionRepo(),
	}
}

func (repo *PostgresWorkspaceRepo) Insert(opts WorkspaceInsertOptions) (model.WorkspaceModel, error) {
	var id string
	if len(opts.Id) > 0 {
		id = opts.Id
	} else {
		id = helpers.NewId()
	}
	workspace := PostgresWorkspace{
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
	if err := repo.populateModelFields([]*PostgresWorkspace{res}); err != nil {
		return nil, err
	}
	return res, nil
}

func (repo *PostgresWorkspaceRepo) findByName(name string) (*PostgresWorkspace, error) {
	var res = PostgresWorkspace{}
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

func (repo *PostgresWorkspaceRepo) FindByName(name string) (model.WorkspaceModel, error) {
	workspace, err := repo.findByName(name)
	if err != nil {
		return nil, err
	}
	if err := repo.populateModelFields([]*PostgresWorkspace{workspace}); err != nil {
		return nil, err
	}
	return workspace, err
}

func (repo *PostgresWorkspaceRepo) findByID(id string) (*PostgresWorkspace, error) {
	var res = PostgresWorkspace{}
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

func (repo *PostgresWorkspaceRepo) FindByID(id string) (model.WorkspaceModel, error) {
	workspace, err := repo.findByID(id)
	if err != nil {
		return nil, err
	}
	if err := repo.populateModelFields([]*PostgresWorkspace{workspace}); err != nil {
		return nil, err
	}
	return workspace, err
}

func (repo *PostgresWorkspaceRepo) UpdateName(id string, name string) (model.WorkspaceModel, error) {
	workspace, err := repo.findByID(id)
	if err != nil {
		return &PostgresWorkspace{}, err
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

func (repo *PostgresWorkspaceRepo) UpdateStorageCapacity(id string, storageCapacity int64) (model.WorkspaceModel, error) {
	workspace, err := repo.findByID(id)
	if err != nil {
		return &PostgresWorkspace{}, err
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

func (repo *PostgresWorkspaceRepo) Delete(id string) error {
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

func (repo *PostgresWorkspaceRepo) GetIDs() ([]string, error) {
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

func (repo *PostgresWorkspaceRepo) GetIdsByOrganization(organizationId string) ([]string, error) {
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

func (repo *PostgresWorkspaceRepo) UpdateRootID(id string, rootNodeId string) error {
	db := repo.db.Exec("UPDATE workspace SET root_id = ? WHERE id = ?", rootNodeId, id)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *PostgresWorkspaceRepo) GrantUserPermission(id string, userId string, permission string) error {
	db := repo.db.Exec(
		"INSERT INTO userpermission (id, user_id, resource_id, permission) "+
			"VALUES (?, ?, ?, ?) ON CONFLICT (user_id, resource_id) DO UPDATE SET permission = ?",
		helpers.NewId(), userId, id, permission, permission)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *PostgresWorkspaceRepo) populateModelFields(workspaces []*PostgresWorkspace) error {
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
