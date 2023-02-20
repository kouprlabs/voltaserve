import fs from 'fs'
import yaml from 'js-yaml'

export type DatabaseConfig = {
  url: string
}

export type TokenConfig = {
  jwtSigningKey: string
  audience: string
  issuer: string
  accessTokenLifetime: number
  refreshTokenLifetime: number
}

export type SearchConfig = {
  url: string
}

export type SmtpConfig = {
  host: string
  port: number
  secure: boolean
  username?: string
  password?: string
  senderAddress: string
  senderName: string
}

export type Config = {
  url: string
  databaseUrl: string
  webUrl: string
  token: TokenConfig
  corsOrigins: string[]
  search: SearchConfig
  smtp: SmtpConfig
}

let config: Config

export function getConfig(): Config {
  if (!config) {
    const filename = fs.existsSync('./config.local.yml')
      ? './config.local.yml'
      : './config.yml'
    config = yaml.load(fs.readFileSync(filename, 'utf8')) as Config
    if (process.env.URL) {
      config.url = process.env.URL
    }
    if (process.env.DATABASE_URL) {
      config.databaseUrl = process.env.DATABASE_URL
    }
    if (process.env.WEB_URL) {
      config.webUrl = process.env.WEB_URL
    }
    overrideToken(config)
    if (process.env.CORS_ORIGINS) {
      config.corsOrigins = process.env.CORS_ORIGINS.split(',')
      config.corsOrigins.forEach((e) => e.trim())
    }
    overrideSearch(config)
    overrideSmtp(config)
  }
  return config
}

function overrideToken(config: Config) {
  if (process.env.TOKEN_JWT_SIGNING_KEY) {
    config.token.jwtSigningKey = process.env.TOKEN_JWT_SIGNING_KEY
  }
  if (process.env.TOKEN_AUDIENCE) {
    config.token.audience = process.env.TOKEN_AUDIENCE
  }
  if (process.env.TOKEN_ISSUER) {
    config.token.audience = process.env.TOKEN_ISSUER
  }
  if (process.env.TOKEN_ACCESS_TOKEN_LIFETIME) {
    config.token.accessTokenLifetime = parseInt(
      process.env.TOKEN_ACCESS_TOKEN_LIFETIME
    )
  }
  if (process.env.TOKEN_REFRESH_TOKEN_LIFETIME) {
    config.token.refreshTokenLifetime = parseInt(
      process.env.TOKEN_REFRESH_TOKEN_LIFETIME
    )
  }
}

function overrideSearch(config: Config) {
  if (process.env.SEARCH_URL) {
    config.search.url = process.env.SEARCH_URL
  }
}

function overrideSmtp(config: Config) {
  if (process.env.SMTP_HOST) {
    config.smtp.host = process.env.SMTP_HOST
  }
  if (process.env.SMTP_PORT) {
    config.smtp.port = parseInt(process.env.SMTP_PORT)
  }
  if (process.env.SMTP_SECURE) {
    config.smtp.secure = Boolean(process.env.SMTP_SECURE)
  }
  if (process.env.SMTP_USERNAME) {
    config.smtp.username = process.env.SMTP_USERNAME
  }
  if (process.env.SMTP_PASSWORD) {
    config.smtp.password = process.env.SMTP_PASSWORD
  }
  if (process.env.SMTP_SENDER_ADDRESS) {
    config.smtp.senderAddress = process.env.SMTP_SENDER_ADDRESS
  }
  if (process.env.SMTP_SENDER_NAME) {
    config.smtp.senderName = process.env.SMTP_SENDER_NAME
  }
}
