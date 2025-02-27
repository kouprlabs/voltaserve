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
	"errors"

	"gorm.io/gorm"

	"github.com/kouprlabs/voltaserve/shared/errorpkg"
	"github.com/kouprlabs/voltaserve/shared/helper"
	"github.com/kouprlabs/voltaserve/shared/infra"
	"github.com/kouprlabs/voltaserve/shared/model"

	"github.com/kouprlabs/voltaserve/api/config"
)

type workspaceEntity struct {
	ID               string                  `gorm:"column:id;size:36"              json:"id"`
	Name             string                  `gorm:"column:name;size:255"           json:"name"`
	StorageCapacity  int64                   `gorm:"column:storage_capacity"        json:"storageCapacity"`
	RootID           *string                 `gorm:"column:root_id;size:36"         json:"rootId"`
	OrganizationID   string                  `gorm:"column:organization_id;size:36" json:"organizationId"`
	UserPermissions  []*UserPermissionValue  `gorm:"-"                              json:"userPermissions"`
	GroupPermissions []*GroupPermissionValue `gorm:"-"                              json:"groupPermissions"`
	Bucket           string                  `gorm:"column:bucket;size:255"         json:"bucket"`
	CreateTime       string                  `gorm:"column:create_time"             json:"createTime"`
	UpdateTime       *string                 `gorm:"column:update_time"             json:"updateTime,omitempty"`
}

func (*workspaceEntity) TableName() string {
	return "workspace"
}

func (w *workspaceEntity) BeforeCreate(*gorm.DB) (err error) {
	w.CreateTime = helper.NewTimeString()
	return nil
}

func (w *workspaceEntity) BeforeSave(*gorm.DB) (err error) {
	w.UpdateTime = helper.ToPtr(helper.NewTimeString())
	return nil
}

func (w *workspaceEntity) GetID() string {
	return w.ID
}

func (w *workspaceEntity) GetName() string {
	return w.Name
}

func (w *workspaceEntity) GetStorageCapacity() int64 {
	return w.StorageCapacity
}

func (w *workspaceEntity) GetRootID() string {
	if w.RootID == nil {
		return ""
	}

	return *w.RootID
}

func (w *workspaceEntity) GetOrganizationID() string {
	return w.OrganizationID
}

func (w *workspaceEntity) GetUserPermissions() []model.CoreUserPermission {
	var res []model.CoreUserPermission
	for _, p := range w.UserPermissions {
		res = append(res, p)
	}
	return res
}

func (w *workspaceEntity) GetGroupPermissions() []model.CoreGroupPermission {
	var res []model.CoreGroupPermission
	for _, p := range w.GroupPermissions {
		res = append(res, p)
	}
	return res
}

func (w *workspaceEntity) GetBucket() string {
	return w.Bucket
}

func (w *workspaceEntity) GetCreateTime() string {
	return w.CreateTime
}

func (w *workspaceEntity) GetUpdateTime() *string {
	return w.UpdateTime
}

func (w *workspaceEntity) SetID(id string) {
	w.ID = id
}

func (w *workspaceEntity) SetName(name string) {
	w.Name = name
}

func (w *workspaceEntity) SetStorageCapacity(storageCapacity int64) {
	w.StorageCapacity = storageCapacity
}

func (w *workspaceEntity) SetRootID(rootID string) {
	w.RootID = &rootID
}

func (w *workspaceEntity) SetOrganizationID(organizationID string) {
	w.OrganizationID = organizationID
}

func (w *workspaceEntity) SetUserPermissions(permissions []model.CoreUserPermission) {
	w.UserPermissions = make([]*UserPermissionValue, len(permissions))
	for i, p := range permissions {
		w.UserPermissions[i] = p.(*UserPermissionValue)
	}
}

func (w *workspaceEntity) SetGroupPermissions(permissions []model.CoreGroupPermission) {
	w.GroupPermissions = make([]*GroupPermissionValue, len(permissions))
	for i, p := range permissions {
		w.GroupPermissions[i] = p.(*GroupPermissionValue)
	}
}

func (w *workspaceEntity) SetBucket(bucket string) {
	w.Bucket = bucket
}

func (w *workspaceEntity) SetCreateTime(createTime string) {
	w.CreateTime = createTime
}

func (w *workspaceEntity) SetUpdateTime(updateTime *string) {
	w.UpdateTime = updateTime
}

func NewWorkspaceModel() model.Workspace {
	return &workspaceEntity{}
}

type WorkspaceNewModelOptions struct {
	ID               string
	Name             string
	StorageCapacity  int64
	Image            *string
	OrganizationID   string
	RootID           *string
	Bucket           string
	UserPermissions  []model.CoreUserPermission
	GroupPermissions []model.CoreGroupPermission
	CreateTime       string
	UpdateTime       *string
}

func NewWorkspaceModelWithOptions(opts WorkspaceNewModelOptions) model.Workspace {
	res := &workspaceEntity{
		ID:              opts.ID,
		Name:            opts.Name,
		StorageCapacity: opts.StorageCapacity,
		RootID:          opts.RootID,
		OrganizationID:  opts.OrganizationID,
		Bucket:          opts.Bucket,
		CreateTime:      opts.CreateTime,
		UpdateTime:      opts.UpdateTime,
	}
	res.SetUserPermissions(opts.UserPermissions)
	res.SetGroupPermissions(opts.GroupPermissions)
	return res
}

type WorkspaceRepo struct {
	db             *gorm.DB
	permissionRepo *PermissionRepo
}

func NewWorkspaceRepo() *WorkspaceRepo {
	return &WorkspaceRepo{
		db: infra.NewPostgresManager(
			config.GetConfig().Postgres,
			config.GetConfig().Environment,
		).GetDBOrPanic(),
		permissionRepo: NewPermissionRepo(),
	}
}

type WorkspaceInsertOptions struct {
	ID              string
	Name            string
	StorageCapacity int64
	Image           *string
	OrganizationID  string
	RootID          *string
	Bucket          string
}

func (repo *WorkspaceRepo) Insert(opts WorkspaceInsertOptions) (model.Workspace, error) {
	var id string
	if len(opts.ID) > 0 {
		id = opts.ID
	} else {
		id = helper.NewID()
	}
	workspace := workspaceEntity{
		ID:              id,
		Name:            opts.Name,
		StorageCapacity: opts.StorageCapacity,
		RootID:          opts.RootID,
		OrganizationID:  opts.OrganizationID,
		Bucket:          opts.Bucket,
	}
	if db := repo.db.Create(&workspace); db.Error != nil {
		return nil, db.Error
	}
	res, err := repo.find(id)
	if err != nil {
		return nil, err
	}
	if err := repo.populateModelFields([]*workspaceEntity{res}); err != nil {
		return nil, err
	}
	return res, nil
}

func (repo *WorkspaceRepo) Find(id string) (model.Workspace, error) {
	workspace, err := repo.find(id)
	if err != nil {
		return nil, err
	}
	if err := repo.populateModelFields([]*workspaceEntity{workspace}); err != nil {
		return nil, err
	}
	return workspace, err
}

func (repo *WorkspaceRepo) FindOrNil(id string) model.Workspace {
	res, err := repo.Find(id)
	if err != nil {
		return nil
	}
	return res
}

func (repo *WorkspaceRepo) Count() (int64, error) {
	var count int64
	db := repo.db.Model(&workspaceEntity{}).Count(&count)
	if db.Error != nil {
		return -1, db.Error
	}
	return count, nil
}

func (repo *WorkspaceRepo) UpdateName(id string, name string) (model.Workspace, error) {
	workspace, err := repo.find(id)
	if err != nil {
		return &workspaceEntity{}, err
	}
	workspace.Name = name
	if db := repo.db.Save(&workspace); db.Error != nil {
		return nil, db.Error
	}
	res, err := repo.Find(id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (repo *WorkspaceRepo) UpdateStorageCapacity(id string, storageCapacity int64) (model.Workspace, error) {
	workspace, err := repo.find(id)
	if err != nil {
		return &workspaceEntity{}, err
	}
	workspace.StorageCapacity = storageCapacity
	db := repo.db.Save(&workspace)
	if db.Error != nil {
		return nil, db.Error
	}
	res, err := repo.Find(id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (repo *WorkspaceRepo) UpdateRootID(id string, rootNodeID string) error {
	db := repo.db.Exec("UPDATE workspace SET root_id = ? WHERE id = ?", rootNodeID, id)
	if db.Error != nil {
		return db.Error
	}
	return nil
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

func (repo *WorkspaceRepo) FindIDs() ([]string, error) {
	type IDResult struct {
		Result string
	}
	var ids []IDResult
	db := repo.db.Raw("SELECT id result FROM workspace ORDER BY create_time DESC").Scan(&ids)
	if db.Error != nil {
		return []string{}, db.Error
	}
	res := make([]string, 0)
	for _, id := range ids {
		res = append(res, id.Result)
	}
	return res, nil
}

func (repo *WorkspaceRepo) FindIDsByOrganization(orgID string) ([]string, error) {
	type IDResult struct {
		Result string
	}
	var ids []IDResult
	db := repo.db.
		Raw("SELECT id result FROM workspace WHERE organization_id = ? ORDER BY create_time DESC", orgID).
		Scan(&ids)
	if db.Error != nil {
		return nil, db.Error
	}
	res := make([]string, 0)
	for _, id := range ids {
		res = append(res, id.Result)
	}
	return res, nil
}

func (repo *WorkspaceRepo) GrantUserPermission(id string, userID string, permission string) error {
	db := repo.db.
		Exec(`INSERT INTO userpermission (id, user_id, resource_id, permission, create_time)
              VALUES (?, ?, ?, ?, ?)
              ON CONFLICT (user_id, resource_id) DO UPDATE SET permission = ?`,
			helper.NewID(), userID, id, permission, helper.NewTimeString(), permission)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *WorkspaceRepo) RevokeUserPermission(id string, userID string) error {
	db := repo.db.Exec("DELETE FROM userpermission WHERE user_id = ? AND resource_id = ?", userID, id)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *WorkspaceRepo) find(id string) (*workspaceEntity, error) {
	res := workspaceEntity{}
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

func (repo *WorkspaceRepo) populateModelFields(workspaces []*workspaceEntity) error {
	for _, w := range workspaces {
		w.UserPermissions = make([]*UserPermissionValue, 0)
		userPermissions, err := repo.permissionRepo.FindUserPermissions(w.ID)
		if err != nil {
			return err
		}
		for _, p := range userPermissions {
			w.UserPermissions = append(w.UserPermissions, &UserPermissionValue{
				UserID: p.GetUserID(),
				Value:  p.GetPermission(),
			})
		}
		w.GroupPermissions = make([]*GroupPermissionValue, 0)
		groupPermissions, err := repo.permissionRepo.FindGroupPermissions(w.ID)
		if err != nil {
			return err
		}
		for _, p := range groupPermissions {
			w.GroupPermissions = append(w.GroupPermissions, &GroupPermissionValue{
				GroupID: p.GetGroupID(),
				Value:   p.GetPermission(),
			})
		}
	}
	return nil
}
