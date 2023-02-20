import { idpFetch } from './fetch'
import { User } from './user'

export type AccountCreateOptions = {
  email: string
  password: string
  fullName: string
  picture?: string
}

export type AccountSendResetPasswordEmailOptions = {
  email: string
}

export type AccountResetPasswordOptions = {
  token: string
  newPassword: string
}

export type AccountConfirmEmailOptions = {
  token: string
}

export default class AccountAPI {
  static async create(options: AccountCreateOptions): Promise<User> {
    return idpFetch(`/accounts`, {
      method: 'POST',
      body: JSON.stringify(options),
      headers: {
        'Content-Type': 'application/json',
      },
    }).then((result) => result.json())
  }

  static async sendResetPasswordEmail(
    options: AccountSendResetPasswordEmailOptions
  ) {
    return idpFetch(`/accounts/send_reset_password_email`, {
      method: 'POST',
      body: JSON.stringify(options),
      headers: {
        'Content-Type': 'application/json',
      },
    })
  }

  static async resetPassword(options: AccountResetPasswordOptions) {
    return idpFetch(`/accounts/reset_password`, {
      method: 'POST',
      body: JSON.stringify(options),
      headers: {
        'Content-Type': 'application/json',
      },
    })
  }

  static async confirmEmail(options: AccountConfirmEmailOptions) {
    return idpFetch(`/accounts/confirm_email`, {
      method: 'POST',
      body: JSON.stringify(options),
      headers: {
        'Content-Type': 'application/json',
      },
    })
  }
}
