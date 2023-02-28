export class Config {
  idpURL: string
  uiURL: string
  databaseURL: string
  token: TokenConfig
  corsOrigins: string[]
  search: SearchConfig
  smtp: SMTPConfig

  constructor() {
    this.token = new TokenConfig()
    this.search = new SearchConfig()
    this.smtp = new SMTPConfig()
  }
}

export class DatabaseConfig {
  url: string
}

export class TokenConfig {
  jwtSigningKey: string
  audience: string
  issuer: string
  accessTokenLifetime: number
  refreshTokenLifetime: number
}

export class SearchConfig {
  url: string
}

export class SMTPConfig {
  host: string
  port: number
  secure: boolean
  username?: string
  password?: string
  senderAddress: string
  senderName: string
}
