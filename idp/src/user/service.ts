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
import {
  UserAdminPostRequest,
  UserIdPostRequest,
  UserSearchResponse,
  UserSuspendPostRequest,
  UserUpdateAdminRequest,
} from '@/infra/admin-requests'
import { ErrorCode, newError } from '@/infra/error'
import { newHyphenlessUuid } from '@/infra/id'
import { sendTemplateMail } from '@/infra/mail'
import { hashPassword, verifyPassword } from '@/infra/password'
import search, { USER_SEARCH_INDEX } from '@/infra/search'
import { UpdateOptions, User } from '@/user/model'
import userRepo from '@/user/repo'

export type UserDTO = {
  id: string
  fullName: string
  picture: string
  email: string
  username: string
  pendingEmail?: string
}

export type UserListDTO = {
  id: string
  fullName: string
  picture: string
  email: string
  username: string
  pendingEmail?: string
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

export async function getUser(id: string): Promise<UserDTO> {
  return mapEntity(await userRepo.findByID(id))
}

export async function getUserByAdmin(id: string): Promise<User> {
  return adminMapEntity(await userRepo.findByID(id))
}

export async function updateAdminUser(
  id: string,
  data: UserUpdateAdminRequest,
): Promise<UserDTO> {
  if (!data.isEmailConfirmed) {
    let usernameUnavailable = false
    try {
      await userRepo.findByUsername(data.email)
      usernameUnavailable = true
    } catch {
      // Ignored
    }
    if (usernameUnavailable) {
      throw newError({ code: ErrorCode.UsernameUnavailable })
    }
    const user = await userRepo.update({ id: id, ...(data as UpdateOptions) })
    try {
      await sendTemplateMail('email-confirmation', user.email, {
        'UI_URL': getConfig().publicUIURL,
        'TOKEN': user.emailConfirmationToken,
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
  } else {
    return mapEntity(
      await userRepo.update({ id: id, ...(data as UpdateOptions) }),
    )
  }
}

export async function searchUserListPaginated(
  query: string,
  size: number,
  page: number,
): Promise<UserSearchResponse> {
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
      data: await userRepo.listAllByIds(
        users.data.map((value) => {
          return value.id
        }),
      ),
      totalElements: users.totalElements,
      size: size,
      page: page,
    }
  } else {
    return {
      data: await userRepo.listAllPaginated(page, size),
      totalElements: await userRepo.getUserCount(),
      size: size,
      page: page,
    }
  }
}

export async function getUserCount(): Promise<number> {
  return await userRepo.getUserCount()
}

export async function getByPicture(picture: string): Promise<UserDTO> {
  return mapEntity(await userRepo.findByPicture(picture))
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

export const raiseSearchError = () => {
  throw newError({ code: ErrorCode.SearchError })
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

export async function suspendUser(options: UserSuspendPostRequest) {
  const user = await userRepo.findByID(options.id)
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

export async function makeAdminUser(options: UserAdminPostRequest) {
  const user = await userRepo.findByID(options.id)
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

export async function forceResetPassword(options: UserIdPostRequest) {
  const user = await userRepo.findByID(options.id)
  if (user) {
    const token = newHyphenlessUuid()
    await userRepo.forceResetPassword(user.id, token)
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
    picture: entity.picture,
    pendingEmail: entity.emailUpdateValue,
  }
  Object.keys(user).forEach(
    (index) => !user[index] && user[index] !== undefined && delete user[index],
  )
  return user
}

export function adminMapEntity(entity: User): User {
  const user: User = {
    id: entity.id,
    email: entity.email,
    username: entity.username,
    fullName: entity.fullName,
    picture: entity.picture,
    createTime: entity.createTime,
    updateTime: entity.updateTime,
    isActive: entity.isActive,
    isAdmin: entity.isAdmin,
    isEmailConfirmed: entity.isEmailConfirmed,
    forceChangePassword: entity.forceChangePassword,
  }
  Object.keys(user).forEach(
    (index) => !user[index] && user[index] !== undefined && delete user[index],
  )
  return user
}
