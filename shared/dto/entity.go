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
	EntitySortByName      = "name"
	EntitySortByFrequency = "frequency"
)

const (
	EntitySortOrderAsc  = "asc"
	EntitySortOrderDesc = "desc"
)

type Entity struct {
	Text      string `json:"text"`
	Label     string `json:"label"`
	Frequency int    `json:"frequency"`
}

type EntityCreateOptions struct {
	Language string `json:"language" validate:"required"`
}

type EntityListOptions struct {
	Query     string `json:"query"`
	Page      uint64 `json:"page"`
	Size      uint64 `json:"size"`
	SortBy    string `json:"sortBy"`
	SortOrder string `json:"sortOrder"`
}

type EntityList struct {
	Data          []*Entity `json:"data"`
	TotalPages    uint64    `json:"totalPages"`
	TotalElements uint64    `json:"totalElements"`
	Page          uint64    `json:"page"`
	Size          uint64    `json:"size"`
}

type EntityProbe struct {
	TotalPages    uint64 `json:"totalPages"`
	TotalElements uint64 `json:"totalElements"`
}
