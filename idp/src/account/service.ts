// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { getConfig } from '@/config/config.ts'
import { newDateTime } from '@/infra/date-time.ts'
import {
  newInternalServerError,
  newUsernameUnavailableError,
} from '@/infra/error/creators.ts'
import { newHashId, newHyphenlessUuid } from '@/infra/id.ts'
import { sendTemplateMail } from '@/infra/mail.ts'
import { hashPassword } from '@/infra/password.ts'
import {
  client as meilisearch,
  USER_SEARCH_INDEX,
} from '../infra/meilisearch.ts'
import { User } from '@/user/model.ts'
import userRepo from '@/user/repo.ts'
import { getUserCount, mapEntity, UserDTO } from '@/user/service.ts'

export type AccountCreateOptions = {
  email: string
  password: string
  fullName: string
  picture?: string
  isAdmin?: boolean
}

export type AccountResetPasswordOptions = {
  token: string
  newPassword: string
}

export type AccountConfirmEmailOptions = {
  token: string
}

export type AccountSendResetPasswordEmailOptions = {
  email: string
}

export type PasswordRequirements = {
  minLength: number
  minLowercase: number
  minUppercase: number
  minNumbers: number
  minSymbols: number
}

export async function createUser(
  options: AccountCreateOptions,
): Promise<UserDTO> {
  const id = newHashId()
  if (!(await userRepo.isUsernameAvailable(options.email))) {
    throw newUsernameUnavailableError()
  }
  if ((await getUserCount()) === 0) {
    options.isAdmin = true
  }
  try {
    const emailConfirmationToken = newHyphenlessUuid()
    const user = await userRepo.insert({
      id,
      username: options.email.toLocaleLowerCase(),
      email: options.email.toLocaleLowerCase(),
      fullName: options.fullName,
      picture: options.picture,
      passwordHash: hashPassword(options.password),
      emailConfirmationToken,
      createTime: newDateTime(),
      isAdmin: options.isAdmin,
    })
    await meilisearch.index(USER_SEARCH_INDEX).addDocuments([
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
    await sendTemplateMail('email-confirmation', options.email, {
      'UI_URL': getConfig().publicUIURL,
      'TOKEN': emailConfirmationToken,
    })
    return mapEntity(user)
  } catch (error) {
    await userRepo.delete(id)
    await meilisearch.index(USER_SEARCH_INDEX).deleteDocuments([id])
    throw newInternalServerError(error)
  }
}

export async function resetPassword(options: AccountResetPasswordOptions) {
  const user = await userRepo.findByResetPasswordToken(options.token)
  await userRepo.update({
    id: user.id,
    passwordHash: hashPassword(options.newPassword),
  })
}

export async function confirmEmail(options: AccountConfirmEmailOptions) {
  let user = await userRepo.findByEmailConfirmationToken(options.token)
  user = await userRepo.update({
    id: user.id,
    isEmailConfirmed: true,
    emailConfirmationToken: null,
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
}

export async function sendResetPasswordEmail(
  options: AccountSendResetPasswordEmailOptions,
) {
  let user: User
  try {
    user = await userRepo.findByEmail(options.email)
    user = await userRepo.update({
      id: user.id,
      resetPasswordToken: newHyphenlessUuid(),
    })
  } catch {
    return
  }
  try {
    await sendTemplateMail('reset-password', user.email, {
      'UI_URL': getConfig().publicUIURL,
      'TOKEN': user.resetPasswordToken,
    })
  } catch (error) {
    const { id } = await userRepo.findByEmail(options.email)
    await userRepo.update({ id, resetPasswordToken: null })
    throw newInternalServerError(error)
  }
}

export function getPasswordRequirements(): PasswordRequirements {
  return {
    minLength: getConfig().password.minLength,
    minLowercase: getConfig().password.minLowercase,
    minUppercase: getConfig().password.minUppercase,
    minNumbers: getConfig().password.minNumbers,
    minSymbols: getConfig().password.minSymbols,
  } as PasswordRequirements
}
