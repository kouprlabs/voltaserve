package repo

import "voltaserve/model"

type FileInsertOptions struct {
	Name        string
	WorkspaceId string
	ParentId    *string
	Type        string
}

type CoreFileRepo interface {
	New() model.CoreFile
	Insert(opts FileInsertOptions) (model.CoreFile, error)
	Find(id string) (model.CoreFile, error)
	FindChildren(id string) ([]model.CoreFile, error)
	FindPath(id string) ([]model.CoreFile, error)
	FindTree(id string) ([]model.CoreFile, error)
	GetIdsByWorkspace(workspaceId string) ([]string, error)
	AssignSnapshots(cloneId string, originalId string) error
	MoveSourceIntoTarget(targetId string, sourceId string) error
	Save(file model.CoreFile) error
	BulkInsert(values []model.CoreFile, chunkSize int) error
	BulkInsertPermissions(values []*UserPermission, chunkSize int) error
	Delete(id string) error
	GetChildrenIDs(id string) ([]string, error)
	GetItemCount(id string) (int64, error)
	IsGrandChildOf(id string, ancestorId string) (bool, error)
	GetSize(id string) (int64, error)
	GrantUserPermission(id string, userId string, permission string) error
	RevokeUserPermission(id string, userId string) error
	GrantGroupPermission(id string, groupId string, permission string) error
	RevokeGroupPermission(id string, groupId string) error
}

func NewFileRepo() CoreFileRepo {
	return NewPostgresFileRepo()
}
