package repo

import (
	"fmt"
	"voltaserve/config"
	"voltaserve/model"
)

type WorkspaceInsertOptions struct {
	ID              string
	Name            string
	StorageCapacity int64
	Image           *string
	OrganizationId  string
	RootId          string
	Bucket          string
}

type CoreWorkspaceRepo interface {
	Insert(opts WorkspaceInsertOptions) (model.CoreWorkspace, error)
	FindByName(name string) (model.CoreWorkspace, error)
	FindByID(id string) (model.CoreWorkspace, error)
	UpdateName(id string, name string) (model.CoreWorkspace, error)
	UpdateStorageCapacity(id string, storageCapacity int64) (model.CoreWorkspace, error)
	Delete(id string) error
	GetIDs() ([]string, error)
	GetIdsByOrganization(organizationId string) ([]string, error)
	UpdateRootID(id string, rootNodeId string) error
	GrantUserPermission(id string, userId string, permission string) error
}

func NewWorkspaceRepo() CoreWorkspaceRepo {
	if config.GetConfig().DatabaseType == config.DATABASE_TYPE_POSTGRES {
		return NewPostgresWorkspaceRepo()
	}
	panic(fmt.Sprintf("database type %s repo not implemented", config.GetConfig().DatabaseType))
}
