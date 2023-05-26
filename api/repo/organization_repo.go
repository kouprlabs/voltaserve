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

type OrganizationInsertOptions struct {
	ID   string
	Name string
}

type OrganizationRepo interface {
	Insert(opts OrganizationInsertOptions) (model.Organization, error)
	Find(id string) (model.Organization, error)
	Save(org model.Organization) error
	Delete(id string) error
	GetIDs() ([]string, error)
	AddUser(id string, userId string) error
	RemoveMember(id string, userId string) error
	GetMembers(id string) ([]model.User, error)
	GetGroups(id string) ([]model.Group, error)
	GetOwnerCount(id string) (int64, error)
	GrantUserPermission(id string, userId string, permission string) error
	RevokeUserPermission(id string, userId string) error
}

func NewOrganizationRepo() OrganizationRepo {
	return newOrganizationRepo()
}

func NewOrganization() model.Organization {
	return &organizationEntity{}
}

type organizationEntity struct {
	ID               string                  `json:"id"`
	Name             string                  `json:"name"`
	UserPermissions  []*userPermissionValue  `json:"userPermissions" gorm:"-"`
	GroupPermissions []*groupPermissionValue `json:"groupPermissions" gorm:"-"`
	Members          []string                `json:"members" gorm:"-"`
	CreateTime       string                  `json:"createTime"`
	UpdateTime       *string                 `json:"updateTime,omitempty"`
}

func (organizationEntity) TableName() string {
	return "organization"
}

func (o *organizationEntity) BeforeCreate(tx *gorm.DB) (err error) {
	o.CreateTime = time.Now().UTC().Format(time.RFC3339)
	return nil
}

func (o *organizationEntity) BeforeSave(tx *gorm.DB) (err error) {
	timeNow := time.Now().UTC().Format(time.RFC3339)
	o.UpdateTime = &timeNow
	return nil
}

func (o organizationEntity) GetID() string {
	return o.ID
}

func (o organizationEntity) GetName() string {
	return o.Name
}

func (o organizationEntity) GetUserPermissions() []model.CoreUserPermission {
	var res []model.CoreUserPermission
	for _, p := range o.UserPermissions {
		res = append(res, p)
	}
	return res
}

func (o organizationEntity) GetGroupPermissions() []model.CoreGroupPermission {
	var res []model.CoreGroupPermission
	for _, p := range o.GroupPermissions {
		res = append(res, p)
	}
	return res
}

func (o organizationEntity) GetUsers() []string {
	return o.Members
}

func (o organizationEntity) GetCreateTime() string {
	return o.CreateTime
}

func (o organizationEntity) GetUpdateTime() *string {
	return o.UpdateTime
}

func (w *organizationEntity) SetName(name string) {
	w.Name = name
}

func (w *organizationEntity) SetUpdateTime(updateTime *string) {
	w.UpdateTime = updateTime
}

type organizationRepo struct {
	db             *gorm.DB
	groupRepo      *groupRepo
	permissionRepo *permissionRepo
}

func newOrganizationRepo() *organizationRepo {
	return &organizationRepo{
		db:             infra.GetDb(),
		groupRepo:      newGroupRepo(),
		permissionRepo: newPermissionRepo(),
	}
}

func (repo *organizationRepo) Insert(opts OrganizationInsertOptions) (model.Organization, error) {
	org := organizationEntity{
		ID:   opts.ID,
		Name: opts.Name,
	}
	if db := repo.db.Save(&org); db.Error != nil {
		return nil, db.Error
	}
	res, err := repo.Find(opts.ID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (repo *organizationRepo) find(id string) (*organizationEntity, error) {
	var res = organizationEntity{}
	db := repo.db.Where("id = ?", id).First(&res)
	if db.Error != nil {
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			return nil, errorpkg.NewOrganizationNotFoundError(db.Error)
		} else {
			return nil, errorpkg.NewInternalServerError(db.Error)
		}
	}
	return &res, nil
}

func (repo *organizationRepo) Find(id string) (model.Organization, error) {
	org, err := repo.find(id)
	if err != nil {
		return nil, err
	}
	if err := repo.populateModelFields([]*organizationEntity{org}); err != nil {
		return nil, err
	}
	return org, nil
}

func (repo *organizationRepo) Save(org model.Organization) error {
	db := repo.db.Save(org)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *organizationRepo) Delete(id string) error {
	db := repo.db.Exec("DELETE FROM organization WHERE id = ?", id)
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

func (repo *organizationRepo) GetIDs() ([]string, error) {
	type Result struct {
		Result string
	}
	var results []Result
	db := repo.db.Raw("SELECT id result FROM organization ORDER BY create_time DESC").Scan(&results)
	if db.Error != nil {
		return []string{}, db.Error
	}
	res := []string{}
	for _, v := range results {
		res = append(res, v.Result)
	}
	return res, nil
}

func (repo *organizationRepo) AddUser(id string, userId string) error {
	db := repo.db.Exec("INSERT INTO organization_user (organization_id, user_id) VALUES (?, ?)", id, userId)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *organizationRepo) RemoveMember(id string, userId string) error {
	db := repo.db.Exec("DELETE FROM organization_user WHERE organization_id = ? AND user_id = ?", id, userId)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *organizationRepo) GetMembers(id string) ([]model.User, error) {
	var entities []*postgresUser
	db := repo.db.
		Raw(`SELECT DISTINCT u.* FROM "user" u INNER JOIN organization_user ou ON u.id = ou.user_id WHERE ou.organization_id = ? ORDER BY u.full_name ASC`, id).
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

func (repo *organizationRepo) GetGroups(id string) ([]model.Group, error) {
	var entities []*groupEntity
	db := repo.db.
		Raw(`SELECT * FROM "group" g WHERE g.organization_id = ? ORDER BY g.name ASC`, id).
		Scan(&entities)
	if db.Error != nil {
		return nil, db.Error
	}
	if err := repo.groupRepo.populateModelFields(entities); err != nil {
		return nil, err
	}
	var res []model.Group
	for _, g := range entities {
		res = append(res, g)
	}
	return res, nil
}

func (repo *organizationRepo) GetOwnerCount(id string) (int64, error) {
	type Result struct {
		Result int64
	}
	var res Result
	db := repo.db.
		Raw("SELECT count(*) as result FROM userpermission WHERE resource_id = ? and permission = ?", id, model.PermissionOwner).
		Scan(&res)
	if db.Error != nil {
		return 0, db.Error
	}
	return res.Result, nil
}

func (repo *organizationRepo) GrantUserPermission(id string, userId string, permission string) error {
	db := repo.db.Exec(
		"INSERT INTO userpermission (id, user_id, resource_id, permission) "+
			"VALUES (?, ?, ?, ?) ON CONFLICT (user_id, resource_id) DO UPDATE SET permission = ?",
		helpers.NewId(), userId, id, permission, permission)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *organizationRepo) RevokeUserPermission(id string, userId string) error {
	db := repo.db.Exec("DELETE FROM userpermission WHERE user_id = ? AND resource_id = ?", userId, id)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *organizationRepo) populateModelFields(organizations []*organizationEntity) error {
	for _, o := range organizations {
		o.UserPermissions = make([]*userPermissionValue, 0)
		userPermissions, err := repo.permissionRepo.GetUserPermissions(o.ID)
		if err != nil {
			return err
		}
		for _, p := range userPermissions {
			o.UserPermissions = append(o.UserPermissions, &userPermissionValue{
				UserId: p.UserID,
				Value:  p.Permission,
			})
		}
		o.GroupPermissions = make([]*groupPermissionValue, 0)
		groupPermissions, err := repo.permissionRepo.GetGroupPermissions(o.ID)
		if err != nil {
			return err
		}
		for _, p := range groupPermissions {
			o.GroupPermissions = append(o.GroupPermissions, &groupPermissionValue{
				GroupID: p.GroupID,
				Value:   p.Permission,
			})
		}
		members, err := repo.GetMembers(o.ID)
		if err != nil {
			return nil
		}
		o.Members = make([]string, 0)
		for _, u := range members {
			o.Members = append(o.Members, u.GetID())
		}
	}
	return nil
}
