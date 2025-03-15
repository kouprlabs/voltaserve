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
	UserSortByEmail        = "email"
	UserSortByFullName     = "full_name"
	UserSortByDateCreated  = "date_created"
	UserSortByDateModified = "date_modified"
)

const (
	UserSortOrderAsc  = "asc"
	UserSortOrderDesc = "desc"
)

type User struct {
	ID         string   `json:"id"`
	FullName   string   `json:"fullName"`
	Picture    *Picture `json:"picture,omitempty"`
	Email      string   `json:"email"`
	Username   string   `json:"username"`
	CreateTime string   `json:"createTime"`
	UpdateTime *string  `json:"updateTime,omitempty"`
}

type Picture struct {
	Extension string `json:"extension"`
}

type UserList struct {
	Data          []*User `json:"data"`
	TotalPages    uint64  `json:"totalPages"`
	TotalElements uint64  `json:"totalElements"`
	Page          uint64  `json:"page"`
	Size          uint64  `json:"size"`
}

type UserProbe struct {
	TotalPages    uint64 `json:"totalPages"`
	TotalElements uint64 `json:"totalElements"`
}
