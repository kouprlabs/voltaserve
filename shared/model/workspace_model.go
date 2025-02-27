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

type Workspace interface {
	GetID() string
	GetName() string
	GetStorageCapacity() int64
	GetRootID() string
	GetOrganizationID() string
	GetUserPermissions() []CoreUserPermission
	GetGroupPermissions() []CoreGroupPermission
	GetBucket() string
	GetCreateTime() string
	GetUpdateTime() *string
	SetID(string)
	SetName(string)
	SetStorageCapacity(int64)
	SetRootID(string)
	SetOrganizationID(string)
	SetUserPermissions([]CoreUserPermission)
	SetGroupPermissions([]CoreGroupPermission)
	SetBucket(string)
	SetCreateTime(string)
	SetUpdateTime(*string)
}
