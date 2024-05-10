package model

const (
	FileTypeFile   = "file"
	FileTypeFolder = "folder"
)

type File interface {
	GetID() string
	GetWorkspaceID() string
	GetName() string
	GetType() string
	GetParentID() *string
	GetCreateTime() string
	GetUpdateTime() *string
	GetUserPermissions() []CoreUserPermission
	GetGroupPermissions() []CoreGroupPermission
	GetText() *string
	GetSnapshotID() *string
	SetID(string)
	SetParentID(*string)
	SetWorkspaceID(string)
	SetType(string)
	SetName(string)
	SetText(*string)
	SetSnapshotID(*string)
	SetCreateTime(string)
	SetUpdateTime(*string)
}
