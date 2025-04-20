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
  newUnsupportedOperationError,
  newUsernameUnavailableError,
} from '@/error/creators.ts'
import { newHashId, newHyphenlessUuid } from '@/infra/id.ts'
import { sendTemplateMail } from '@/infra/mail.ts'
import { hashPassword } from '@/infra/password.ts'
import {
  client as meilisearch,
  USER_SEARCH_INDEX,
} from '@/infra/meilisearch.ts'
import { logger } from '@/infra/logger.ts'
import { isLocalStrategy, User } from '@/user/model.ts'
import userRepo from '@/user/repo.ts'
import { getCount, mapEntity, UserDTO } from '@/user/service.ts'
import { call as callWebhook, UserWebhookEventType } from '@/user/webhook.ts'

export type AccountCreateOptions = {
  username: string
  email: string
  password?: string
  fullName: string
  picture?: string
  emailConfirmationToken?: string
  isAdmin?: boolean
  isEmailConfirmed?: boolean
}

export type AccountSignUpWithLocalOptions = {
  email: string
  password: string
  fullName: string
  picture?: string
  isAdmin?: boolean
}

export type AccountSignUpWithAppleOptions = {
  payload: any
  appleFullName?: string
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

export async function createUser(options: AccountCreateOptions): Promise<User> {
  const id = newHashId()
  if (!(await userRepo.isUsernameAvailable(options.username))) {
    throw newUsernameUnavailableError()
  }
  try {
    const user = await userRepo.insert({
      id,
      username: options.username,
      email: options.email,
      fullName: options.fullName,
      picture: options.picture,
      passwordHash: options.password
        ? hashPassword(options.password)
        : undefined,
      emailConfirmationToken: options.emailConfirmationToken,
      createTime: newDateTime(),
      isAdmin: options.isAdmin,
      isEmailConfirmed: options.isEmailConfirmed,
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
    if (options.emailConfirmationToken) {
      await sendTemplateMail('email-confirmation', options.email, {
        'UI_URL': getConfig().publicUIURL,
        'TOKEN': options.emailConfirmationToken,
      })
    }
    if (getConfig().userWebhooks.length > 0) {
      for (const url of getConfig().userWebhooks) {
        try {
          await callWebhook(url, {
            eventType: UserWebhookEventType.Create,
            user: mapEntity(user),
          })
        } catch (error) {
          logger.error(error)
        }
      }
    }
    return user
  } catch (error) {
    await userRepo.delete(id)
    await meilisearch.index(USER_SEARCH_INDEX).deleteDocuments([id])
    throw newInternalServerError(error)
  }
}

export async function signUpWithLocal(
  options: AccountSignUpWithLocalOptions,
): Promise<UserDTO> {
  return mapEntity(
    await createUser({
      username: options.email,
      email: options.email,
      password: options.password,
      fullName: options.fullName,
      picture: options.picture,
      emailConfirmationToken: newHyphenlessUuid(),
      isAdmin: (await getCount()) === 0,
    }),
  )
}

export async function signUpWithApple(
  options: AccountSignUpWithAppleOptions,
): Promise<User> {
  return await createUser({
    username: options.payload.sub,
    email: options.payload.email.toLocaleLowerCase(),
    fullName: options.appleFullName ?? options.payload.email,
    isEmailConfirmed: true,
  })
}

export async function resetPassword(options: AccountResetPasswordOptions) {
  const user = await userRepo.findByResetPasswordToken(options.token)
  if (!isLocalStrategy(user)) {
    throw newUnsupportedOperationError()
  }
  await userRepo.update({
    id: user.id,
    passwordHash: hashPassword(options.newPassword),
  })
}

export async function confirmEmail(options: AccountConfirmEmailOptions) {
  let user = await userRepo.findByEmailConfirmationToken(options.token)
  if (!isLocalStrategy(user)) {
    throw newUnsupportedOperationError()
  }
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
  let user = await userRepo.findByEmail(options.email)
  if (!isLocalStrategy(user)) {
    throw newUnsupportedOperationError()
  }
  try {
    user = await userRepo.update({
      id: user.id,
      resetPasswordToken: newHyphenlessUuid(),
    })
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
