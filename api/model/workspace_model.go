package model

type WorkspaceModel interface {
	GetId() string
	GetName() string
	GetStorageCapacity() int64
	GetRootId() string
	GetOrganizationId() string
	GetUserPermissions() []UserPermissionModel
	GetGroupPermissions() []GroupPermissionModel
	GetBucket() string
	GetCreateTime() string
	GetUpdateTime() *string
	SetName(string)
	SetUpdateTime(*string)
}
