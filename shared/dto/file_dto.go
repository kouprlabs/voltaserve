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

import (
	"slices"
)

const (
	FileSortByName         = "name"
	FileSortByKind         = "kind"
	FileSortBySize         = "size"
	FileSortByDateCreated  = "date_created"
	FileSortByDateModified = "date_modified"
)

const (
	FileSortOrderAsc  = "asc"
	FileSortOrderDesc = "desc"
)

type File struct {
	ID          string    `json:"id"`
	WorkspaceID string    `json:"workspaceId"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	ParentID    *string   `json:"parentId,omitempty"`
	Permission  string    `json:"permission"`
	IsShared    *bool     `json:"isShared,omitempty"`
	Snapshot    *Snapshot `json:"snapshot,omitempty"`
	CreateTime  string    `json:"createTime"`
	UpdateTime  *string   `json:"updateTime,omitempty"`
}

type FileQuery struct {
	Text             *string `json:"text"                       validate:"required"`
	Type             *string `json:"type,omitempty"             validate:"omitempty,oneof=file folder"`
	CreateTimeAfter  *int64  `json:"createTimeAfter,omitempty"`
	CreateTimeBefore *int64  `json:"createTimeBefore,omitempty"`
	UpdateTimeAfter  *int64  `json:"updateTimeAfter,omitempty"`
	UpdateTimeBefore *int64  `json:"updateTimeBefore,omitempty"`
}

type FileList struct {
	Data          []*File    `json:"data"`
	TotalPages    uint64     `json:"totalPages"`
	TotalElements uint64     `json:"totalElements"`
	Page          uint64     `json:"page"`
	Size          uint64     `json:"size"`
	Query         *FileQuery `json:"query,omitempty"`
}

type FileProbe struct {
	TotalPages    uint64 `json:"totalPages"`
	TotalElements uint64 `json:"totalElements"`
}

type FilePatchNameOptions struct {
	Name string `json:"name" validate:"required,max=255"`
}

type FileCopyManyOptions struct {
	SourceIDs []string `json:"sourceIds" validate:"required"`
	TargetID  string   `json:"targetId"  validate:"required"`
}

type FileCopyManyResult struct {
	New       []string `json:"new"`
	Succeeded []string `json:"succeeded"`
	Failed    []string `json:"failed"`
}

type FileMoveManyOptions struct {
	SourceIDs []string `json:"sourceIds" validate:"required"`
	TargetID  string   `json:"targetId"  validate:"required"`
}

type FileMoveManyResult struct {
	Succeeded []string `json:"succeeded"`
	Failed    []string `json:"failed"`
}

type FileDeleteManyOptions struct {
	IDs []string `json:"ids" validate:"required"`
}

type FileDeleteManyResult struct {
	Succeeded []string `json:"succeeded"`
	Failed    []string `json:"failed"`
}

type FileGrantUserPermissionOptions struct {
	UserID     string   `json:"userId"     validate:"required"`
	IDs        []string `json:"ids"        validate:"required"`
	Permission string   `json:"permission" validate:"required,oneof=viewer editor owner"`
}

type FileRevokeUserPermissionOptions struct {
	IDs    []string `json:"ids"    validate:"required"`
	UserID string   `json:"userId" validate:"required"`
}

type FileGrantGroupPermissionOptions struct {
	GroupID    string   `json:"groupId"    validate:"required"`
	IDs        []string `json:"ids"        validate:"required"`
	Permission string   `json:"permission" validate:"required,oneof=viewer editor owner"`
}

type FileRevokeGroupPermissionOptions struct {
	IDs     []string `json:"ids"     validate:"required"`
	GroupID string   `json:"groupId" validate:"required"`
}

type UserPermission struct {
	ID         string `json:"id"`
	User       *User  `json:"user"`
	Permission string `json:"permission"`
}

type GroupPermission struct {
	ID         string `json:"id"`
	Group      *Group `json:"group"`
	Permission string `json:"permission"`
}

type FileReprocessResult struct {
	Accepted []string `json:"accepted"`
	Rejected []string `json:"rejected"`
}

func (r *FileReprocessResult) AppendAccepted(id string) {
	if !slices.Contains(r.Accepted, id) {
		r.Accepted = append(r.Accepted, id)
	}
}

func (r *FileReprocessResult) AppendRejected(id string) {
	if !slices.Contains(r.Rejected, id) {
		r.Rejected = append(r.Rejected, id)
	}
}
