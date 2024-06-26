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

type GroupRepo interface {
	Insert(opts GroupInsertOptions) (model.Group, error)
	Find(id string) (model.Group, error)
	GetIDsForFile(fileID string) ([]string, error)
	GetIDsForUser(userID string) ([]string, error)
	GetIDsForOrganization(id string) ([]string, error)
	Save(group model.Group) error
	Delete(id string) error
	AddUser(id string, userID string) error
	RemoveMember(id string, userID string) error
	GetIDs() ([]string, error)
	GetMembers(id string) ([]model.User, error)
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
	ID               string                  `json:"id" gorm:"column:id"`
	Name             string                  `json:"name" gorm:"column:name"`
	OrganizationID   string                  `json:"organizationId" gorm:"column:organization_id"`
	UserPermissions  []*UserPermissionValue  `json:"userPermissions" gorm:"-"`
	GroupPermissions []*GroupPermissionValue `json:"groupPermissions" gorm:"-"`
	Members          []string                `json:"members" gorm:"-"`
	CreateTime       string                  `json:"createTime" gorm:"column:create_time"`
	UpdateTime       *string                 `json:"updateTime" gorm:"column:update_time"`
}

func (*groupEntity) TableName() string {
	return "group"
}

func (g *groupEntity) BeforeCreate(*gorm.DB) (err error) {
	g.CreateTime = time.Now().UTC().Format(time.RFC3339)
	return nil
}

func (g *groupEntity) BeforeSave(*gorm.DB) (err error) {
	timeNow := time.Now().UTC().Format(time.RFC3339)
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
	var res = groupEntity{}
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

func (repo *groupRepo) GetIDsForFile(fileID string) ([]string, error) {
	type Value struct {
		Result string
	}
	var values []Value
	db := repo.db.
		Raw(`SELECT DISTINCT g.id as result FROM "group" g INNER JOIN grouppermission p ON p.resource_id = ? WHERE p.group_id = g.id`, fileID).
		Scan(&values)
	if db.Error != nil {
		return []string{}, db.Error
	}
	res := []string{}
	for _, v := range values {
		res = append(res, v.Result)
	}
	return res, nil
}

func (repo *groupRepo) GetIDsForUser(userID string) ([]string, error) {
	type Value struct {
		Result string
	}
	var values []Value
	db := repo.db.Raw(`SELECT group_id from group_user WHERE user_id = ?`, userID).Scan(&values)
	if db.Error != nil {
		return []string{}, db.Error
	}
	res := []string{}
	for _, v := range values {
		res = append(res, v.Result)
	}
	return res, nil
}

func (repo *groupRepo) GetIDsForOrganization(id string) ([]string, error) {
	type Value struct {
		Result string
	}
	var values []Value
	db := repo.db.Raw(`SELECT id as result from "group" WHERE organization_id = ?`, id).Scan(&values)
	if db.Error != nil {
		return []string{}, db.Error
	}
	res := []string{}
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

func (repo *groupRepo) AddUser(id string, userID string) error {
	db := repo.db.Exec("INSERT INTO group_user (group_id, user_id) VALUES (?, ?)", id, userID)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *groupRepo) RemoveMember(id string, userID string) error {
	db := repo.db.Exec("DELETE FROM group_user WHERE group_id = ? AND user_id = ?", id, userID)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *groupRepo) GetIDs() ([]string, error) {
	type Value struct {
		Result string
	}
	var values []Value
	db := repo.db.Raw(`SELECT id result FROM "group" ORDER BY create_time DESC`).Scan(&values)
	if db.Error != nil {
		return []string{}, db.Error
	}
	res := []string{}
	for _, v := range values {
		res = append(res, v.Result)
	}
	return res, nil
}

func (repo *groupRepo) GetMembers(id string) ([]model.User, error) {
	var entities []*userEntity
	db := repo.db.
		Raw(`SELECT DISTINCT u.* FROM "user" u INNER JOIN group_user gu ON u.id = gu.user_id WHERE gu.group_id = ?`, id).
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

func (repo *groupRepo) GrantUserPermission(id string, userID string, permission string) error {
	db := repo.db.Exec(
		"INSERT INTO userpermission (id, user_id, resource_id, permission) VALUES (?, ?, ?, ?) ON CONFLICT (user_id, resource_id) DO UPDATE SET permission = ?",
		helper.NewID(), userID, id, permission, permission)
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
		userPermissions, err := repo.permissionRepo.GetUserPermissions(g.ID)
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
		groupPermissions, err := repo.permissionRepo.GetGroupPermissions(g.ID)
		if err != nil {
			return err
		}
		for _, p := range groupPermissions {
			g.GroupPermissions = append(g.GroupPermissions, &GroupPermissionValue{
				GroupID: p.GetGroupID(),
				Value:   p.GetPermission(),
			})
		}
		members, err := repo.GetMembers(g.ID)
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
