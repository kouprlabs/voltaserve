// Copyright (c) 2024 Mateusz Kaźmierczak.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package webhook

import "github.com/kouprlabs/voltaserve/shared/dto"

type UserWebhook struct{}

func NewUserWebhook() *UserWebhook {
	return &UserWebhook{}
}

func (wh *UserWebhook) Handle(opts dto.UserWebhookOptions) error {
	if opts.EventType == dto.UserWebhookEventTypeCreate {
		return wh.handleCreate(opts)
	} else if opts.EventType == dto.UserWebhookEventTypeDelete {
		return wh.handleDelete(opts)
	}
	return nil
}

func (wh *UserWebhook) handleCreate(opts dto.UserWebhookOptions) error {
	return nil
}

func (wh *UserWebhook) handleDelete(opts dto.UserWebhookOptions) error {
	return nil
}
