// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

import hashids from 'hashids'
import { v4 as uuidv4 } from 'uuid'

export function newHashId(): string {
  return new hashids(uuidv4()).encode(Date.now())
}

export function newHyphenlessUuid(): string {
  return uuidv4().replaceAll('-', '')
}
