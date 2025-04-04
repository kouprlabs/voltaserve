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
import { exchange, TokenExchangeOptions } from '@/token/service.ts'
import { handleValidationError } from '@/lib/validation.ts'

const router = new Hono()

router.post(
  '/',
  zValidator(
    'form',
    z.object({
      grant_type: z.union([
        z.literal('password'),
        z.literal('refresh_token'),
        z.literal('apple'),
      ]),
      username: z.string().optional(),
      password: z.string().optional(),
      refresh_token: z.string().optional(),
      apple_jwt: z.string().optional(),
      apple_full_name: z.string().optional(),
    }),
    handleValidationError,
  ),
  async (c) => {
    const options = c.req.valid('form') as TokenExchangeOptions
    return c.json(await exchange(options))
  },
)

export default router
