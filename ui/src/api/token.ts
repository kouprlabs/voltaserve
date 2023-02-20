import { idpFetch } from './fetch'

export type TokenGrantType = 'password' | 'refresh_token'

export type Token = {
  access_token: string
  expires_in: number
  token_type: string
  refresh_token: string
}

export type TokenExchangeOptions = {
  grant_type: TokenGrantType
  username: string
  password: string
  refresh_token?: string
  locale?: string
}

export default class TokenAPI {
  static async exchange(options: TokenExchangeOptions): Promise<Token> {
    const formBody = []
    formBody.push(`grant_type=${options.grant_type}`)
    formBody.push(`username=${encodeURIComponent(options.username)}`)
    formBody.push(`password=${encodeURIComponent(options.password)}`)
    if (options.refresh_token) {
      formBody.push(`username=${encodeURIComponent(options.refresh_token)}`)
    }
    if (options.locale) {
      formBody.push(`&locale=${encodeURIComponent(options.locale)}`)
    }
    return idpFetch(
      `/token`,
      {
        method: 'POST',
        body: formBody.join('&'),
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
      },
      false
    ).then((result) => result.json())
  }
}
