package model

const (
	PermissionViewer = "viewer"
	PermissionEditor = "editor"
	PermissionOwner  = "owner"
)

type CoreUserPermission interface {
	GetUserID() string
	GetValue() string
}

type CoreGroupPermission interface {
	GetGroupID() string
	GetValue() string
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
