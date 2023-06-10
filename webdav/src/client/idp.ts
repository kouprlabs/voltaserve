import { IDP_URL } from '@/config'

export type IdPErrorResponse = {
  code: string
  status: number
  message: string
  userMessage: string
  moreInfo: string
}

export class IdPError extends Error {
  constructor(readonly error: IdPErrorResponse) {
    super(JSON.stringify(error, null, 2))
  }
}

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
    const response = await fetch(`${IDP_URL}/v1/token`, {
      method: 'POST',
      body: formBody.join('&'),
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded',
      },
    })
    return this.jsonResponseOrThrow(response)
  }

  private async jsonResponseOrThrow<T>(response: Response): Promise<T> {
    if (response.headers.get('content-type')?.includes('application/json')) {
      const json = await response.json()
      if (response.status > 299) {
        throw new IdPError(json)
      }
      return json
    } else {
      if (response.status > 299) {
        throw new Error(response.statusText)
      }
    }
  }
}
