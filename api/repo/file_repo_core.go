package repo

import "voltaserve/model"

type FileInsertOptions struct {
	Name        string
	WorkspaceId string
	ParentId    *string
	Type        string
}

type CoreFileRepo interface {
	New() model.FileModel
	Insert(opts FileInsertOptions) (model.FileModel, error)
	Find(id string) (model.FileModel, error)
	FindChildren(id string) ([]model.FileModel, error)
	FindPath(id string) ([]model.FileModel, error)
	FindTree(id string) ([]model.FileModel, error)
	GetIdsByWorkspace(workspaceId string) ([]string, error)
	AssignSnapshots(cloneId string, originalId string) error
	MoveSourceIntoTarget(targetId string, sourceId string) error
	Save(file model.FileModel) error
	BulkInsert(values []model.FileModel, chunkSize int) error
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
