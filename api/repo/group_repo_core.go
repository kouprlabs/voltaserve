package repo

import "voltaserve/model"

type GroupInsertOptions struct {
	Id             string
	Name           string
	OrganizationId string
	OwnerId        string
}

type CoreGroupRepo interface {
	Insert(opts GroupInsertOptions) (model.GroupModel, error)
	Find(id string) (model.GroupModel, error)
	GetIdsForFile(fileId string) ([]string, error)
	GetIdsForUser(userId string) ([]string, error)
	GetIdsForOrganization(id string) ([]string, error)
	Save(group model.GroupModel) error
	Delete(id string) error
	AddUser(id string, userId string) error
	RemoveMember(id string, userId string) error
	GetIDs() ([]string, error)
	GetMembers(id string) ([]model.UserModel, error)
	GrantUserPermission(id string, userId string, permission string) error
	RevokeUserPermission(id string, userId string) error
}
