package model

const (
	FileTypeFile   = "file"
	FileTypeFolder = "folder"
)

type FileModel interface {
	GetID() string
	GetWorkspaceId() string
	GetName() string
	GetType() string
	GetParentID() *string
	GetCreateTime() string
	GetUpdateTime() *string
	GetSnapshots() []SnapshotModel
	GetUserPermissions() []UserPermissionModel
	GetGroupPermissions() []GroupPermissionModel
	GetText() *string
	SetID(string)
	SetParentID(*string)
	SetWorkspaceID(string)
	SetType(string)
	SetName(string)
	SetText(*string)
	SetCreateTime(string)
	SetUpdateTime(*string)
}
