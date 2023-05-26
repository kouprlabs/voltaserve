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
	GetSnapshots() []Snapshot
	GetUserPermissions() []CoreUserPermission
	GetGroupPermissions() []CoreGroupPermission
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
