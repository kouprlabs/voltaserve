package model

const (
	PermissionViewer = "viewer"
	PermissionEditor = "editor"
	PermissionOwner  = "owner"
)

type UserPermissionModel interface {
	GetUserID() string
	GetValue() string
}

type UserPermission struct {
	UserId string `json:"userId,omitempty"`
	Value  string `json:"value,omitempty"`
}

func (p UserPermission) GetUserID() string {
	return p.UserId
}
func (p UserPermission) GetValue() string {
	return p.Value
}

type GroupPermissionModel interface {
	GetGroupID() string
	GetValue() string
}

type GroupPermission struct {
	GroupID string `json:"groupId,omitempty"`
	Value   string `json:"value,omitempty"`
}

func (p GroupPermission) GetGroupID() string {
	return p.GroupID
}
func (p GroupPermission) GetValue() string {
	return p.Value
}

func GteViewerPermission(permission string) bool {
	return permission == PermissionViewer || permission == PermissionEditor || permission == PermissionOwner
}

func GteEditorPermission(permission string) bool {
	return permission == PermissionEditor || permission == PermissionOwner
}

func GteOwnerPermission(permission string) bool {
	return permission == PermissionOwner
}

func IsEquivalentPermission(permission string, otherPermission string) bool {
	if permission == otherPermission {
		return true
	}
	if otherPermission == PermissionViewer && GteViewerPermission(permission) {
		return true
	}
	if otherPermission == PermissionEditor && GteEditorPermission(permission) {
		return true
	}
	if otherPermission == PermissionOwner && GteOwnerPermission(permission) {
		return true
	}
	return false
}

func GetPermissionWeight(permission string) int {
	if permission == PermissionViewer {
		return 1
	}
	if permission == PermissionEditor {
		return 2
	}
	if permission == PermissionOwner {
		return 3
	}
	return 0
}
