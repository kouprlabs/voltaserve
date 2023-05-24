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

type PostgresGroup struct {
	Id               string                   `json:"id"`
	Name             string                   `json:"name"`
	OrganizationId   string                   `json:"organizationId"`
	UserPermissions  []*model.UserPermission  `json:"userPermissions" gorm:"-"`
	GroupPermissions []*model.GroupPermission `json:"groupPermissions" gorm:"-"`
	Members          []string                 `json:"members" gorm:"-"`
	CreateTime       string                   `json:"createTime"`
	UpdateTime       *string                  `json:"updateTime"`
}

func (PostgresGroup) TableName() string {
	return "group"
}

func (g *PostgresGroup) BeforeCreate(tx *gorm.DB) (err error) {
	g.CreateTime = time.Now().UTC().Format(time.RFC3339)
	return nil
}

func (g *PostgresGroup) BeforeSave(tx *gorm.DB) (err error) {
	timeNow := time.Now().UTC().Format(time.RFC3339)
	g.UpdateTime = &timeNow
	return nil
}

func (g PostgresGroup) GetId() string {
	return g.Id
}

func (g PostgresGroup) GetName() string {
	return g.Name
}

func (g PostgresGroup) GetOrganizationId() string {
	return g.OrganizationId
}

func (g PostgresGroup) GetUserPermissions() []model.UserPermissionModel {
	var res []model.UserPermissionModel
	for _, p := range g.UserPermissions {
		res = append(res, p)
	}
	return res
}

func (g PostgresGroup) GetGroupPermissions() []model.GroupPermissionModel {
	var res []model.GroupPermissionModel
	for _, p := range g.GroupPermissions {
		res = append(res, p)
	}
	return res
}

func (g PostgresGroup) GetUsers() []string {
	return g.Members
}

func (g PostgresGroup) GetCreateTime() string {
	return g.CreateTime
}

func (g PostgresGroup) GetUpdateTime() *string {
	return g.UpdateTime
}

func (g *PostgresGroup) SetName(name string) {
	g.Name = name
}

func (g *PostgresGroup) SetUpdateTime(updateTime *string) {
	g.UpdateTime = updateTime
}

type PostgresGroupRepo struct {
	db             *gorm.DB
	permissionRepo *PostgresPermissionRepo
}

func NewPostgresGroupRepo() *PostgresGroupRepo {
	return &PostgresGroupRepo{
		db:             infra.GetDb(),
		permissionRepo: NewPostgresPermissionRepo(),
	}
}

func (repo *PostgresGroupRepo) Insert(opts GroupInsertOptions) (model.GroupModel, error) {
	group := PostgresGroup{
		Id:             opts.Id,
		Name:           opts.Name,
		OrganizationId: opts.OrganizationId,
	}
	if db := repo.db.Save(&group); db.Error != nil {
		return nil, db.Error
	}
	res, err := repo.Find(opts.Id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (repo *PostgresGroupRepo) find(id string) (*PostgresGroup, error) {
	var res = PostgresGroup{}
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

func (repo *PostgresGroupRepo) Find(id string) (model.GroupModel, error) {
	group, err := repo.find(id)
	if err != nil {
		return nil, err
	}
	if err := repo.populateModelFields([]*PostgresGroup{group}); err != nil {
		return nil, err
	}
	return group, nil
}

func (repo *PostgresGroupRepo) GetIdsForFile(fileId string) ([]string, error) {
	type Result struct {
		Result string
	}
	var results []Result
	db := repo.db.
		Raw(`SELECT DISTINCT g.id as result FROM "group" g INNER JOIN grouppermission p ON p.resource_id = ? WHERE p.group_id = g.id`, fileId).
		Scan(&results)
	if db.Error != nil {
		return []string{}, db.Error
	}
	res := []string{}
	for _, v := range results {
		res = append(res, v.Result)
	}
	return res, nil
}

func (repo *PostgresGroupRepo) GetIdsForUser(userId string) ([]string, error) {
	type Result struct {
		Result string
	}
	var results []Result
	db := repo.db.Raw(`SELECT id from group_user WHERE user_id = ?`, userId).Scan(&results)
	if db.Error != nil {
		return []string{}, db.Error
	}
	res := []string{}
	for _, v := range results {
		res = append(res, v.Result)
	}
	return res, nil
}

func (repo *PostgresGroupRepo) GetIdsForOrganization(id string) ([]string, error) {
	type Result struct {
		Result string
	}
	var results []Result
	db := repo.db.Raw(`SELECT id as result from "group" WHERE organization_id = ?`, id).Scan(&results)
	if db.Error != nil {
		return []string{}, db.Error
	}
	res := []string{}
	for _, v := range results {
		res = append(res, v.Result)
	}
	return res, nil
}

func (repo *PostgresGroupRepo) Save(group model.GroupModel) error {
	db := repo.db.Save(group)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *PostgresGroupRepo) Delete(id string) error {
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

func (repo *PostgresGroupRepo) AddUser(id string, userId string) error {
	db := repo.db.Exec("INSERT INTO group_user (group_id, user_id) VALUES (?, ?)", id, userId)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *PostgresGroupRepo) RemoveMember(id string, userId string) error {
	db := repo.db.Exec("DELETE FROM group_user WHERE group_id = ? AND user_id = ?", id, userId)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *PostgresGroupRepo) GetIDs() ([]string, error) {
	type Result struct {
		Result string
	}
	var results []Result
	db := repo.db.Raw(`SELECT id result FROM "group" ORDER BY create_time DESC`).Scan(&results)
	if db.Error != nil {
		return []string{}, db.Error
	}
	res := []string{}
	for _, v := range results {
		res = append(res, v.Result)
	}
	return res, nil
}

func (repo *PostgresGroupRepo) GetMembers(id string) ([]model.UserModel, error) {
	var entities []*PostgresUser
	db := repo.db.
		Raw(`SELECT DISTINCT u.* FROM "user" u INNER JOIN group_user gu ON u.id = gu.user_id WHERE gu.group_id = ?`, id).
		Scan(&entities)
	if db.Error != nil {
		return nil, db.Error
	}
	var res []model.UserModel
	for _, u := range entities {
		res = append(res, u)
	}
	return res, nil
}

func (repo *PostgresGroupRepo) GrantUserPermission(id string, userId string, permission string) error {
	db := repo.db.Exec(
		"INSERT INTO userpermission (id, user_id, resource_id, permission) "+
			"VALUES (?, ?, ?, ?) ON CONFLICT (user_id, resource_id) DO UPDATE SET permission = ?",
		helpers.NewId(), userId, id, permission, permission)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *PostgresGroupRepo) RevokeUserPermission(id string, userId string) error {
	db := repo.db.Exec("DELETE FROM userpermission WHERE user_id = ? AND resource_id = ?", userId, id)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *PostgresGroupRepo) populateModelFields(groups []*PostgresGroup) error {
	for _, g := range groups {
		g.UserPermissions = make([]*model.UserPermission, 0)
		userPermissions, err := repo.permissionRepo.GetUserPermissions(g.Id)
		if err != nil {
			return err
		}
		for _, p := range userPermissions {
			g.UserPermissions = append(g.UserPermissions, &model.UserPermission{
				UserId: p.UserId,
				Value:  p.Permission,
			})
		}
		g.GroupPermissions = make([]*model.GroupPermission, 0)
		groupPermissions, err := repo.permissionRepo.GetGroupPermissions(g.Id)
		if err != nil {
			return err
		}
		for _, p := range groupPermissions {
			g.GroupPermissions = append(g.GroupPermissions, &model.GroupPermission{
				GroupId: p.GroupId,
				Value:   p.Permission,
			})
		}
		members, err := repo.GetMembers(g.Id)
		if err != nil {
			return nil
		}
		g.Members = make([]string, 0)
		for _, u := range members {
			g.Members = append(g.Members, u.GetId())
		}
	}
	return nil
}
