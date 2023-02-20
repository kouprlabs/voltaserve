package model

type GroupModel interface {
	GetId() string
	GetName() string
	GetOrganizationId() string
	GetUserPermissions() []UserPermissionModel
	GetGroupPermissions() []GroupPermissionModel
	GetUsers() []string
	GetCreateTime() string
	GetUpdateTime() *string
	SetName(string)
	SetUpdateTime(*string)
}
