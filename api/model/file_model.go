package model

const (
	FileTypeFile   = "file"
	FileTypeFolder = "folder"
)

type FileModel interface {
	GetId() string
	GetWorkspaceId() string
	GetName() string
	GetType() string
	GetParentId() *string
	GetCreateTime() string
	GetUpdateTime() *string
	GetSnapshots() []SnapshotModel
	GetUserPermissions() []UserPermissionModel
	GetGroupPermissions() []GroupPermissionModel
	GetText() *string
	SetId(string)
	SetParentId(*string)
	SetWorkspaceId(string)
	SetType(string)
	SetName(string)
	SetText(*string)
	SetCreateTime(string)
	SetUpdateTime(*string)
}
