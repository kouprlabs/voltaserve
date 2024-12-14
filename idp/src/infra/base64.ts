// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { Buffer } from 'node:buffer'

export function base64ToBuffer(value: string): Buffer | null {
  let withoutPrefix: string
  if (value.includes(',')) {
    withoutPrefix = value.split(',')[1]
  } else {
    withoutPrefix = value
  }
  try {
    return Buffer.from(withoutPrefix, 'base64')
  } catch {
    return null
  }
}

export function base64ToMIME(value: string): string | null {
  if (!value.startsWith('data:image/')) {
    return ''
  }
  const colonIndex = value.indexOf(':')
  const semicolonIndex = value.indexOf(';')
  if (colonIndex === -1 || semicolonIndex === -1) {
    return ''
  }
  return value.substring(colonIndex + 1, semicolonIndex)
}

export function base64ToExtension(value: string): string {
  const mime = base64ToMIME(value)
  switch (mime) {
    case 'image/jpeg':
      return '.jpg'
    case 'image/png':
      return '.png'
    default:
      return ''
  }
}
