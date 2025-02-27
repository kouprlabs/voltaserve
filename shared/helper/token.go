// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package helper

import (
	"github.com/kouprlabs/voltaserve/shared/dto"
	"time"
)

func NewExpiry(token *dto.Token) time.Time {
	return time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)
}
