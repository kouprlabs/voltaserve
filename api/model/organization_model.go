package model

type OrganizationModel interface {
	GetID() string
	GetName() string
	GetUserPermissions() []UserPermissionModel
	GetGroupPermissions() []GroupPermissionModel
	GetUsers() []string
	GetCreateTime() string
	GetUpdateTime() *string
	SetName(string)
	SetUpdateTime(*string)
}
