// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import fs from 'node:fs/promises'
import { getConfig } from '@/config/config.ts'
import {
  base64ToBuffer,
  base64ToExtension,
  base64ToMIME,
} from '@/infra/base64.ts'
import {
  newCannotDemoteSoleAdminError,
  newCannotSuspendSoleAdminError,
  newInternalServerError,
  newInvalidPasswordError,
  newPasswordValidationFailedError,
  newPictureNotFoundError,
  newUsernameUnavailableError,
  newUserNotFoundError,
} from '@/infra/error/index.ts'
import { ErrorCode, newError } from '@/infra/error/core.ts'
import { newHyphenlessUuid } from '@/infra/id.ts'
import { sendTemplateMail } from '@/infra/mail.ts'
import { hashPassword, verifyPassword } from '@/infra/password.ts'
import search, { USER_SEARCH_INDEX } from '@/infra/search.ts'
import { User } from '@/user/model.ts'
import userRepo from '@/user/repo.ts'
import { Buffer } from 'node:buffer'

export type UserDTO = {
  id: string
  fullName: string
  picture?: PictureDTO
  email: string
  username: string
  pendingEmail?: string
}

export type UserAdminDTO = {
  id: string
  fullName: string
  username: string
  email: string
  passwordHash?: string
  refreshTokenValue?: string
  refreshTokenExpiry?: string
  resetPasswordToken?: string
  emailConfirmationToken?: string
  isEmailConfirmed: boolean
  isAdmin: boolean
  isActive: boolean
  emailUpdateToken?: string
  emailUpdateValue?: string
  picture?: PictureDTO
  createTime: string
  updateTime?: string
}

export interface UserAdminList {
  data: UserAdminDTO[]
  page: number
  size: number
  totalElements: number
  totalPages: number
}

export type PictureDTO = {
  extension: string
}

export type UserUpdateEmailRequestOptions = {
  email: string
}

export type UserUpdateEmailConfirmationOptions = {
  token: string
}

export type UserUpdateFullNameOptions = {
  fullName: string
}

export type UserUpdatePasswordOptions = {
  currentPassword: string
  newPassword: string
}

export type UserDeleteOptions = {
  password: string
}

export type UserPictureResponse = {
  buffer: Buffer
  extension: string
  mime: string
}

export type UserSuspendOptions = {
  suspend: boolean
}

export type UserMakeAdminOptions = {
  makeAdmin: boolean
}

export type UserListOptions = {
  query?: string
  size: number
  page: number
}

export async function getUser(id: string): Promise<UserDTO> {
  return mapEntity(await userRepo.findById(id))
}

export async function getUserByAdmin(id: string): Promise<UserAdminDTO> {
  return adminMapEntity(await userRepo.findById(id))
}

export async function getUserPicture(id: string): Promise<UserPictureResponse> {
  const user = await userRepo.findById(id)
  if (!user.picture) {
    throw newPictureNotFoundError()
  }
  const buffer = base64ToBuffer(user.picture)
  if (!buffer) {
    throw newPictureNotFoundError()
  }
  const extension = base64ToExtension(user.picture)
  if (!extension) {
    throw newPictureNotFoundError()
  }
  const mime = base64ToMIME(user.picture)
  if (!mime) {
    throw newPictureNotFoundError()
  }
  return { buffer, extension, mime }
}

export async function list({
  query,
  size,
  page,
}: UserListOptions): Promise<UserAdminList> {
  if (query && query.length >= 3) {
    const users = await search
      .index(USER_SEARCH_INDEX)
      .search(query, { page: page, hitsPerPage: size })
      .then((value) => {
        return {
          data: value.hits,
          totalElements: value.totalHits,
        }
      })
    return {
      data: (
        await userRepo.findMany(
          users.data.map((value) => {
            return value.id
          }),
        )
      ).map((value) => adminMapEntity(value)),
      totalElements: users.totalElements,
      totalPages: Math.floor((users.totalElements + size - 1) / size),
      size: size,
      page: page,
    }
  } else {
    return {
      data: (await userRepo.list(page, size)).map((value) =>
        adminMapEntity(value)
      ),
      totalElements: await userRepo.getCount(),
      totalPages: Math.floor(((await userRepo.getCount()) + size - 1) / size),
      size: size,
      page: page,
    }
  }
}

export async function getUserCount(): Promise<number> {
  return await userRepo.getCount()
}

export async function updateFullName(
  id: string,
  options: UserUpdateFullNameOptions,
): Promise<UserDTO> {
  let user = await userRepo.findById(id)
  user = await userRepo.update({ id: user.id, fullName: options.fullName })
  await search.index(USER_SEARCH_INDEX).updateDocuments([
    {
      id: user.id,
      username: user.username,
      email: user.email,
      fullName: user.fullName,
      isEmailConfirmed: user.isEmailConfirmed,
      createTime: user.createTime,
      updateTime: user.updateTime,
      picture: user.picture,
    },
  ])
  return mapEntity(user)
}

export async function updateEmailRequest(
  id: string,
  options: UserUpdateEmailRequestOptions,
): Promise<UserDTO> {
  let user = await userRepo.findById(id)
  if (options.email === user.email) {
    user = await userRepo.update({
      id: user.id,
      emailUpdateToken: null,
      emailUpdateValue: null,
    })
    return mapEntity(user)
  } else {
    let usernameUnavailable = false
    try {
      await userRepo.findByUsername(options.email)
      usernameUnavailable = true
    } catch {
      // Ignored
    }
    if (usernameUnavailable) {
      throw newUsernameUnavailableError()
    }
    user = await userRepo.update({
      id: user.id,
      emailUpdateToken: newHyphenlessUuid(),
      emailUpdateValue: options.email,
    })
    try {
      await sendTemplateMail('email-update', options.email, {
        'EMAIL': options.email,
        'UI_URL': getConfig().publicUIURL,
        'TOKEN': user.emailUpdateToken,
      })
      return mapEntity(user)
    } catch (error) {
      await userRepo.update({
        id,
        emailUpdateToken: null,
        emailUpdateValue: null,
      })
      throw newInternalServerError(error)
    }
  }
}

export async function updateEmailConfirmation(
  options: UserUpdateEmailConfirmationOptions,
) {
  let user = await userRepo.findByEmailUpdateToken(options.token)
  user = await userRepo.update({
    id: user.id,
    email: user.emailUpdateValue,
    username: user.emailUpdateValue,
    emailUpdateToken: null,
    emailUpdateValue: null,
  })
  await search.index(USER_SEARCH_INDEX).updateDocuments([
    {
      id: user.id,
      username: user.username,
      email: user.email,
      fullName: user.fullName,
      isEmailConfirmed: user.isEmailConfirmed,
      createTime: user.createTime,
      updateTime: user.updateTime,
      picture: user.picture,
    },
  ])
  return mapEntity(user)
}

export async function updatePassword(
  id: string,
  options: UserUpdatePasswordOptions,
): Promise<UserDTO> {
  let user = await userRepo.findById(id)
  if (verifyPassword(options.currentPassword, user.passwordHash)) {
    user = await userRepo.update({
      id: user.id,
      passwordHash: hashPassword(options.newPassword),
    })
    return mapEntity(user)
  } else {
    throw newPasswordValidationFailedError()
  }
}

export async function updatePicture(
  id: string,
  path: string,
  contentType: string,
): Promise<UserDTO> {
  const picture = await fs.readFile(path, { encoding: 'base64' })
  const { id: userId } = await userRepo.findById(id)
  const user = await userRepo.update({
    id: userId,
    picture: `data:${contentType};base64,${picture}`,
  })
  await search.index(USER_SEARCH_INDEX).updateDocuments([
    {
      id: user.id,
      username: user.username,
      email: user.email,
      fullName: user.fullName,
      isEmailConfirmed: user.isEmailConfirmed,
      createTime: user.createTime,
      updateTime: user.updateTime,
      picture: user.picture,
    },
  ])
  return mapEntity(user)
}

export async function deletePicture(id: string): Promise<UserDTO> {
  let user = await userRepo.findById(id)
  user = await userRepo.update({ id: user.id, picture: null })
  await search.index(USER_SEARCH_INDEX).updateDocuments([
    {
      id: user.id,
      username: user.username,
      email: user.email,
      fullName: user.fullName,
      isEmailConfirmed: user.isEmailConfirmed,
      createTime: user.createTime,
      updateTime: user.updateTime,
      picture: user.picture,
    },
  ])
  return mapEntity(user)
}

export async function deleteUser(id: string, options: UserDeleteOptions) {
  const user = await userRepo.findById(id)
  if (verifyPassword(options.password, user.passwordHash)) {
    await userRepo.delete(user.id)
    await search.index(USER_SEARCH_INDEX).deleteDocuments([user.id])
  } else {
    throw newInvalidPasswordError()
  }
}

export async function suspendUser(id: string, options: UserSuspendOptions) {
  const user = await userRepo.findById(id)
  if (
    user.isAdmin &&
    !(await userRepo.enoughActiveAdmins()) &&
    options.suspend
  ) {
    throw newCannotSuspendSoleAdminError()
  }
  if (user) {
    await userRepo.suspend(user.id, options.suspend)
    await search.index(USER_SEARCH_INDEX).updateDocuments([
      {
        id: user.id,
        username: user.username,
        email: user.email,
        fullName: user.fullName,
        isEmailConfirmed: user.isEmailConfirmed,
        createTime: user.createTime,
        updateTime: user.updateTime,
        picture: user.picture,
      },
    ])
  } else {
    throw newUserNotFoundError()
  }
}

export async function makeAdminUser(id: string, options: UserMakeAdminOptions) {
  const user = await userRepo.findById(id)
  if (
    user.isAdmin &&
    !(await userRepo.enoughActiveAdmins()) &&
    !options.makeAdmin
  ) {
    throw newCannotDemoteSoleAdminError()
  }
  if (user) {
    await userRepo.makeAdmin(user.id, options.makeAdmin)
    await search.index(USER_SEARCH_INDEX).updateDocuments([
      {
        id: user.id,
        username: user.username,
        email: user.email,
        fullName: user.fullName,
        isEmailConfirmed: user.isEmailConfirmed,
        createTime: user.createTime,
        updateTime: user.updateTime,
        picture: user.picture,
      },
    ])
  } else {
    throw newError({ code: ErrorCode.UserNotFound })
  }
}

export function mapEntity(entity: User): UserDTO {
  const user: UserDTO = {
    id: entity.id,
    email: entity.email,
    username: entity.username,
    fullName: entity.fullName,
    pendingEmail: entity.emailUpdateValue,
  }
  if (entity.picture) {
    user.picture = {
      extension: base64ToExtension(entity.picture),
    }
  }
  return user
}

export function adminMapEntity(entity: User): UserAdminDTO {
  const user: UserAdminDTO = {
    id: entity.id,
    email: entity.email,
    username: entity.username,
    fullName: entity.fullName,
    createTime: entity.createTime,
    updateTime: entity.updateTime,
    isActive: entity.isActive,
    isAdmin: entity.isAdmin,
    isEmailConfirmed: entity.isEmailConfirmed,
  }
  if (entity.picture) {
    user.picture = {
      extension: base64ToExtension(entity.picture),
    }
  }
  return user
}
