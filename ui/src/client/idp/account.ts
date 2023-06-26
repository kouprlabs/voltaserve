import { idpFetch } from '@/client/fetch'
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
  static async create(options: CreateOptions): Promise<User> {
    return idpFetch(`/accounts`, {
      method: 'POST',
      body: JSON.stringify(options),
      headers: {
        'Content-Type': 'application/json',
      },
    }).then((result) => result.json())
  }

  static async sendResetPasswordEmail(options: SendResetPasswordEmailOptions) {
    return idpFetch(`/accounts/send_reset_password_email`, {
      method: 'POST',
      body: JSON.stringify(options),
      headers: {
        'Content-Type': 'application/json',
      },
    })
  }

  static async resetPassword(options: ResetPasswordOptions) {
    return idpFetch(`/accounts/reset_password`, {
      method: 'POST',
      body: JSON.stringify(options),
      headers: {
        'Content-Type': 'application/json',
      },
    })
  }

  static async confirmEmail(options: ConfirmEmailOptions) {
    return idpFetch(`/accounts/confirm_email`, {
      method: 'POST',
      body: JSON.stringify(options),
      headers: {
        'Content-Type': 'application/json',
      },
    })
  }
}
