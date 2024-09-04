// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { NextFunction, Router, Response } from 'express'
import { body, validationResult } from 'express-validator'
import fs from 'fs/promises'
import multer from 'multer'
import os from 'os'
import passport from 'passport'
import {
  SearchPaginatedRequest,
  UserAdminPostRequest,
  UserIdRequest,
  UserSuspendPostRequest,
} from '@/infra/admin-requests'
import { parseValidationError } from '@/infra/error'
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
  searchUserListPaginated,
} from './service'

const router = Router()

router.get(
  '/',
  passport.authenticate('jwt', { session: false }),
  async (req: PassportRequest, res: Response, next: NextFunction) => {
    try {
      res.json(await getUser(req.user.id))
    } catch (err) {
      next(err)
    }
  },
)

router.post(
  '/update_full_name',
  passport.authenticate('jwt', { session: false }),
  body('fullName').isString().notEmpty().trim().escape().isLength({ max: 255 }),
  async (req: PassportRequest, res: Response, next: NextFunction) => {
    try {
      const result = validationResult(req)
      if (!result.isEmpty()) {
        throw parseValidationError(result)
      }
      res.json(
        await updateFullName(
          req.user.id,
          req.body as UserUpdateFullNameOptions,
        ),
      )
    } catch (err) {
      next(err)
    }
  },
)

router.post(
  '/update_email_request',
  passport.authenticate('jwt', { session: false }),
  body('email').isEmail().isLength({ max: 255 }),
  async (req: PassportRequest, res: Response, next: NextFunction) => {
    try {
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
    } catch (err) {
      next(err)
    }
  },
)

router.post(
  '/update_email_confirmation',
  passport.authenticate('jwt', { session: false }),
  body('token').isString().notEmpty().trim(),
  async (req: PassportRequest, res: Response, next: NextFunction) => {
    try {
      const result = validationResult(req)
      if (!result.isEmpty()) {
        throw parseValidationError(result)
      }
      res.json(
        await updateEmailConfirmation(
          req.body as UserUpdateEmailConfirmationOptions,
        ),
      )
    } catch (err) {
      next(err)
    }
  },
)

router.post(
  '/update_password',
  passport.authenticate('jwt', { session: false }),
  body('currentPassword').notEmpty(),
  body('newPassword').isStrongPassword(),
  async (req: PassportRequest, res: Response, next: NextFunction) => {
    try {
      const result = validationResult(req)
      if (!result.isEmpty()) {
        throw parseValidationError(result)
      }
      res.json(
        await updatePassword(
          req.user.id,
          req.body as UserUpdatePasswordOptions,
        ),
      )
    } catch (err) {
      if (err === 'invalid_password') {
        res.status(400).json({ error: err })
        return
      } else {
        next(err)
      }
    }
  },
)

router.post(
  '/update_picture',
  passport.authenticate('jwt', { session: false }),
  multer({
    dest: os.tmpdir(),
    limits: { fileSize: 3000000, fields: 0, files: 1 },
  }).single('file'),
  async (req: PassportRequest, res: Response, next: NextFunction) => {
    try {
      const user = await updatePicture(
        req.user.id,
        req.file.path,
        req.file.mimetype,
      )
      await fs.rm(req.file.path)
      res.json(user)
    } catch (err) {
      next(err)
    }
  },
)

router.post(
  '/delete_picture',
  passport.authenticate('jwt', { session: false }),
  async (req: PassportRequest, res: Response, next: NextFunction) => {
    try {
      res.json(await deletePicture(req.user.id))
    } catch (err) {
      next(err)
    }
  },
)

router.delete(
  '/',
  passport.authenticate('jwt', { session: false }),
  body('password').isString().notEmpty(),
  async (req: PassportRequest, res: Response, next: NextFunction) => {
    try {
      const result = validationResult(req)
      if (!result.isEmpty()) {
        throw parseValidationError(result)
      }
      await deleteUser(req.user.id, req.body as UserDeleteOptions)
      res.sendStatus(204)
    } catch (err) {
      if (err === 'invalid_password') {
        res.status(400).json({ error: err })
        return
      } else {
        next(err)
      }
    }
  },
)

router.get(
  '/all',
  passport.authenticate('jwt', { session: false }),
  async (req: SearchPaginatedRequest, res: Response, next: NextFunction) => {
    try {
      checkAdmin(req.header('Authorization'))
      res.json(
        await searchUserListPaginated(
          req.query.query,
          parseInt(req.query.size),
          parseInt(req.query.page),
        ),
      )
    } catch (err) {
      next(err)
    }
  },
)

router.patch(
  '/suspend',
  passport.authenticate('jwt', { session: false }),
  body('id').isString(),
  body('suspend').isBoolean(),
  async (req: PassportRequest, res: Response, next: NextFunction) => {
    try {
      checkAdmin(req.header('Authorization'))
      const result = validationResult(req)
      if (!result.isEmpty()) {
        throw parseValidationError(result)
      }
      await suspendUser(req.body as UserSuspendPostRequest)
      res.sendStatus(200)
    } catch (err) {
      next(err)
    }
  },
)

router.patch(
  '/admin',
  passport.authenticate('jwt', { session: false }),
  body('id').isString(),
  body('makeAdmin').isBoolean(),
  async (req: PassportRequest, res: Response, next: NextFunction) => {
    try {
      checkAdmin(req.header('Authorization'))
      const result = validationResult(req)
      if (!result.isEmpty()) {
        throw parseValidationError(result)
      }
      await makeAdminUser(req.body as UserAdminPostRequest)
      res.sendStatus(200)
    } catch (err) {
      next(err)
    }
  },
)

router.get(
  '/:id',
  passport.authenticate('jwt', { session: false }),
  async (req: UserIdRequest, res: Response, next: NextFunction) => {
    checkAdmin(req.header('Authorization'))
    try {
      res.json(await getUserByAdmin(req.params.id))
    } catch (err) {
      next(err)
    }
  },
)

export default router
