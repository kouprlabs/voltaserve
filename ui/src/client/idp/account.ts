import { idpFetcher } from '@/client/fetcher'
import { User } from './user'

export type CreateOptions = {
  email: string
  password: string
  fullName: string
  picture?: string
}

export type SendResetPasswordEmailOptions = {
  email: string
}

export type ResetPasswordOptions = {
  token: string
  newPassword: string
}

export type ConfirmEmailOptions = {
  token: string
}

export default class AccountAPI {
  static async create(options: CreateOptions) {
    return idpFetcher({
      url: `/accounts`,
      method: 'POST',
      body: JSON.stringify(options),
      redirect: false,
      authenticate: false,
    }) as Promise<User>
  }

  static async sendResetPasswordEmail(options: SendResetPasswordEmailOptions) {
    return idpFetcher({
      url: `/accounts/send_reset_password_email`,
      method: 'POST',
      body: JSON.stringify(options),
    })
  }

  static async resetPassword(options: ResetPasswordOptions) {
    return idpFetcher({
      url: `/accounts/reset_password`,
      method: 'POST',
      body: JSON.stringify(options),
    })
  }

  static async confirmEmail(options: ConfirmEmailOptions) {
    return idpFetcher({
      url: `/accounts/confirm_email`,
      method: 'POST',
      body: JSON.stringify(options),
    })
  }
}
