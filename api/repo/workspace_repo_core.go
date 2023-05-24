package repo

import "voltaserve/model"

type WorkspaceInsertOptions struct {
	Id              string
	Name            string
	StorageCapacity int64
	Image           *string
	OrganizationId  string
	RootId          string
	Bucket          string
}

type CoreWorkspaceRepo interface {
	Insert(opts WorkspaceInsertOptions) (model.WorkspaceModel, error)
	FindByName(name string) (model.WorkspaceModel, error)
	FindByID(id string) (model.WorkspaceModel, error)
	UpdateName(id string, name string) (model.WorkspaceModel, error)
	UpdateStorageCapacity(id string, storageCapacity int64) (model.WorkspaceModel, error)
	Delete(id string) error
	GetIDs() ([]string, error)
	GetIdsByOrganization(organizationId string) ([]string, error)
	UpdateRootID(id string, rootNodeId string) error
	GrantUserPermission(id string, userId string, permission string) error
}
