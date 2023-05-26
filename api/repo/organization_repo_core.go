package repo

import (
	"fmt"
	"voltaserve/config"
	"voltaserve/model"
)

type OrganizationInsertOptions struct {
	ID   string
	Name string
}

type CoreOrganizationRepo interface {
	Insert(opts OrganizationInsertOptions) (model.CoreOrganization, error)
	Find(id string) (model.CoreOrganization, error)
	Save(org model.CoreOrganization) error
	Delete(id string) error
	GetIDs() ([]string, error)
	AddUser(id string, userId string) error
	RemoveMember(id string, userId string) error
	GetMembers(id string) ([]model.CoreUser, error)
	GetGroups(id string) ([]model.CoreGroup, error)
	GetOwnerCount(id string) (int64, error)
	GrantUserPermission(id string, userId string, permission string) error
	RevokeUserPermission(id string, userId string) error
}

func NewOrganizationRepo() CoreOrganizationRepo {
	if config.GetConfig().DatabaseType == config.DATABASE_TYPE_POSTGRES {
		return NewPostgresOrganizationRepo()
	}
	panic(fmt.Sprintf("database type %s repo not implemented", config.GetConfig().DatabaseType))
}
