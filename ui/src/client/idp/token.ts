import { idpFetch } from '@/client/fetch'

export type GrantType = 'password' | 'refresh_token'

export type Token = {
  access_token: string
  expires_in: number
  token_type: string
  refresh_token: string
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
    return idpFetch(
      '/token',
      {
        method: 'POST',
        body: formBody.join('&'),
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
      },
      false,
    ).then((result) => result.json())
  }
}
