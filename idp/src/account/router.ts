// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { Request, Response, Router } from 'express'
import { body, validationResult } from 'express-validator'
import { getConfig } from '@/config/config.ts'
import { parseValidationError } from '@/infra/error/index.ts'
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
} from './service.ts'

const router = Router()

router.post(
  '/',
  body('email').isEmail().isLength({ max: 255 }),
  body('password')
    .isStrongPassword({
      minLength: getConfig().password.minLength,
      minLowercase: getConfig().password.minLowercase,
      minUppercase: getConfig().password.minUppercase,
      minNumbers: getConfig().password.minNumbers,
      minSymbols: getConfig().password.minSymbols,
    })
    .isLength({ max: 10000 }),
  body('fullName').isString().notEmpty().trim().escape().isLength({ max: 255 }),
  body('picture').optional().isBase64().isByteLength({ max: 3000000 }),
  async (req: Request, res: Response) => {
    const result = validationResult(req)
    if (!result.isEmpty()) {
      throw parseValidationError(result)
    }
    res.json(await createUser(req.body as AccountCreateOptions))
  },
)

router.get('/password_requirements', (_: Request, res: Response) => {
  res.json(getPasswordRequirements())
})

router.post(
  '/reset_password',
  body('token').isString().notEmpty().trim(),
  body('newPassword').isStrongPassword(),
  async (req: Request, res: Response) => {
    const result = validationResult(req)
    if (!result.isEmpty()) {
      throw parseValidationError(result)
    }
    await resetPassword(req.body as AccountResetPasswordOptions)
    res.sendStatus(200)
  },
)

router.post(
  '/confirm_email',
  body('token').isString().notEmpty().trim(),
  async (req: Request, res: Response) => {
    const result = validationResult(req)
    if (!result.isEmpty()) {
      throw parseValidationError(result)
    }
    await confirmEmail(req.body as AccountConfirmEmailOptions)
    res.sendStatus(200)
  },
)

router.post(
  '/send_reset_password_email',
  body('email').isEmail().isLength({ max: 255 }),
  async (req: Request, res: Response) => {
    const result = validationResult(req)
    if (!result.isEmpty()) {
      throw parseValidationError(result)
    }
    res.json(
      await sendResetPasswordEmail(
        req.body as AccountSendResetPasswordEmailOptions,
      ),
    )
  },
)

export default router
