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
import fs, { writeFile } from 'node:fs/promises'
import { tmpdir } from 'node:os'
import { getConfig } from '@/config/config.ts'
import {
  newMissingQueryParamError,
  newPictureNotFoundError,
  newUserIsNotAdminError,
} from '@/error/creators.ts'
import { handleValidationError, ZodFactory } from '@/lib/validation.ts'
import {
  deletePicture,
  deleteUser,
  find,
  findAsAdmin,
  getPicture,
  list,
  makeAdmin,
  suspend,
  updateEmailConfirmation,
  updateEmailRequest,
  updateFullName,
  updatePassword,
  updatePicture,
  UserMakeAdminOptions,
  UserSuspendOptions,
  UserUpdateEmailConfirmationOptions,
  UserUpdateEmailRequestOptions,
  UserUpdateFullNameOptions,
  UserUpdatePasswordOptions,
} from '@/user/service.ts'
import { basename, extname, join } from 'node:path'
import { Buffer } from 'node:buffer'
import { UserListOptions } from '@/user/service.ts'
import { getUser, getUserIdFromAccessToken } from '@/lib/router.ts'

const router = new Hono()

router.get('/me', async (c) => {
  return c.json(await find(getUser(c).id))
})

router.get('/me/:filename', async (c) => {
  const { filename } = c.req.param()
  if (basename(filename, extname(filename)) !== 'picture') {
    return c.body(null, 404)
  }
  const accessToken = c.req.query('access_token') || c.req.query('access_key')
  if (!accessToken) {
    throw newMissingQueryParamError('access_token')
  }
  const userId = await getUserIdFromAccessToken(accessToken)
  const { buffer, extension, mime } = await getPicture(userId)
  if (extension !== extname(c.req.param('filename'))) {
    throw newPictureNotFoundError()
  }
  return c.body(buffer as any, 200, {
    'Content-Type': mime,
    'Content-Disposition': `attachment; filename=picture${extension}`,
  })
})

router.post(
  '/me/update_full_name',
  zValidator(
    'json',
    z.object({ fullName: ZodFactory.fullName() }),
    handleValidationError,
  ),
  async (c) => {
    const body = c.req.valid('json') as UserUpdateFullNameOptions
    return c.json(await updateFullName(getUser(c).id, body))
  },
)

router.post(
  '/me/update_email_request',
  zValidator(
    'json',
    z.object({ email: ZodFactory.email() }),
    handleValidationError,
  ),
  async (c) => {
    if (!getConfig().isLocalStrategy()) {
      return c.notFound()
    }
    const body = c.req.valid('json') as UserUpdateEmailRequestOptions
    return c.json(await updateEmailRequest(getUser(c).id, body))
  },
)

router.post(
  '/me/update_email_confirmation',
  zValidator(
    'json',
    z.object({ token: ZodFactory.token() }),
    handleValidationError,
  ),
  async (c) => {
    if (!getConfig().isLocalStrategy()) {
      return c.notFound()
    }
    const body = c.req.valid('json') as UserUpdateEmailConfirmationOptions
    return c.json(await updateEmailConfirmation(body))
  },
)

router.post(
  '/me/update_password',
  zValidator(
    'json',
    z.object({
      currentPassword: ZodFactory.password(),
      newPassword: ZodFactory.password(),
    }),
    handleValidationError,
  ),
  async (c) => {
    if (!getConfig().isLocalStrategy()) {
      return c.notFound()
    }
    const body = c.req.valid('json') as UserUpdatePasswordOptions
    return c.json(await updatePassword(getUser(c).id, body))
  },
)

router.post(
  '/me/update_picture',
  zValidator(
    'form',
    z.object({
      file: z.instanceof(File)
        .refine((file) => file.size <= 3 * 1024 * 1024, 'File too large.')
        .refine((file) =>
          [
            'image/jpeg',
            'image/png',
            'image/gif',
            'image/webp',
            'image/bmp',
            'image/tiff',
            'image/svg+xml',
            'image/x-icon',
          ].includes(file.type), 'File is not an image.'),
    }),
    handleValidationError,
  ),
  async (c) => {
    const { file }: { file: File } = c.req.valid('form')

    const path = join(tmpdir(), `${extname(file.name)}`)
    const arrayBuffer = await file.arrayBuffer()
    await writeFile(path, Buffer.from(arrayBuffer))

    try {
      return c.json(await updatePicture(getUser(c).id, path, file.type))
    } finally {
      await fs.rm(path)
    }
  },
)

router.post('/me/delete_picture', async (c) => {
  return c.json(await deletePicture(getUser(c).id))
})

router.delete(
  '/me',
  async (c) => {
    await deleteUser(getUser(c).id)
    return c.body(null, 204)
  },
)

router.get(
  '/',
  zValidator(
    'query',
    z.object({
      query: z.string().optional(),
      page: ZodFactory.page(),
      size: ZodFactory.size(),
    }),
    handleValidationError,
  ),
  async (c) => {
    if (!getConfig().isLocalStrategy()) {
      return c.notFound()
    }
    if (!c.get('user').isAdmin) {
      throw newUserIsNotAdminError()
    }
    const { query, size, page } = c.req.valid('query') as UserListOptions
    return c.json(await list({ query, size, page }))
  },
)

router.post(
  '/:id/suspend',
  zValidator(
    'json',
    z.object({ suspend: z.boolean() }),
    handleValidationError,
  ),
  async (c) => {
    if (!getConfig().isLocalStrategy()) {
      return c.notFound()
    }
    if (!c.get('user').isAdmin) {
      throw newUserIsNotAdminError()
    }
    const { id } = c.req.param()
    const body = c.req.valid('json') as UserSuspendOptions
    await suspend(id, body)
    return c.body(null, 200)
  },
)

router.post(
  '/:id/make_admin',
  zValidator(
    'json',
    z.object({ makeAdmin: z.boolean() }),
    handleValidationError,
  ),
  async (c) => {
    if (!getConfig().isLocalStrategy()) {
      return c.notFound()
    }
    if (!c.get('user').isAdmin) {
      throw newUserIsNotAdminError()
    }
    const { id } = c.req.param()
    const body = c.req.valid('json') as UserMakeAdminOptions
    await makeAdmin(id, body)
    return c.body(null, 200)
  },
)

router.get('/:id', async (c) => {
  if (!getConfig().isLocalStrategy()) {
    return c.notFound()
  }
  if (!c.get('user').isAdmin) {
    throw newUserIsNotAdminError()
  }
  const { id } = c.req.param()
  return c.json(await findAsAdmin(id))
})

export default router
