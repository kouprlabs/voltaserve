package model

type CoreOrganization interface {
	GetID() string
	GetName() string
	GetUserPermissions() []CoreUserPermission
	GetGroupPermissions() []CoreGroupPermission
	GetUsers() []string
	GetCreateTime() string
	GetUpdateTime() *string
	SetName(string)
	SetUpdateTime(*string)
}
