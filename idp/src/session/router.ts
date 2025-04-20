// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { Hono } from 'hono'
import { z } from 'zod'
import { zValidator } from '@hono/zod-validator'
import { exchange, TokenGrantType } from '@/token/service.ts'
import { handleValidationError } from '@/lib/validation.ts'

const router = new Hono()

type SessionExchangeOptions = {
  grant_type: TokenGrantType
  username?: string
  password?: string
  refresh_key?: string
  apple_key?: string
  apple_full_name?: string
}

type Session = {
  access_key: string
  key_type: string
  expires_in: number
  refresh_key: string
}

router.post(
  '/',
  zValidator(
    'form',
    z.object({
      grant_type: z.union([
        z.literal('password'),
        z.literal('refresh_key'),
        z.literal('apple'),
      ]),
      username: z.string().optional(),
      password: z.string().optional(),
      refresh_token: z.string().optional(),
      apple_key: z.string().optional(),
      apple_full_name: z.string().optional(),
    }),
    handleValidationError,
  ),
  async (c) => {
    const options = c.req.valid('form') as SessionExchangeOptions
    const token = await exchange({
      grant_type: options.grant_type,
      username: options.username,
      password: options.password,
      refresh_token: options.refresh_key,
      apple_token: options.apple_key,
      apple_full_name: options.apple_full_name,
    })
    return c.json({
      access_key: token.access_token,
      key_type: token.token_type,
      expires_in: token.expires_in,
      refresh_key: token.refresh_token,
    } as Session)
  },
)

export default router
