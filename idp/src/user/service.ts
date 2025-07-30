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
  newPasswordValidationFailedError,
  newPictureNotFoundError,
  newUnsupportedOperationError,
  newUsernameUnavailableError,
  newUserNotFoundError,
} from '@/error/creators.ts'
import { ErrorCode, newError } from '@/error/core.ts'
import { newHyphenlessUuid } from '@/infra/id.ts'
import { sendTemplateMail } from '@/infra/mail.ts'
import { hashPassword, verifyPassword } from '@/infra/password.ts'
import {
  client as meilisearch,
  USER_SEARCH_INDEX,
} from '@/infra/meilisearch.ts'
import { isLocalStrategy, User } from '@/user/model.ts'
import userRepo from '@/user/repo.ts'
import { Buffer } from 'node:buffer'
import { call as callWebhook, UserWebhookEventType } from './webhook.ts'
import { logger } from '@/infra/logger.ts'

export type UserDTO = {
  id: string
  fullName: string
  picture?: PictureDTO
  email: string
  pendingEmail?: string
  username: string
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

export type UserUpdatePictureRawOptions = {
  picture?: string
}

export type UserUpdatePasswordOptions = {
  currentPassword: string
  newPassword: string
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

export function find(user: User): UserDTO {
  return mapEntity(user)
}

export async function findAsAdmin(id: string): Promise<UserAdminDTO> {
  return mapAdminEntity(await userRepo.findById(id))
}

export async function getPicture(id: string): Promise<UserPictureResponse> {
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

export async function getCount(): Promise<number> {
  return await userRepo.getCount()
}

export async function list({
  query,
  size,
  page,
}: UserListOptions): Promise<UserAdminList> {
  if (query) {
    const hits = await meilisearch
      .index(USER_SEARCH_INDEX)
      .search(query, { page: page, hitsPerPage: size })
      .then((value) => {
        return {
          data: value.hits,
          totalElements: value.totalHits,
        }
      })
    const users = await userRepo.findMany(hits.data.map((value) => value.id))
    return {
      data: users.map(mapAdminEntity),
      totalElements: hits.totalElements,
      totalPages: Math.floor((hits.totalElements + size - 1) / size),
      size: size,
      page: page,
    }
  } else {
    const users = await userRepo.list(page, size)
    return {
      data: users.map(mapAdminEntity),
      totalElements: await userRepo.getCount(),
      totalPages: Math.floor(((await userRepo.getCount()) + size - 1) / size),
      size: size,
      page: page,
    }
  }
}

export async function updateFullName(
  user: User,
  options: UserUpdateFullNameOptions,
): Promise<UserDTO> {
  user = await userRepo.update({ id: user.id, fullName: options.fullName })
  await meilisearch.index(USER_SEARCH_INDEX).updateDocuments([
    {
      id: user.id,
      username: user.username,
      email: user.email,
      fullName: user.fullName,
      isEmailConfirmed: user.isEmailConfirmed,
      createTime: user.createTime,
      updateTime: user.updateTime,
    },
  ])
  return mapEntity(user)
}

export async function updateEmailRequest(
  user: User,
  options: UserUpdateEmailRequestOptions,
): Promise<UserDTO> {
  if (!isLocalStrategy(user)) {
    throw newUnsupportedOperationError()
  }
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
        id: user.id,
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
  if (!isLocalStrategy(user)) {
    throw newUnsupportedOperationError()
  }
  user = await userRepo.update({
    id: user.id,
    email: user.emailUpdateValue,
    username: user.emailUpdateValue,
    emailUpdateToken: null,
    emailUpdateValue: null,
  })
  await meilisearch.index(USER_SEARCH_INDEX).updateDocuments([
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
  user: User,
  options: UserUpdatePasswordOptions,
): Promise<UserDTO> {
  if (!isLocalStrategy(user)) {
    throw newUnsupportedOperationError()
  }
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
  { id: userId }: User,
  path: string,
  contentType: string,
): Promise<UserDTO> {
  const picture = await fs.readFile(path, { encoding: 'base64' })
  const user = await userRepo.update({
    id: userId,
    picture: `data:${contentType};base64,${picture}`,
  })
  return mapEntity(user)
}

export async function updatePictureRaw(
  { id: userId }: User,
  picture?: string,
): Promise<UserDTO> {
  const user = await userRepo.update({
    id: userId,
    picture,
  })
  return mapEntity(user)
}

export async function deletePicture({ id }: User): Promise<UserDTO> {
  const user = await userRepo.update({ id, picture: null })
  return mapEntity(user)
}

export async function deleteUser(user: User) {
  if (getConfig().userWebhooks.length > 0) {
    const dto = mapEntity(user)
    if (getConfig().userWebhooks.length > 0) {
      for (const url of getConfig().userWebhooks) {
        try {
          await callWebhook(url, {
            eventType: UserWebhookEventType.Delete,
            user: dto,
          })
        } catch (error) {
          logger.error(error)
        }
      }
    }
  }
  await userRepo.delete(user.id)
  await meilisearch.index(USER_SEARCH_INDEX).deleteDocuments([user.id])
}

export async function suspend(id: string, options: UserSuspendOptions) {
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
    await meilisearch.index(USER_SEARCH_INDEX).updateDocuments([
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

export async function makeAdmin(id: string, options: UserMakeAdminOptions) {
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
    await meilisearch.index(USER_SEARCH_INDEX).updateDocuments([
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
  }
  if (entity.emailUpdateValue) {
    user.pendingEmail = entity.emailUpdateValue
  }
  if (entity.picture) {
    user.picture = {
      extension: base64ToExtension(entity.picture),
    }
  }
  return user
}

export function mapAdminEntity(entity: User): UserAdminDTO {
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
