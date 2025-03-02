// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package model

const (
	FileTypeFile   = "file"
	FileTypeFolder = "folder"
)

type File interface {
	GetID() string
	GetParentID() *string
	GetWorkspaceID() string
	GetName() string
	GetType() string
	GetUserPermissions() []CoreUserPermission
	GetGroupPermissions() []CoreGroupPermission
	GetText() *string
	GetSummary() *string
	GetSnapshotID() *string
	GetCreateTime() string
	GetUpdateTime() *string
	SetID(string)
	SetParentID(*string)
	SetWorkspaceID(string)
	SetType(string)
	SetName(string)
	SetText(*string)
	SetSummary(*string)
	SetSnapshotID(*string)
	SetUserPermissions([]CoreUserPermission)
	SetGroupPermissions([]CoreGroupPermission)
	SetCreateTime(string)
	SetUpdateTime(*string)
}
