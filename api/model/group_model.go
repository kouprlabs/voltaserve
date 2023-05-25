package model

type GroupModel interface {
	GetID() string
	GetName() string
	GetOrganizationID() string
	GetUserPermissions() []UserPermissionModel
	GetGroupPermissions() []GroupPermissionModel
	GetUsers() []string
	GetCreateTime() string
	GetUpdateTime() *string
	SetName(string)
	SetUpdateTime(*string)
}
