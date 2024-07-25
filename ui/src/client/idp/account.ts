// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import useSWR, { SWRConfiguration } from 'swr'
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

export type PasswordRequirements = {
  minLength: number
  minLowercase: number
  minUppercase: number
  minNumbers: number
  minSymbols: number
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

  static async getPasswordRequirements() {
    return idpFetcher({
      url: `/accounts/password_requirements`,
      method: 'GET',
    })
  }

  static useGetPasswordRequirements(swrOptions?: SWRConfiguration) {
    const url = `/accounts/password_requirements`
    return useSWR<PasswordRequirements | undefined>(
      url,
      () => idpFetcher({ url, method: 'GET' }) as Promise<PasswordRequirements>,
      swrOptions,
    )
  }
}
