// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import Hashids from 'hashids'
import { v4 as uuidv4 } from 'uuid'

export function newHashId(): string {
  return new Hashids(uuidv4()).encode(Date.now())
}
