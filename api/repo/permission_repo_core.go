package repo

type UserPermission struct {
	Id         string `json:"id"`
	UserId     string `json:"userId"`
	ResourceId string `json:"resourceId"`
	Permission string `json:"permission"`
	CreateTime string `json:"createTime"`
}

type GroupPermission struct {
	Id         string `json:"id"`
	GroupID    string `json:"groupId"`
	ResourceId string `json:"resourceId"`
	Permission string `json:"permission"`
	CreateTime string `json:"createTime"`
}

type CorePermissionRepo interface {
	GetUserPermissions(id string) ([]*UserPermission, error)
	GetGroupPermissions(id string) ([]*GroupPermission, error)
}
