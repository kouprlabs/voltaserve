import { Config } from './types'

let config: Config

export function getConfig(): Config {
  if (!config) {
    config = new Config()
    config.port = parseInt(process.env.PORT)
    readURLs(config)
    readToken(config)
    readCORS(config)
    readSearch(config)
    readSMTP(config)
  }
  return config
}

export function readURLs(config: Config) {
  config.publicUIURL = process.env.PUBLIC_UI_URL
  config.databaseURL = process.env.POSTGRES_URL
}

export function readToken(config: Config) {
  config.token.jwtSigningKey = process.env.TOKEN_JWT_SIGNING_KEY
  config.token.audience = process.env.TOKEN_AUDIENCE
  config.token.issuer = process.env.TOKEN_ISSUER
  if (process.env.TOKEN_ACCESS_TOKEN_LIFETIME) {
    config.token.accessTokenLifetime = parseInt(
      process.env.TOKEN_ACCESS_TOKEN_LIFETIME,
    )
  }
  if (process.env.TOKEN_REFRESH_TOKEN_LIFETIME) {
    config.token.refreshTokenLifetime = parseInt(
      process.env.TOKEN_REFRESH_TOKEN_LIFETIME,
    )
  }
}

export function readCORS(config: Config) {
  if (process.env.CORS_ORIGINS) {
    config.corsOrigins = process.env.CORS_ORIGINS.split(',')
    config.corsOrigins.forEach((e) => e.trim())
  }
}

export function readSearch(config: Config) {
  config.search.url = process.env.SEARCH_URL
}

export function readSMTP(config: Config) {
  config.smtp.host = process.env.SMTP_HOST
  if (process.env.SMTP_PORT) {
    config.smtp.port = parseInt(process.env.SMTP_PORT)
  }
  if (process.env.SMTP_SECURE) {
    config.smtp.secure = process.env.SMTP_SECURE === 'true'
  }
  config.smtp.username = process.env.SMTP_USERNAME
  config.smtp.password = process.env.SMTP_PASSWORD
  config.smtp.senderAddress = process.env.SMTP_SENDER_ADDRESS
  config.smtp.senderName = process.env.SMTP_SENDER_NAME
}
