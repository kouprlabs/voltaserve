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

	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/infra"
	"github.com/kouprlabs/voltaserve/api/model"
)

type groupEntity struct {
	ID               string                  `gorm:"column:id"              json:"id"`
	Name             string                  `gorm:"column:name"            json:"name"`
	OrganizationID   string                  `gorm:"column:organization_id" json:"organizationId"`
	UserPermissions  []*UserPermissionValue  `gorm:"-"                      json:"userPermissions"`
	GroupPermissions []*GroupPermissionValue `gorm:"-"                      json:"groupPermissions"`
	Members          []string                `gorm:"-"                      json:"members"`
	CreateTime       string                  `gorm:"column:create_time"     json:"createTime"`
	UpdateTime       *string                 `gorm:"column:update_time"     json:"updateTime"`
}

func (*groupEntity) TableName() string {
	return "group"
}

func (g *groupEntity) BeforeCreate(*gorm.DB) (err error) {
	g.CreateTime = helper.NewTimeString()
	return nil
}

func (g *groupEntity) BeforeSave(*gorm.DB) (err error) {
	g.UpdateTime = helper.ToPtr(helper.NewTimeString())
	return nil
}

func (g *groupEntity) GetID() string {
	return g.ID
}

func (g *groupEntity) GetName() string {
	return g.Name
}

func (g *groupEntity) GetOrganizationID() string {
	return g.OrganizationID
}

func (g *groupEntity) GetUserPermissions() []model.CoreUserPermission {
	var res []model.CoreUserPermission
	for _, p := range g.UserPermissions {
		res = append(res, p)
	}
	return res
}

func (g *groupEntity) GetGroupPermissions() []model.CoreGroupPermission {
	var res []model.CoreGroupPermission
	for _, p := range g.GroupPermissions {
		res = append(res, p)
	}
	return res
}

func (g *groupEntity) GetMembers() []string {
	return g.Members
}

func (g *groupEntity) GetCreateTime() string {
	return g.CreateTime
}

func (g *groupEntity) GetUpdateTime() *string {
	return g.UpdateTime
}

func (g *groupEntity) SetID(id string) {
	g.ID = id
}

func (g *groupEntity) SetName(name string) {
	g.Name = name
}

func (g *groupEntity) SetOrganizationID(id string) {
	g.OrganizationID = id
}

func (w *groupEntity) SetUserPermissions(permissions []model.CoreUserPermission) {
	w.UserPermissions = make([]*UserPermissionValue, len(permissions))
	for i, p := range permissions {
		w.UserPermissions[i] = p.(*UserPermissionValue)
	}
}

func (w *groupEntity) SetGroupPermissions(permissions []model.CoreGroupPermission) {
	w.GroupPermissions = make([]*GroupPermissionValue, len(permissions))
	for i, p := range permissions {
		w.GroupPermissions[i] = p.(*GroupPermissionValue)
	}
}

func (g *groupEntity) SetMembers(members []string) {
	g.Members = members
}

func (g *groupEntity) SetCreateTime(createTime string) {
	g.CreateTime = createTime
}

func (g *groupEntity) SetUpdateTime(updateTime *string) {
	g.UpdateTime = updateTime
}

func NewGroupModel() model.Group {
	return &groupEntity{}
}

type GroupNewModelOptions struct {
	ID               string
	Name             string
	OrganizationID   string
	UserPermissions  []model.CoreUserPermission
	GroupPermissions []model.CoreGroupPermission
	Members          []string
	CreateTime       string
	UpdateTime       *string
}

func NewGroupModelWithOptions(opts GroupNewModelOptions) model.Group {
	res := &groupEntity{
		ID:             opts.ID,
		Name:           opts.Name,
		OrganizationID: opts.OrganizationID,
		Members:        make([]string, len(opts.Members)),
		CreateTime:     opts.CreateTime,
		UpdateTime:     opts.UpdateTime,
	}
	res.SetUserPermissions(opts.UserPermissions)
	res.SetGroupPermissions(opts.GroupPermissions)
	return res
}

type GroupRepo struct {
	db             *gorm.DB
	permissionRepo *PermissionRepo
}

func NewGroupRepo() *GroupRepo {
	return &GroupRepo{
		db:             infra.NewPostgresManager().GetDBOrPanic(),
		permissionRepo: NewPermissionRepo(),
	}
}

type GroupInsertOptions struct {
	ID             string
	Name           string
	OrganizationID string
	OwnerID        string
}

func (repo *GroupRepo) Insert(opts GroupInsertOptions) (model.Group, error) {
	group := groupEntity{
		ID:             opts.ID,
		Name:           opts.Name,
		OrganizationID: opts.OrganizationID,
	}
	if db := repo.db.Create(&group); db.Error != nil {
		return nil, db.Error
	}
	res, err := repo.Find(opts.ID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (repo *GroupRepo) find(id string) (*groupEntity, error) {
	res := groupEntity{}
	db := repo.db.Where("id = ?", id).First(&res)
	if db.Error != nil {
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			return nil, errorpkg.NewGroupNotFoundError(db.Error)
		} else {
			return nil, errorpkg.NewInternalServerError(db.Error)
		}
	}
	return &res, nil
}

func (repo *GroupRepo) Find(id string) (model.Group, error) {
	group, err := repo.find(id)
	if err != nil {
		return nil, err
	}
	if err := repo.populateModelFields([]*groupEntity{group}); err != nil {
		return nil, err
	}
	return group, nil
}

func (repo *GroupRepo) FindOrNil(id string) model.Group {
	res, err := repo.Find(id)
	if err != nil {
		return nil
	}
	return res
}

func (repo *GroupRepo) Count() (int64, error) {
	var count int64
	db := repo.db.Model(&groupEntity{}).Count(&count)
	if db.Error != nil {
		return -1, db.Error
	}
	return count, nil
}

func (repo *GroupRepo) FindIDsByFile(fileID string) ([]string, error) {
	type Value struct {
		Result string
	}
	var values []Value
	db := repo.db.
		Raw(`SELECT DISTINCT g.id as result FROM "group" g
             INNER JOIN grouppermission p ON p.resource_id = ?
			 WHERE p.group_id = g.id`,
			fileID).
		Scan(&values)
	if db.Error != nil {
		return []string{}, db.Error
	}
	res := make([]string, 0)
	for _, v := range values {
		res = append(res, v.Result)
	}
	return res, nil
}

func (repo *GroupRepo) FindIDsByOrganization(id string) ([]string, error) {
	type Value struct {
		Result string
	}
	var values []Value
	db := repo.db.Raw(`SELECT id as result from "group" WHERE organization_id = ?`, id).Scan(&values)
	if db.Error != nil {
		return []string{}, db.Error
	}
	res := make([]string, 0)
	for _, v := range values {
		res = append(res, v.Result)
	}
	return res, nil
}

func (repo *GroupRepo) Save(group model.Group) error {
	db := repo.db.Save(group)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *GroupRepo) Delete(id string) error {
	db := repo.db.Exec(`DELETE FROM "group" WHERE id = ?`, id)
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

func (repo *GroupRepo) FindIDs() ([]string, error) {
	type Value struct {
		Result string
	}
	var values []Value
	db := repo.db.Raw(`SELECT id result FROM "group" ORDER BY create_time DESC`).Scan(&values)
	if db.Error != nil {
		return []string{}, db.Error
	}
	res := make([]string, 0)
	for _, v := range values {
		res = append(res, v.Result)
	}
	return res, nil
}

func (repo *GroupRepo) FindMembers(id string) ([]model.User, error) {
	var entities []*userEntity
	db := repo.db.
		Raw(`SELECT u.* FROM "user" u INNER JOIN userpermission up on
             u.id = up.user_id AND up.resource_id = ?`, id).
		Scan(&entities)
	if db.Error != nil {
		return nil, db.Error
	}
	var res []model.User
	for _, u := range entities {
		res = append(res, u)
	}
	return res, nil
}

func (repo *GroupRepo) CountOwners(id string) (int64, error) {
	var count int64
	db := repo.db.Model(&userPermissionEntity{}).
		Where("resource_id = ?", id).
		Where("permission = ?", model.PermissionOwner).
		Count(&count)
	if db.Error != nil {
		return -1, db.Error
	}
	return count, nil
}

func (repo *GroupRepo) GrantUserPermission(id string, userID string, permission string) error {
	db := repo.db.
		Exec(`INSERT INTO userpermission (id, user_id, resource_id, permission, create_time)
              VALUES (?, ?, ?, ?, ?) ON CONFLICT (user_id, resource_id) DO UPDATE SET permission = ?`,
			helper.NewID(), userID, id, permission, helper.NewTimeString(), permission)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *GroupRepo) RevokeUserPermission(id string, userID string) error {
	db := repo.db.Exec("DELETE FROM userpermission WHERE user_id = ? AND resource_id = ?", userID, id)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *GroupRepo) populateModelFields(groups []*groupEntity) error {
	for _, g := range groups {
		g.UserPermissions = make([]*UserPermissionValue, 0)
		userPermissions, err := repo.permissionRepo.FindUserPermissions(g.ID)
		if err != nil {
			return err
		}
		for _, p := range userPermissions {
			g.UserPermissions = append(g.UserPermissions, &UserPermissionValue{
				UserID: p.GetUserID(),
				Value:  p.GetPermission(),
			})
		}
		g.GroupPermissions = make([]*GroupPermissionValue, 0)
		groupPermissions, err := repo.permissionRepo.FindGroupPermissions(g.ID)
		if err != nil {
			return err
		}
		for _, p := range groupPermissions {
			g.GroupPermissions = append(g.GroupPermissions, &GroupPermissionValue{
				GroupID: p.GetGroupID(),
				Value:   p.GetPermission(),
			})
		}
		members, err := repo.FindMembers(g.ID)
		if err != nil {
			return nil
		}
		g.Members = make([]string, 0)
		for _, u := range members {
			g.Members = append(g.Members, u.GetID())
		}
	}
	return nil
}
