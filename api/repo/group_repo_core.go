package repo

import "voltaserve/model"

type GroupInsertOptions struct {
	ID             string
	Name           string
	OrganizationId string
	OwnerId        string
}

type CoreGroupRepo interface {
	Insert(opts GroupInsertOptions) (model.CoreGroup, error)
	Find(id string) (model.CoreGroup, error)
	GetIDsForFile(fileId string) ([]string, error)
	GetIDsForUser(userId string) ([]string, error)
	GetIDsForOrganization(id string) ([]string, error)
	Save(group model.CoreGroup) error
	Delete(id string) error
	AddUser(id string, userId string) error
	RemoveMember(id string, userId string) error
	GetIDs() ([]string, error)
	GetMembers(id string) ([]model.CoreUser, error)
	GrantUserPermission(id string, userId string, permission string) error
	RevokeUserPermission(id string, userId string) error
}
