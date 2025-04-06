// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

import { User } from '@/user/model.ts'
import { Context } from 'hono'
import { verify } from 'hono/jwt'
import { newInvalidJwtError, newUserNotFoundError } from '@/error/creators.ts'
import { getConfig } from '@/config/config.ts'

export function getUser(c: Context): User {
  const user = c.get('user')
  if (!user) {
    throw newUserNotFoundError()
  }
  return user
}

export async function getUserIdFromAccessToken(
  accessToken: string,
): Promise<string> {
  try {
    const payload = await verify(
      accessToken,
      getConfig().token.jwtSigningKey,
      'HS256',
    )
    if (payload.sub) {
      return payload.sub as string
    } else {
      // noinspection ExceptionCaughtLocallyJS
      throw newInvalidJwtError()
    }
  } catch {
    throw newInvalidJwtError()
  }
}
