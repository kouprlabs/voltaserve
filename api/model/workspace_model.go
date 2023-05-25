package model

type CoreWorkspace interface {
	GetID() string
	GetName() string
	GetStorageCapacity() int64
	GetRootID() string
	GetOrganizationID() string
	GetUserPermissions() []CoreUserPermission
	GetGroupPermissions() []CoreGroupPermission
	GetBucket() string
	GetCreateTime() string
	GetUpdateTime() *string
	SetName(string)
	SetUpdateTime(*string)
}
