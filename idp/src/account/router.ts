import { Router, Request, Response, NextFunction } from 'express'
import { body, validationResult } from 'express-validator'
import { parseValidationError } from '@/infra/error'
import {
  confirmEmail,
  createUser,
  resetPassword,
  sendResetPasswordEmail,
  AccountConfirmEmailOptions,
  AccountCreateOptions,
  AccountResetPasswordOptions,
  AccountSendResetPasswordEmailOptions,
} from './service'

const router = Router()

router.post(
  '/',
  body('email').isEmail().isLength({ max: 255 }),
  body('password').isStrongPassword().isLength({ max: 10000 }),
  body('fullName').isString().notEmpty().trim().escape().isLength({ max: 255 }),
  body('picture').optional().isBase64().isByteLength({ max: 3000000 }),
  async (req: Request, res: Response, next: NextFunction) => {
    try {
      const result = validationResult(req)
      if (!result.isEmpty()) {
        throw parseValidationError(result)
      }
      res.json(await createUser(req.body as AccountCreateOptions))
    } catch (err) {
      next(err)
    }
  },
)

router.post(
  '/reset_password',
  body('token').isString().notEmpty().trim(),
  body('newPassword').isStrongPassword(),
  async (req: Request, res: Response, next: NextFunction) => {
    try {
      const result = validationResult(req)
      if (!result.isEmpty()) {
        throw parseValidationError(result)
      }
      await resetPassword(req.body as AccountResetPasswordOptions)
      res.sendStatus(200)
    } catch (err) {
      next(err)
    }
  },
)

router.post(
  '/confirm_email',
  body('token').isString().notEmpty().trim(),
  async (req: Request, res: Response, next: NextFunction) => {
    try {
      const result = validationResult(req)
      if (!result.isEmpty()) {
        throw parseValidationError(result)
      }
      await confirmEmail(req.body as AccountConfirmEmailOptions)
      res.sendStatus(200)
    } catch (err) {
      next(err)
    }
  },
)

router.post(
  '/send_reset_password_email',
  body('email').isEmail().isLength({ max: 255 }),
  async (req: Request, res: Response, next: NextFunction) => {
    try {
      const result = validationResult(req)
      if (!result.isEmpty()) {
        throw parseValidationError(result)
      }
      res.json(
        await sendResetPasswordEmail(
          req.body as AccountSendResetPasswordEmailOptions,
        ),
      )
    } catch (err) {
      next(err)
    }
  },
)

export default router
