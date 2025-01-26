// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package service

type Invitation struct {
	ID           string        `json:"id"`
	Owner        *User         `json:"owner,omitempty"`
	Email        string        `json:"email"`
	Organization *Organization `json:"organization,omitempty"`
	Status       string        `json:"status"`
	CreateTime   string        `json:"createTime"`
	UpdateTime   *string       `json:"updateTime"`
}

type InvitationList struct {
	Data          []*Invitation `json:"data"`
	TotalPages    uint64        `json:"totalPages"`
	TotalElements uint64        `json:"totalElements"`
	Page          uint64        `json:"page"`
	Size          uint64        `json:"size"`
}

type InvitationProbe struct {
	TotalPages    uint64 `json:"totalPages"`
	TotalElements uint64 `json:"totalElements"`
}
