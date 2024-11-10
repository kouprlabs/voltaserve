// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { Router, Response } from 'express'
import { body, query, validationResult } from 'express-validator'
import fs from 'fs/promises'
import { jwtVerify } from 'jose'
import multer from 'multer'
import os from 'os'
import passport from 'passport'
import { getConfig } from '@/config/config'
import { newInvalidJwtError, newUserPictureNotFoundError } from '@/infra/error'
import { ErrorCode, newError, parseValidationError } from '@/infra/error/core'
import { PassportRequest } from '@/infra/passport-request'
import { checkAdmin } from '@/token/service'
import {
  deleteUser,
  getUser,
  updateFullName,
  updatePicture,
  updatePassword,
  UserDeleteOptions,
  UserUpdateFullNameOptions,
  UserUpdatePasswordOptions,
  deletePicture,
  UserUpdateEmailRequestOptions,
  UserUpdateEmailConfirmationOptions,
  updateEmailRequest,
  updateEmailConfirmation,
  suspendUser,
  makeAdminUser,
  getUserByAdmin,
  list,
  getUserPicture,
  UserMakeAdminOptions,
  UserSuspendOptions,
} from './service'

const router = Router()

router.get(
  '/me',
  passport.authenticate('jwt', { session: false }),
  async (req: PassportRequest, res: Response) => {
    res.json(await getUser(req.user.id))
  },
)

router.get(
  '/me/picture:extension',
  async (req: PassportRequest, res: Response) => {
    if (!req.query.access_token) {
      throw newError({
        code: ErrorCode.InvalidRequest,
        message: "Query param 'access_token' is required",
      })
    }
    const userId = await getUserIdFromAccessToken(
      req.query.access_token as string,
    )
    const { buffer, extension, mime } = await getUserPicture(userId)
    if (extension !== req.params.extension) {
      throw newUserPictureNotFoundError(userId)
    }
    res.setHeader(
      'Content-Disposition',
      `attachment; filename=picture.${extension}`,
    )
    res.setHeader('Content-Type', mime)
    res.send(buffer)
  },
)

router.post(
  '/me/update_full_name',
  passport.authenticate('jwt', { session: false }),
  body('fullName').isString().notEmpty().trim().escape().isLength({ max: 255 }),
  async (req: PassportRequest, res: Response) => {
    const result = validationResult(req)
    if (!result.isEmpty()) {
      throw parseValidationError(result)
    }
    res.json(
      await updateFullName(req.user.id, req.body as UserUpdateFullNameOptions),
    )
  },
)

router.post(
  '/me/update_email_request',
  passport.authenticate('jwt', { session: false }),
  body('email').isEmail().isLength({ max: 255 }),
  async (req: PassportRequest, res: Response) => {
    const result = validationResult(req)
    if (!result.isEmpty()) {
      throw parseValidationError(result)
    }
    res.json(
      await updateEmailRequest(
        req.user.id,
        req.body as UserUpdateEmailRequestOptions,
      ),
    )
  },
)

router.post(
  '/me/update_email_confirmation',
  passport.authenticate('jwt', { session: false }),
  body('token').isString().notEmpty().trim(),
  async (req: PassportRequest, res: Response) => {
    const result = validationResult(req)
    if (!result.isEmpty()) {
      throw parseValidationError(result)
    }
    res.json(
      await updateEmailConfirmation(
        req.body as UserUpdateEmailConfirmationOptions,
      ),
    )
  },
)

router.post(
  '/me/update_password',
  passport.authenticate('jwt', { session: false }),
  body('currentPassword').notEmpty(),
  body('newPassword').isStrongPassword(),
  async (req: PassportRequest, res: Response) => {
    const result = validationResult(req)
    if (!result.isEmpty()) {
      throw parseValidationError(result)
    }
    res.json(
      await updatePassword(req.user.id, req.body as UserUpdatePasswordOptions),
    )
  },
)

router.post(
  '/me/update_picture',
  passport.authenticate('jwt', { session: false }),
  multer({
    dest: os.tmpdir(),
    limits: { fileSize: 3000000, fields: 0, files: 1 },
  }).single('file'),
  async (req: PassportRequest, res: Response) => {
    const user = await updatePicture(
      req.user.id,
      req.file.path,
      req.file.mimetype,
    )
    await fs.rm(req.file.path)
    res.json(user)
  },
)

router.post(
  '/me/delete_picture',
  passport.authenticate('jwt', { session: false }),
  async (req: PassportRequest, res: Response) => {
    res.json(await deletePicture(req.user.id))
  },
)

router.delete(
  '/me',
  passport.authenticate('jwt', { session: false }),
  body('password').isString().notEmpty(),
  async (req: PassportRequest, res: Response) => {
    const result = validationResult(req)
    if (!result.isEmpty()) {
      throw parseValidationError(result)
    }
    await deleteUser(req.user.id, req.body as UserDeleteOptions)
    res.sendStatus(204)
  },
)

router.get(
  '/',
  passport.authenticate('jwt', { session: false }),
  query('query').isString().optional(),
  query('page').isInt(),
  query('size').isInt(),
  async (req: PassportRequest, res: Response) => {
    checkAdmin(req.header('Authorization'))
    const result = validationResult(req)
    if (!result.isEmpty()) {
      throw parseValidationError(result)
    }
    res.json(
      await list({
        query: req.query.query as string,
        size: parseInt(req.query.size as string),
        page: parseInt(req.query.page as string),
      }),
    )
  },
)

router.post(
  '/:id/suspend',
  passport.authenticate('jwt', { session: false }),
  body('suspend').isBoolean(),
  async (req: PassportRequest, res: Response) => {
    checkAdmin(req.header('Authorization'))
    const result = validationResult(req)
    if (!result.isEmpty()) {
      throw parseValidationError(result)
    }
    await suspendUser(req.params.id, req.body as UserSuspendOptions)
    res.sendStatus(200)
  },
)

router.post(
  '/:id/make_admin',
  passport.authenticate('jwt', { session: false }),
  body('makeAdmin').isBoolean(),
  async (req: PassportRequest, res: Response) => {
    checkAdmin(req.header('Authorization'))
    const result = validationResult(req)
    if (!result.isEmpty()) {
      throw parseValidationError(result)
    }
    await makeAdminUser(req.params.id, req.body as UserMakeAdminOptions)
    res.sendStatus(200)
  },
)

router.get(
  '/:id',
  passport.authenticate('jwt', { session: false }),
  async (req: PassportRequest, res: Response) => {
    checkAdmin(req.header('Authorization'))
    res.json(await getUserByAdmin(req.params.id))
  },
)

async function getUserIdFromAccessToken(accessToken: string): Promise<string> {
  try {
    const { payload } = await jwtVerify(
      accessToken,
      new TextEncoder().encode(getConfig().token.jwtSigningKey),
    )
    return payload.sub
  } catch {
    throw newInvalidJwtError()
  }
}

export default router
