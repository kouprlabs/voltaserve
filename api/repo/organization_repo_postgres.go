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

type OrganizationEntity struct {
	Id               string                   `json:"id"`
	Name             string                   `json:"name"`
	UserPermissions  []*model.UserPermission  `json:"userPermissions" gorm:"-"`
	GroupPermissions []*model.GroupPermission `json:"groupPermissions" gorm:"-"`
	Members          []string                 `json:"members" gorm:"-"`
	CreateTime       string                   `json:"createTime"`
	UpdateTime       *string                  `json:"updateTime,omitempty"`
}

func (OrganizationEntity) TableName() string {
	return "organization"
}

func (o *OrganizationEntity) BeforeCreate(tx *gorm.DB) (err error) {
	o.CreateTime = time.Now().UTC().Format(time.RFC3339)
	return nil
}

func (o *OrganizationEntity) BeforeSave(tx *gorm.DB) (err error) {
	timeNow := time.Now().UTC().Format(time.RFC3339)
	o.UpdateTime = &timeNow
	return nil
}

func (o OrganizationEntity) GetID() string {
	return o.Id
}

func (o OrganizationEntity) GetName() string {
	return o.Name
}

func (o OrganizationEntity) GetUserPermissions() []model.UserPermissionModel {
	var res []model.UserPermissionModel
	for _, p := range o.UserPermissions {
		res = append(res, p)
	}
	return res
}

func (o OrganizationEntity) GetGroupPermissions() []model.GroupPermissionModel {
	var res []model.GroupPermissionModel
	for _, p := range o.GroupPermissions {
		res = append(res, p)
	}
	return res
}

func (o OrganizationEntity) GetUsers() []string {
	return o.Members
}

func (o OrganizationEntity) GetCreateTime() string {
	return o.CreateTime
}

func (o OrganizationEntity) GetUpdateTime() *string {
	return o.UpdateTime
}

func (w *OrganizationEntity) SetName(name string) {
	w.Name = name
}

func (w *OrganizationEntity) SetUpdateTime(updateTime *string) {
	w.UpdateTime = updateTime
}

type PostgresOrganizationRepo struct {
	db             *gorm.DB
	groupRepo      *PostgresGroupRepo
	permissionRepo *PostgresPermissionRepo
}

func NewPostgresOrganizationRepo() *PostgresOrganizationRepo {
	return &PostgresOrganizationRepo{
		db:             infra.GetDb(),
		groupRepo:      NewPostgresGroupRepo(),
		permissionRepo: NewPostgresPermissionRepo(),
	}
}

func (repo *PostgresOrganizationRepo) Insert(opts OrganizationInsertOptions) (model.OrganizationModel, error) {
	org := OrganizationEntity{
		Id:   opts.Id,
		Name: opts.Name,
	}
	if db := repo.db.Save(&org); db.Error != nil {
		return nil, db.Error
	}
	res, err := repo.Find(opts.Id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (repo *PostgresOrganizationRepo) find(id string) (*OrganizationEntity, error) {
	var res = OrganizationEntity{}
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

func (repo *PostgresOrganizationRepo) Find(id string) (model.OrganizationModel, error) {
	org, err := repo.find(id)
	if err != nil {
		return nil, err
	}
	if err := repo.populateModelFields([]*OrganizationEntity{org}); err != nil {
		return nil, err
	}
	return org, nil
}

func (repo *PostgresOrganizationRepo) Save(org model.OrganizationModel) error {
	db := repo.db.Save(org)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *PostgresOrganizationRepo) Delete(id string) error {
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

func (repo *PostgresOrganizationRepo) GetIDs() ([]string, error) {
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

func (repo *PostgresOrganizationRepo) AddUser(id string, userId string) error {
	db := repo.db.Exec("INSERT INTO organization_user (organization_id, user_id) VALUES (?, ?)", id, userId)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *PostgresOrganizationRepo) RemoveMember(id string, userId string) error {
	db := repo.db.Exec("DELETE FROM organization_user WHERE organization_id = ? AND user_id = ?", id, userId)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *PostgresOrganizationRepo) GetMembers(id string) ([]model.UserModel, error) {
	var entities []*PostgresUser
	db := repo.db.
		Raw(`SELECT DISTINCT u.* FROM "user" u INNER JOIN organization_user ou ON u.id = ou.user_id WHERE ou.organization_id = ? ORDER BY u.full_name ASC`, id).
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

func (repo *PostgresOrganizationRepo) GetGroups(id string) ([]model.GroupModel, error) {
	var entities []*PostgresGroup
	db := repo.db.
		Raw(`SELECT * FROM "group" g WHERE g.organization_id = ? ORDER BY g.name ASC`, id).
		Scan(&entities)
	if db.Error != nil {
		return nil, db.Error
	}
	if err := repo.groupRepo.populateModelFields(entities); err != nil {
		return nil, err
	}
	var res []model.GroupModel
	for _, g := range entities {
		res = append(res, g)
	}
	return res, nil
}

func (repo *PostgresOrganizationRepo) GetOwnerCount(id string) (int64, error) {
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

func (repo *PostgresOrganizationRepo) GrantUserPermission(id string, userId string, permission string) error {
	db := repo.db.Exec(
		"INSERT INTO userpermission (id, user_id, resource_id, permission) "+
			"VALUES (?, ?, ?, ?) ON CONFLICT (user_id, resource_id) DO UPDATE SET permission = ?",
		helpers.NewId(), userId, id, permission, permission)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *PostgresOrganizationRepo) RevokeUserPermission(id string, userId string) error {
	db := repo.db.Exec("DELETE FROM userpermission WHERE user_id = ? AND resource_id = ?", userId, id)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *PostgresOrganizationRepo) populateModelFields(organizations []*OrganizationEntity) error {
	for _, o := range organizations {
		o.UserPermissions = make([]*model.UserPermission, 0)
		userPermissions, err := repo.permissionRepo.GetUserPermissions(o.Id)
		if err != nil {
			return err
		}
		for _, p := range userPermissions {
			o.UserPermissions = append(o.UserPermissions, &model.UserPermission{
				UserId: p.UserId,
				Value:  p.Permission,
			})
		}
		o.GroupPermissions = make([]*model.GroupPermission, 0)
		groupPermissions, err := repo.permissionRepo.GetGroupPermissions(o.Id)
		if err != nil {
			return err
		}
		for _, p := range groupPermissions {
			o.GroupPermissions = append(o.GroupPermissions, &model.GroupPermission{
				GroupId: p.GroupId,
				Value:   p.Permission,
			})
		}
		members, err := repo.GetMembers(o.Id)
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
