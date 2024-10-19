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

	"gorm.io/gorm"

	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/infra"
	"github.com/kouprlabs/voltaserve/api/model"
)

type GroupRepo interface {
	Insert(opts GroupInsertOptions) (model.Group, error)
	Find(id string) (model.Group, error)
	Count() (int64, error)
	FindIDs() ([]string, error)
	FindIDsByFile(fileID string) ([]string, error)
	FindIDsByOrganization(id string) ([]string, error)
	Save(group model.Group) error
	Delete(id string) error
	FindMembers(id string) ([]model.User, error)
	CountOwners(id string) (int64, error)
	GrantUserPermission(id string, userID string, permission string) error
	RevokeUserPermission(id string, userID string) error
}

func NewGroupRepo() GroupRepo {
	return newGroupRepo()
}

func NewGroup() model.Group {
	return &groupEntity{}
}

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
	g.CreateTime = helper.NewTimestamp()
	return nil
}

func (g *groupEntity) BeforeSave(*gorm.DB) (err error) {
	timeNow := helper.NewTimestamp()
	g.UpdateTime = &timeNow
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

func (g *groupEntity) GetUsers() []string {
	return g.Members
}

func (g *groupEntity) GetCreateTime() string {
	return g.CreateTime
}

func (g *groupEntity) GetUpdateTime() *string {
	return g.UpdateTime
}

func (g *groupEntity) SetName(name string) {
	g.Name = name
}

func (g *groupEntity) SetUpdateTime(updateTime *string) {
	g.UpdateTime = updateTime
}

type groupRepo struct {
	db             *gorm.DB
	permissionRepo *permissionRepo
}

func newGroupRepo() *groupRepo {
	return &groupRepo{
		db:             infra.NewPostgresManager().GetDBOrPanic(),
		permissionRepo: newPermissionRepo(),
	}
}

type GroupInsertOptions struct {
	ID             string
	Name           string
	OrganizationID string
	OwnerID        string
}

func (repo *groupRepo) Insert(opts GroupInsertOptions) (model.Group, error) {
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

func (repo *groupRepo) find(id string) (*groupEntity, error) {
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

func (repo *groupRepo) Find(id string) (model.Group, error) {
	group, err := repo.find(id)
	if err != nil {
		return nil, err
	}
	if err := repo.populateModelFields([]*groupEntity{group}); err != nil {
		return nil, err
	}
	return group, nil
}

func (repo *groupRepo) Count() (int64, error) {
	var count int64
	db := repo.db.Model(&groupEntity{}).Count(&count)
	if db.Error != nil {
		return -1, db.Error
	}
	return count, nil
}

func (repo *groupRepo) FindIDsByFile(fileID string) ([]string, error) {
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

func (repo *groupRepo) FindIDsByOrganization(id string) ([]string, error) {
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

func (repo *groupRepo) Save(group model.Group) error {
	db := repo.db.Save(group)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *groupRepo) Delete(id string) error {
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

func (repo *groupRepo) FindIDs() ([]string, error) {
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

func (repo *groupRepo) FindMembers(id string) ([]model.User, error) {
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

func (repo *groupRepo) CountOwners(id string) (int64, error) {
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

func (repo *groupRepo) GrantUserPermission(id string, userID string, permission string) error {
	db := repo.db.
		Exec(`INSERT INTO userpermission (id, user_id, resource_id, permission, create_time)
              VALUES (?, ?, ?, ?, ?) ON CONFLICT (user_id, resource_id) DO UPDATE SET permission = ?`,
			helper.NewID(), userID, id, permission, helper.NewTimestamp(), permission)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *groupRepo) RevokeUserPermission(id string, userID string) error {
	db := repo.db.Exec("DELETE FROM userpermission WHERE user_id = ? AND resource_id = ?", userID, id)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *groupRepo) populateModelFields(groups []*groupEntity) error {
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
