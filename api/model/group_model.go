package model

type Group interface {
	GetID() string
	GetName() string
	GetOrganizationID() string
	GetUserPermissions() []CoreUserPermission
	GetGroupPermissions() []CoreGroupPermission
	GetUsers() []string
	GetCreateTime() string
	GetUpdateTime() *string
	SetName(string)
	SetUpdateTime(*string)
}
