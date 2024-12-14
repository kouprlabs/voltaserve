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
import { verify } from 'hono/jwt'
import { z } from 'zod'
import { zValidator } from '@hono/zod-validator'
import fs, { writeFile } from 'node:fs/promises'
import { tmpdir } from 'node:os'
import { getConfig } from '@/config/config.ts'
import {
  newInvalidJwtError,
  newMissingQueryParamError,
  newPictureNotFoundError,
} from '@/infra/error/creators.ts'
import {
  email,
  fullName,
  handleValidationError,
  password,
  token,
} from '@/infra/error/validation.ts'
import { checkAdmin } from '@/token/service.ts'
import {
  deletePicture,
  deleteUser,
  getUser,
  getUserByAdmin,
  getUserPicture,
  list,
  makeAdminUser,
  suspendUser,
  updateEmailConfirmation,
  updateEmailRequest,
  updateFullName,
  updatePassword,
  updatePicture,
  UserDeleteOptions,
  UserMakeAdminOptions,
  UserSuspendOptions,
  UserUpdateEmailConfirmationOptions,
  UserUpdateEmailRequestOptions,
  UserUpdateFullNameOptions,
  UserUpdatePasswordOptions,
} from '@/user/service.ts'
import { extname, join } from 'node:path'
import { Buffer } from 'node:buffer'
import { UserListOptions } from '@/user/service.ts'

const router = new Hono()

router.get('/me', async (c) => {
  c.json(await getUser(c.get('user').id))
})

router.get('/me/picture:extension', async (c) => {
  const accessToken = c.req.query('access_token')
  if (!accessToken) {
    throw newMissingQueryParamError('access_token')
  }
  const userId = await getUserIdFromAccessToken(accessToken)
  const { buffer, extension, mime } = await getUserPicture(userId)
  if (extension !== c.req.param('extension')) {
    throw newPictureNotFoundError()
  }
  c.res.headers.append(
    'Content-Disposition',
    `attachment; filename=picture.${extension}`,
  )
  c.res.headers.append('Content-Type', mime)
  return c.body(buffer)
})

router.post(
  '/me/update_full_name',
  zValidator(
    'json',
    z.object({ fullName }),
    handleValidationError,
  ),
  async (c) => {
    const body = c.req.valid('json') as UserUpdateFullNameOptions
    return c.json(await updateFullName(c.get('user').id, body))
  },
)

router.post(
  '/me/update_email_request',
  zValidator(
    'json',
    z.object({ email }),
    handleValidationError,
  ),
  async (c) => {
    const body = c.req.valid('json') as UserUpdateEmailRequestOptions
    return c.json(await updateEmailRequest(c.get('user').id, body))
  },
)

router.post(
  '/me/update_email_confirmation',
  zValidator(
    'json',
    z.object({ token }),
    handleValidationError,
  ),
  async (c) => {
    const body = c.req.valid('json') as UserUpdateEmailConfirmationOptions
    c.json(await updateEmailConfirmation(body))
  },
)

router.post(
  '/me/update_password',
  zValidator(
    'json',
    z.object({
      currentPassword: password,
      newPassword: password,
    }),
    handleValidationError,
  ),
  async (c) => {
    const body = c.req.valid('json') as UserUpdatePasswordOptions
    return c.json(await updatePassword(c.get('user').id, body))
  },
)

router.post(
  '/me/update_picture',
  zValidator(
    'form',
    z.object({
      file: z.instanceof(File)
        .refine((file) => file.size <= 3_000_000, 'File too large.')
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
      return c.json(await updatePicture(c.get('user').id, path, file.type))
    } finally {
      await fs.rm(path)
    }
  },
)

router.post('/me/delete_picture', async (c) => {
  return c.json(await deletePicture(c.get('user').id))
})

router.delete(
  '/me',
  zValidator(
    'json',
    z.object({ password }),
    handleValidationError,
  ),
  async (c) => {
    const body = c.req.valid('json') as UserDeleteOptions
    await deleteUser(c.get('user').id, body)
    c.status(204)
    return c.body(null)
  },
)

router.get(
  '/',
  zValidator(
    'query',
    z.object({
      query: z.string().optional(),
      page: z.string().regex(/^\d+$/, 'Must be a numeric value.').transform(
        Number,
      ),
      size: z.string().regex(/^\d+$/, 'Must be a numeric value.').transform(
        Number,
      ),
    }),
    handleValidationError,
  ),
  async (c) => {
    checkAdmin(c.req.header('Authorization')!)
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
    checkAdmin(c.req.header('Authorization')!)
    const { id } = c.req.param()
    const body = c.req.valid('json') as UserSuspendOptions
    await suspendUser(id, body)
    c.status(200)
    return c.body(null)
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
    checkAdmin(c.req.header('Authorization')!)
    const { id } = c.req.param()
    const body = c.req.valid('json') as UserMakeAdminOptions
    await makeAdminUser(id, body)
    c.status(200)
    return c.body(null)
  },
)

router.get('/:id', async (c) => {
  checkAdmin(c.req.header('Authorization')!)
  const { id } = c.req.param()
  return c.json(await getUserByAdmin(id))
})

async function getUserIdFromAccessToken(accessToken: string): Promise<string> {
  try {
    const payload = await verify(
      accessToken,
      getConfig().token.jwtSigningKey,
      'HS256',
    )
    if (payload.sub) {
      return payload.sub as string
    } else {
      throw newInvalidJwtError()
    }
  } catch {
    throw newInvalidJwtError()
  }
}

export default router
