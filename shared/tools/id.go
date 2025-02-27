// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package tools

import (
	"time"

	"github.com/google/uuid"
	"github.com/speps/go-hashids/v2"
)

func NewID() string {
	hd := hashids.NewData()
	hd.Salt = uuid.NewString()
	h, err := hashids.NewWithData(hd)
	if err != nil {
		panic(err)
	}
	id, err := h.EncodeInt64([]int64{time.Now().UTC().UnixNano()})
	if err != nil {
		panic(err)
	}
	return id
}
