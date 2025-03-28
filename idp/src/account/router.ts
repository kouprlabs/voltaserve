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
import { handleValidationError, ZodFactory } from '@/lib/validation.ts'
import {
  AccountConfirmEmailOptions,
  AccountCreateOptions,
  AccountResetPasswordOptions,
  AccountSendResetPasswordEmailOptions,
  confirmEmail,
  createUser,
  getPasswordRequirements,
  resetPassword,
  sendResetPasswordEmail,
} from '@/account/service.ts'

const router = new Hono()

router.post(
  '/',
  zValidator(
    'json',
    z.object({
      email: ZodFactory.email(),
      password: ZodFactory.password(),
      fullName: ZodFactory.fullName(),
      picture: ZodFactory.picture(),
    }),
    handleValidationError,
  ),
  async (c) => {
    const body = c.req.valid('json') as AccountCreateOptions
    return c.json(await createUser(body))
  },
)

router.get('/password_requirements', (c) => {
  return c.json(getPasswordRequirements())
})

router.post(
  '/reset_password',
  zValidator(
    'json',
    z.object({
      token: ZodFactory.token(),
      newPassword: ZodFactory.password(),
    }),
    handleValidationError,
  ),
  async (c) => {
    const body = c.req.valid('json') as AccountResetPasswordOptions
    await resetPassword(body)
    return c.body(null, 200)
  },
)

router.post(
  '/confirm_email',
  zValidator(
    'json',
    z.object({ token: ZodFactory.token() }),
    handleValidationError,
  ),
  async (c) => {
    const body = c.req.valid('json') as AccountConfirmEmailOptions
    await confirmEmail(body)
    return c.body(null, 200)
  },
)

router.post(
  '/send_reset_password_email',
  zValidator(
    'json',
    z.object({ email: ZodFactory.email() }),
    handleValidationError,
  ),
  async (c) => {
    const body = c.req.valid('json') as AccountSendResetPasswordEmailOptions
    await sendResetPasswordEmail(body)
    return c.body(null, 204)
  },
)

export default router
