// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

import { newInternalServerError } from '@/error/creators.ts'

const APPLE_KEYS_URL = 'https://appleid.apple.com/auth/keys'

let cachedKeys: Record<string, CryptoKey> | null = null

export async function getApplePublicKey(header: any): Promise<CryptoKey> {
  if (cachedKeys && cachedKeys[header.kid]) {
    return cachedKeys[header.kid]
  }
  const res = await fetch(APPLE_KEYS_URL)
  const { keys } = await res.json()
  cachedKeys = {}
  for (const key of keys) {
    const jwk = {
      alg: 'RS256',
      ...key,
    }
    cachedKeys[key.kid] = await crypto.subtle.importKey(
      'jwk',
      jwk,
      { name: 'RSASSA-PKCS1-v1_5', hash: 'SHA-256' },
      true,
      ['verify'],
    )
  }
  if (!cachedKeys[header.kid]) {
    throw newInternalServerError('Apple public key not found.')
  }
  return cachedKeys[header.kid]
}
