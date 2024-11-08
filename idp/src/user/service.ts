// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import fs from 'fs/promises'
import { getConfig } from '@/config/config'
import { base64ToBuffer, base64ToExtension, base64ToMIME } from '@/infra/base64'
import { ErrorCode, newError } from '@/infra/error'
import { newHyphenlessUuid } from '@/infra/id'
import { sendTemplateMail } from '@/infra/mail'
import { hashPassword, verifyPassword } from '@/infra/password'
import search, { USER_SEARCH_INDEX } from '@/infra/search'
import { User } from '@/user/model'
import userRepo from '@/user/repo'

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

export type SearchRequest = {
  page: string
  size: string
  query: string
}

export interface UserSuspendOptions {
  suspend: boolean
}

export interface UserMakeAdminOptions {
  makeAdmin: boolean
}

export async function getUser(id: string): Promise<UserDTO> {
  return mapEntity(await userRepo.findByID(id))
}

export async function getUserByAdmin(id: string): Promise<UserAdminDTO> {
  return adminMapEntity(await userRepo.findByID(id))
}

export async function getUserPicture(id: string): Promise<UserPictureResponse> {
  const user = await userRepo.findByID(id)
  if (!user.picture) {
    throw newError({
      code: ErrorCode.ResourceNotFound,
      message: 'Picture not found',
      userMessage: 'Picture not found',
    })
  }
  return {
    buffer: base64ToBuffer(user.picture),
    extension: base64ToExtension(user.picture),
    mime: base64ToMIME(user.picture),
  }
}

export async function searchUserListPaginated(
  query: string,
  size: number,
  page: number,
): Promise<UserAdminList> {
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
        await userRepo.listAllByIds(
          users.data.map((value) => {
            return value.id
          }),
        )
      ).map((value) => adminMapEntity(value)),
      totalElements: users.totalElements,
      size: size,
      page: page,
    }
  } else {
    return {
      data: (await userRepo.listAllPaginated(page, size)).map((value) =>
        adminMapEntity(value),
      ),
      totalElements: await userRepo.getUserCount(),
      size: size,
      page: page,
    }
  }
}

export async function getUserCount(): Promise<number> {
  return await userRepo.getUserCount()
}

export async function updateFullName(
  id: string,
  options: UserUpdateFullNameOptions,
): Promise<UserDTO> {
  let user = await userRepo.findByID(id)
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
  let user = await userRepo.findByID(id)
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
      throw newError({ code: ErrorCode.UsernameUnavailable })
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
      throw newError({ code: ErrorCode.InternalServerError, error })
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
  let user = await userRepo.findByID(id)
  if (verifyPassword(options.currentPassword, user.passwordHash)) {
    user = await userRepo.update({
      id: user.id,
      passwordHash: hashPassword(options.newPassword),
    })
    return mapEntity(user)
  } else {
    throw newError({ code: ErrorCode.PasswordValidationFailed })
  }
}

export async function updatePicture(
  id: string,
  path: string,
  contentType: string,
): Promise<UserDTO> {
  const picture = await fs.readFile(path, { encoding: 'base64' })
  const { id: userId } = await userRepo.findByID(id)
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
  let user = await userRepo.findByID(id)
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
  const user = await userRepo.findByID(id)
  if (verifyPassword(options.password, user.passwordHash)) {
    await userRepo.delete(user.id)
    await search.index(USER_SEARCH_INDEX).deleteDocuments([user.id])
  } else {
    throw newError({ code: ErrorCode.InvalidPassword })
  }
}

export async function suspendUser(id: string, options: UserSuspendOptions) {
  const user = await userRepo.findByID(id)
  if (
    user.isAdmin &&
    !(await userRepo.enoughActiveAdmins()) &&
    options.suspend
  ) {
    throw newError({ code: ErrorCode.OrphanError })
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
    throw newError({ code: ErrorCode.UserNotFound })
  }
}

export async function makeAdminUser(id: string, options: UserMakeAdminOptions) {
  const user = await userRepo.findByID(id)
  if (
    user.isAdmin &&
    !(await userRepo.enoughActiveAdmins()) &&
    options.makeAdmin
  ) {
    throw newError({ code: ErrorCode.OrphanError })
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
  Object.keys(user).forEach(
    (index) => !user[index] && user[index] !== undefined && delete user[index],
  )
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
  Object.keys(user).forEach(
    (index) => !user[index] && user[index] !== undefined && delete user[index],
  )
  return user
}
