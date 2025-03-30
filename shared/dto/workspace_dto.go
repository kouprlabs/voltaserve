// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package dto

const (
	WorkspaceSortByName         = "name"
	WorkspaceSortByDateCreated  = "date_created"
	WorkspaceSortByDateModified = "date_modified"
)

const (
	WorkspaceSortOrderAsc  = "asc"
	WorkspaceSortOrderDesc = "desc"
)

type Workspace struct {
	ID              string       `json:"id"`
	Image           *string      `json:"image,omitempty"`
	Name            string       `json:"name"`
	RootID          string       `json:"rootId,omitempty"`
	StorageCapacity int64        `json:"storageCapacity"`
	Permission      string       `json:"permission"`
	Organization    Organization `json:"organization"`
	CreateTime      string       `json:"createTime"`
	UpdateTime      *string      `json:"updateTime,omitempty"`
}

type WorkspaceCreateOptions struct {
	Name            string  `json:"name"            validate:"required,max=255"`
	Image           *string `json:"image"`
	OrganizationID  string  `json:"organizationId"  validate:"required"`
	StorageCapacity int64   `json:"storageCapacity"`
}

type WorkspaceList struct {
	Data          []*Workspace `json:"data"`
	TotalPages    uint64       `json:"totalPages"`
	TotalElements uint64       `json:"totalElements"`
	Page          uint64       `json:"page"`
	Size          uint64       `json:"size"`
}

type WorkspaceProbe struct {
	TotalPages    uint64 `json:"totalPages"`
	TotalElements uint64 `json:"totalElements"`
}

type WorkspacePatchNameOptions struct {
	Name string `json:"name" validate:"required,max=255"`
}

type WorkspacePatchStorageCapacityOptions struct {
	StorageCapacity int64 `json:"storageCapacity" validate:"required,min=1"`
}

const (
	WorkspaceWebhookEventTypeCreate               = "create"
	WorkspaceWebhookEventTypePatchStorageCapacity = "patch_storage_capacity"
)

type WorkspaceWebhookOptions struct {
	EventType            string                                `json:"eventType"            validate:"required"`
	UserID               string                                `json:"userId"               validate:"required"`
	Create               *WorkspaceCreateOptions               `json:"create,omitempty"`
	PatchStorageCapacity *WorkspacePatchStorageCapacityOptions `json:"patchStorageCapacity"`
}
