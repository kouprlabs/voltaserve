// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

package model

const (
	PermissionNone   = "none"
	PermissionViewer = "viewer"
	PermissionEditor = "editor"
	PermissionOwner  = "owner"
)

type UserPermission interface {
	GetID() string
	GetUserID() string
	GetResourceID() string
	GetPermission() string
	GetCreateTime() string
	SetID(string)
	SetUserID(string)
	SetResourceID(string)
	SetPermission(string)
	SetCreateTime(string)
}

type GroupPermission interface {
	GetID() string
	GetGroupID() string
	GetResourceID() string
	GetPermission() string
	GetCreateTime() string
	SetID(string)
	SetGroupID(string)
	SetResourceID(string)
	SetPermission(string)
	SetCreateTime(string)
}

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
