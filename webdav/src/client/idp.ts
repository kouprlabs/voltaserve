import { IDP_URL } from '@/config'
import { ClientError } from './error'

export type Token = {
  access_token: string
  expires_in: number
  token_type: string
  refresh_token: string
}

export type TokenGrantType = 'password' | 'refresh_token'

export type TokenExchangeOptions = {
  grant_type: TokenGrantType
  username?: string
  password?: string
  refresh_token?: string
  locale?: string
}

export class TokenAPI {
  async exchange(options: TokenExchangeOptions): Promise<Token> {
    const formBody = []
    formBody.push(`grant_type=${options.grant_type}`)
    formBody.push(`username=${encodeURIComponent(options.username)}`)
    formBody.push(`password=${encodeURIComponent(options.password)}`)
    if (options.refresh_token) {
      formBody.push(`refresh_token=${options.refresh_token}`)
    }
    const result = await fetch(`${IDP_URL}/v1/token`, {
      method: 'POST',
      body: formBody.join('&'),
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded',
      },
    })
    const json = await result.json()
    if (result.status > 299) {
      throw new ClientError(json)
    }
    return json
  }
}
