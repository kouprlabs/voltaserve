package model

type WorkspaceModel interface {
	GetID() string
	GetName() string
	GetStorageCapacity() int64
	GetRootID() string
	GetOrganizationID() string
	GetUserPermissions() []UserPermissionModel
	GetGroupPermissions() []GroupPermissionModel
	GetBucket() string
	GetCreateTime() string
	GetUpdateTime() *string
	SetName(string)
	SetUpdateTime(*string)
}
