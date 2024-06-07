package model

type Organization interface {
	GetID() string
	GetName() string
	GetUserPermissions() []CoreUserPermission
	GetGroupPermissions() []CoreGroupPermission
	GetUsers() []string
	GetCreateTime() string
	GetUpdateTime() *string
	SetName(string)
}
