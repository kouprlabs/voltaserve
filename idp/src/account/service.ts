import { getConfig } from '@/config/config'
import { newDateTime } from '@/infra/date-time'
import { UserRepo } from '@/infra/db'
import { ErrorCode, newError } from '@/infra/error'
import { newHashId, newHyphenlessUuid } from '@/infra/id'
import { sendTemplateMail } from '@/infra/mail'
import { hashPassword } from '@/infra/password'
import search, { USER_SEARCH_INDEX } from '@/infra/search'
import { mapEntity, User } from '@/user/service'

export type AccountCreateOptions = {
  email: string
  password: string
  fullName: string
  picture?: string
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

export async function createUser(options: AccountCreateOptions): Promise<User> {
  const id = newHashId()
  const existingUser = await UserRepo.find('username', options.email)
  if (existingUser) {
    throw newError({ code: ErrorCode.UsernameUnavailable })
  }
  try {
    const emailConfirmationToken = newHyphenlessUuid()
    const user = await UserRepo.insert({
      id,
      username: options.email,
      email: options.email,
      fullName: options.fullName,
      picture: options.picture,
      passwordHash: hashPassword(options.password),
      emailConfirmationToken,
      createTime: newDateTime(),
    })
    await search.index(USER_SEARCH_INDEX).addDocuments([
      {
        id: user.id,
        username: user.username,
        email: user.email,
        fullName: user.fullName,
        isEmailConfirmed: user.isEmailConfirmed,
        createTime: user.createTime,
      },
    ])
    await sendTemplateMail('email-confirmation', options.email, {
      'UI_URL': getConfig().uiURL,
      'TOKEN': emailConfirmationToken,
    })
    return mapEntity(user)
  } catch (error) {
    await UserRepo.delete(id)
    await search.index(USER_SEARCH_INDEX).deleteDocuments([id])
    throw newError({ code: ErrorCode.InternalServerError, error })
  }
}

export async function resetPassword(options: AccountResetPasswordOptions) {
  const user = await UserRepo.find('reset_password_token', options.token, true)
  await UserRepo.update({
    id: user.id,
    passwordHash: hashPassword(options.newPassword),
  })
}

export async function confirmEmail(options: AccountConfirmEmailOptions) {
  let user = await UserRepo.find(
    'email_confirmation_token',
    options.token,
    true
  )
  user = await UserRepo.update({
    id: user.id,
    isEmailConfirmed: true,
    emailConfirmationToken: null,
  })
  await search.index(USER_SEARCH_INDEX).updateDocuments([
    {
      ...user,
      isEmailConfirmed: user.isEmailConfirmed,
    },
  ])
}

export async function sendResetPasswordEmail(
  options: AccountSendResetPasswordEmailOptions
) {
  let user = await UserRepo.find('email', options.email)
  if (user) {
    user = await UserRepo.update({
      id: user.id,
      resetPasswordToken: newHyphenlessUuid(),
    })
  } else {
    return
  }
  try {
    await sendTemplateMail('reset-password', user.email, {
      'UI_URL': getConfig().uiURL,
      'TOKEN': user.resetPasswordToken,
    })
  } catch (error) {
    const { id } = await UserRepo.find('email', options.email, true)
    await UserRepo.update({ id, resetPasswordToken: null })
    throw newError({ code: ErrorCode.InternalServerError, error })
  }
}
