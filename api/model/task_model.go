// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

package model

const (
	TaskStatusWaiting = "waiting"
	TaskStatusRunning = "running"
)

type Task interface {
	GetID() string
	GetName() string
	GetError() *string
	GetPercentage() *int
	GetIsIndeterminate() bool
	GetUserID() string
	GetStatus() string
	GetPayload() map[string]string
	HasError() bool
	SetName(string)
	SetError(*string)
	SetPercentage(*int)
	SetIsIndeterminate(bool)
	SetUserID(string)
	SetStatus(string)
	SetPayload(map[string]string)
	GetCreateTime() string
	GetUpdateTime() *string
}
