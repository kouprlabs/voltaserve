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
	OrganizationSortByName         = "name"
	OrganizationSortByDateCreated  = "date_created"
	OrganizationSortByDateModified = "date_modified"
)

const (
	OrganizationSortOrderAsc  = "asc"
	OrganizationSortOrderDesc = "desc"
)

type Organization struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Image      *string `json:"image,omitempty"`
	Permission string  `json:"permission"`
	CreateTime string  `json:"createTime"`
	UpdateTime *string `json:"updateTime,omitempty"`
}

type OrganizationList struct {
	Data          []*Organization `json:"data"`
	TotalPages    uint64          `json:"totalPages"`
	TotalElements uint64          `json:"totalElements"`
	Page          uint64          `json:"page"`
	Size          uint64          `json:"size"`
}

type OrganizationProbe struct {
	TotalPages    uint64 `json:"totalPages"`
	TotalElements uint64 `json:"totalElements"`
}

type OrganizationCreateOptions struct {
	Name  string  `json:"name"  validate:"required,max=255"`
	Image *string `json:"image"`
}

type OrganizationPatchNameOptions struct {
	Name string `json:"name" validate:"required,max=255"`
}

type OrganizationRemoveMemberOptions struct {
	UserID string `json:"userId" validate:"required"`
}
