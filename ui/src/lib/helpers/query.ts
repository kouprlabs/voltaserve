// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { encode, decode } from 'js-base64'
import { FileQuery } from '@/client/api/file'

export function encodeQuery(value: string) {
  return encode(value, true)
}

export function decodeQuery(value: string): string | undefined {
  if (!value) {
    return undefined
  }
  return decode(value)
}

export function encodeFileQuery(value: FileQuery) {
  return encode(JSON.stringify(value), true)
}

export function decodeFileQuery(value: string): FileQuery | undefined {
  if (!value) {
    return undefined
  }
  return JSON.parse(decode(value))
}
