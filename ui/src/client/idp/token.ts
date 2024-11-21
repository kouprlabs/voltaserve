// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { baseFetcher } from '@/client/fetcher'
import { getConfig } from '@/config/config'

export type GrantType = 'password' | 'refresh_token'

export type Token = {
  access_token: string
  expires_in: number
  token_type: string
  refresh_token: string
  is_admin: boolean
}

export type ExchangeOptions = {
  grant_type: GrantType
  username?: string
  password?: string
  refresh_token?: string
  locale?: string
}

export default class TokenAPI {
  static async exchange(options: ExchangeOptions): Promise<Token> {
    const formBody = []
    formBody.push(`grant_type=${options.grant_type}`)
    if (options.grant_type === 'password') {
      if (options.username && options.password) {
        formBody.push(`username=${encodeURIComponent(options.username)}`)
        formBody.push(`password=${encodeURIComponent(options.password)}`)
      } else {
        throw new Error('Username or password missing!')
      }
    }
    if (options.grant_type === 'refresh_token') {
      if (options.refresh_token) {
        formBody.push(
          `refresh_token=${encodeURIComponent(options.refresh_token)}`,
        )
      } else {
        throw new Error('Refresh token missing!')
      }
    }
    if (options.locale) {
      formBody.push(`&locale=${encodeURIComponent(options.locale)}`)
    }
    const response = await baseFetcher(
      `${getConfig().idpURL}/token`,
      {
        method: 'POST',
        body: formBody.join('&'),
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
      },
      false,
    )
    if (response) {
      return (await response.json()) as Token
    } else {
      throw new Error('Invalid token response')
    }
  }
}
