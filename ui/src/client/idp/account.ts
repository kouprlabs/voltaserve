// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import useSWR, { SWRConfiguration } from 'swr'
import { idpFetcher } from '@/client/fetcher'
import { AuthUser } from './user'

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

export type AccountPasswordRequirements = {
  minLength: number
  minLowercase: number
  minUppercase: number
  minNumbers: number
  minSymbols: number
}

export class AccountAPI {
  static create(options: AccountCreateOptions) {
    return idpFetcher({
      url: `/accounts`,
      method: 'POST',
      body: JSON.stringify(options),
      redirect: false,
      authenticate: false,
    }) as Promise<AuthUser>
  }

  static async sendResetPasswordEmail(
    options: AccountSendResetPasswordEmailOptions,
  ) {
    return idpFetcher({
      url: `/accounts/send_reset_password_email`,
      method: 'POST',
      body: JSON.stringify(options),
    })
  }

  static async resetPassword(options: AccountResetPasswordOptions) {
    return idpFetcher({
      url: `/accounts/reset_password`,
      method: 'POST',
      body: JSON.stringify(options),
    })
  }

  static async confirmEmail(options: AccountConfirmEmailOptions) {
    return idpFetcher({
      url: `/accounts/confirm_email`,
      method: 'POST',
      body: JSON.stringify(options),
    })
  }

  static async getPasswordRequirements() {
    return idpFetcher({
      url: `/accounts/password_requirements`,
      method: 'GET',
    })
  }

  static useGetPasswordRequirements(swrOptions?: SWRConfiguration) {
    const url = `/accounts/password_requirements`
    return useSWR<AccountPasswordRequirements | undefined>(
      url,
      () =>
        idpFetcher({
          url,
          method: 'GET',
        }) as Promise<AccountPasswordRequirements>,
      swrOptions,
    )
  }
}
