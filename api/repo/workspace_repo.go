// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

package repo

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/infra"
	"github.com/kouprlabs/voltaserve/api/model"
)

type WorkspaceRepo interface {
	Insert(opts WorkspaceInsertOptions) (model.Workspace, error)
	Find(id string) (model.Workspace, error)
	Count() (int64, error)
	UpdateName(id string, name string) (model.Workspace, error)
	UpdateStorageCapacity(id string, storageCapacity int64) (model.Workspace, error)
	UpdateRootID(id string, rootNodeID string) error
	Delete(id string) error
	GetIDs() ([]string, error)
	GetIDsByOrganization(orgID string) ([]string, error)
	GrantUserPermission(id string, userID string, permission string) error
}

func NewWorkspaceRepo() WorkspaceRepo {
	return newWorkspaceRepo()
}

func NewWorkspace() model.Workspace {
	return &workspaceEntity{}
}

type workspaceEntity struct {
	ID               string                  `gorm:"column:id;size:36"              json:"id"`
	Name             string                  `gorm:"column:name;size:255"           json:"name"`
	StorageCapacity  int64                   `gorm:"column:storage_capacity"        json:"storageCapacity"`
	RootID           string                  `gorm:"column:root_id;size:36"         json:"rootId"`
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
	w.CreateTime = time.Now().UTC().Format(time.RFC3339)
	return nil
}

func (w *workspaceEntity) BeforeSave(*gorm.DB) (err error) {
	timeNow := time.Now().UTC().Format(time.RFC3339)
	w.UpdateTime = &timeNow
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
	return w.RootID
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

func (w *workspaceEntity) SetName(name string) {
	w.Name = name
}

type workspaceRepo struct {
	db             *gorm.DB
	permissionRepo *permissionRepo
}

func newWorkspaceRepo() *workspaceRepo {
	return &workspaceRepo{
		db:             infra.NewPostgresManager().GetDBOrPanic(),
		permissionRepo: newPermissionRepo(),
	}
}

type WorkspaceInsertOptions struct {
	ID              string
	Name            string
	StorageCapacity int64
	Image           *string
	OrganizationID  string
	RootID          string
	Bucket          string
}

func (repo *workspaceRepo) Insert(opts WorkspaceInsertOptions) (model.Workspace, error) {
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

func (repo *workspaceRepo) find(id string) (*workspaceEntity, error) {
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

func (repo *workspaceRepo) Find(id string) (model.Workspace, error) {
	workspace, err := repo.find(id)
	if err != nil {
		return nil, err
	}
	if err := repo.populateModelFields([]*workspaceEntity{workspace}); err != nil {
		return nil, err
	}
	return workspace, err
}

func (repo *workspaceRepo) Count() (int64, error) {
	type Result struct {
		Result int64
	}
	var res Result
	db := repo.db.
		Raw("SELECT count(*) as result FROM workspace").
		Scan(&res)
	if db.Error != nil {
		return 0, db.Error
	}
	return res.Result, nil
}

func (repo *workspaceRepo) UpdateName(id string, name string) (model.Workspace, error) {
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

func (repo *workspaceRepo) UpdateStorageCapacity(id string, storageCapacity int64) (model.Workspace, error) {
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

func (repo *workspaceRepo) UpdateRootID(id string, rootNodeID string) error {
	db := repo.db.Exec("UPDATE workspace SET root_id = ? WHERE id = ?", rootNodeID, id)
	if db.Error != nil {
		return db.Error
	}
	return nil
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
	type IDResult struct {
		Result string
	}
	var ids []IDResult
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

func (repo *workspaceRepo) GetIDsByOrganization(orgID string) ([]string, error) {
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
	res := []string{}
	for _, id := range ids {
		res = append(res, id.Result)
	}
	return res, nil
}

func (repo *workspaceRepo) GrantUserPermission(id string, userID string, permission string) error {
	db := repo.db.
		Exec(`INSERT INTO userpermission (id, user_id, resource_id, permission, create_time)
              VALUES (?, ?, ?, ?, ?)
              ON CONFLICT (user_id, resource_id) DO UPDATE SET permission = ?`,
			helper.NewID(), userID, id, permission, time.Now().UTC().Format(time.RFC3339), permission)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *workspaceRepo) populateModelFields(workspaces []*workspaceEntity) error {
	for _, w := range workspaces {
		w.UserPermissions = make([]*UserPermissionValue, 0)
		userPermissions, err := repo.permissionRepo.GetUserPermissions(w.ID)
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
		groupPermissions, err := repo.permissionRepo.GetGroupPermissions(w.ID)
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
