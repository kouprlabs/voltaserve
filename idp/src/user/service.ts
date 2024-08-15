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
import {getConfig} from '@/config/config'
import {ErrorCode, newError} from '@/infra/error'
import {newHyphenlessUuid} from '@/infra/id'
import {sendTemplateMail} from '@/infra/mail'
import {hashPassword, verifyPassword} from '@/infra/password'
import search, {USER_SEARCH_INDEX} from '@/infra/search'
import {User} from '@/user/model'
import userRepo from '@/user/repo'
import {UserAdminRequest, UserSuspendRequest} from "@/infra/admin-requests";

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

export async function getUserListPaginated(
  page: number,
  size: number,
): Promise<User[]> {
  return await userRepo.listAllPaginated(page, size)
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
      ...user,
      fullName: user.fullName,
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
      ...user,
      email: user.email,
      username: user.email,
      emailUpdateToken: null,
      emailUpdateValue: null,
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
  return mapEntity(user)
}

export async function deletePicture(id: string): Promise<UserDTO> {
  let user = await userRepo.findByID(id)
  user = await userRepo.update({ id: user.id, picture: null })
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

export async function suspendUser(options: UserSuspendRequest) {
  const user = await userRepo.findByID(options.id)
  if (user.isAdmin && !await userRepo.enoughActiveAdmins() && options.suspend) {
    throw newError({ code: ErrorCode.OrphanError })
  }
  if (user) {
    await userRepo.suspend(user.id, options.suspend)
  } else {
    throw newError({ code: ErrorCode.UserNotFound })
  }
}

export async function makeAdminUser(options: UserAdminRequest) {
  const user = await userRepo.findByID(options.id)
  if (user.isAdmin && !await userRepo.enoughActiveAdmins() && options.makeAdmin) {
    throw newError({ code: ErrorCode.OrphanError })
  }
  if (user) {
    await userRepo.makeAdmin(user.id, options.makeAdmin)
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
