package repo

import "voltaserve/model"

type OrganizationInsertOptions struct {
	Id   string
	Name string
}

type CoreOrganizationRepo interface {
	Insert(opts OrganizationInsertOptions) (model.OrganizationModel, error)
	Find(id string) (model.OrganizationModel, error)
	Save(org model.OrganizationModel) error
	Delete(id string) error
	GetIDs() ([]string, error)
	AddUser(id string, userId string) error
	RemoveMember(id string, userId string) error
	GetMembers(id string) ([]model.UserModel, error)
	GetGroups(id string) ([]model.GroupModel, error)
	GetOwnerCount(id string) (int64, error)
	GrantUserPermission(id string, userId string, permission string) error
	RevokeUserPermission(id string, userId string) error
}
