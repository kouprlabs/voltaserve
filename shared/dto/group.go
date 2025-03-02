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
	GroupSortByName         = "name"
	GroupSortByDateCreated  = "date_created"
	GroupSortByDateModified = "date_modified"
)

const (
	GroupSortOrderAsc  = "asc"
	GroupSortOrderDesc = "desc"
)

type Group struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Image        *string      `json:"image,omitempty"`
	Organization Organization `json:"organization"`
	Permission   string       `json:"permission"`
	CreateTime   string       `json:"createTime,omitempty"`
	UpdateTime   *string      `json:"updateTime"`
}

type GroupCreateOptions struct {
	Name           string  `json:"name"           validate:"required,max=255"`
	Image          *string `json:"image"`
	OrganizationID string  `json:"organizationId" validate:"required"`
}

type GroupPatchNameOptions struct {
	Name string `json:"name" validate:"required,max=255"`
}

type GroupAddMemberOptions struct {
	UserID string `json:"userId" validate:"required"`
}

type GroupRemoveMemberOptions struct {
	UserID string `json:"userId" validate:"required"`
}
type GroupList struct {
	Data          []*Group `json:"data"`
	TotalPages    uint64   `json:"totalPages"`
	TotalElements uint64   `json:"totalElements"`
	Page          uint64   `json:"page"`
	Size          uint64   `json:"size"`
}

type GroupProbe struct {
	TotalPages    uint64 `json:"totalPages"`
	TotalElements uint64 `json:"totalElements"`
}
