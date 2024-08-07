// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { scryptSync, randomBytes } from 'crypto'

export function hashPassword(password: string): string {
  const salt = randomBytes(16).toString('hex')
  const key = scryptSync(password, salt, 64).toString('hex')
  return `${key}:${salt}`
}

export function verifyPassword(password: string, hash: string): boolean {
  const [key, salt] = hash.split(':')
  const newKey = scryptSync(password, salt, 64).toString('hex')
  return newKey === key
}
