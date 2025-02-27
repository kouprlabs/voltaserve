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
	InvitationStatusPending  = "pending"
	InvitationStatusAccepted = "accepted"
	InvitationStatusDeclined = "declined"
)

type Invitation interface {
	GetID() string
	GetOrganizationID() string
	GetOwnerID() string
	GetEmail() string
	GetStatus() string
	GetCreateTime() string
	GetUpdateTime() *string
	SetID(string)
	SetOrganizationID(string)
	SetOwnerID(string)
	SetEmail(string)
	SetStatus(string)
	SetCreateTime(string)
	SetUpdateTime(*string)
}
