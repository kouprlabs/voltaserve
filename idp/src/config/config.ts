// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import * as process from 'node:process'
import { Config, Strategy } from '@/config/types.ts'

let config: Config

export function getConfig(): Config {
  if (!process.env.PORT) {
    throw newEnvironmentVariableNotSetError('PORT')
  }
  if (!config) {
    config = new Config()
    config.port = parseInt(process.env.PORT)
    if (process.env.STRATEGY && isValidStrategy(process.env.STRATEGY)) {
      config.strategy = process.env.STRATEGY as Strategy
    }
    readURLs(config)
    readToken(config)
    readPassword(config)
    readCORS(config)
    readSearch(config)
    readSMTP(config)
    readSecurity(config)
    readWebhooks(config)
  }
  return config
}

function isValidStrategy(value: string): boolean {
  return value === Strategy.Local || value === Strategy.Apple
}

export function readURLs(config: Config) {
  if (!process.env.PUBLIC_UI_URL) {
    throw newEnvironmentVariableNotSetError('PUBLIC_UI_URL')
  }
  if (!process.env.POSTGRES_URL) {
    throw newEnvironmentVariableNotSetError('POSTGRES_URL')
  }
  config.publicUIURL = process.env.PUBLIC_UI_URL
  config.databaseURL = process.env.POSTGRES_URL
}

export function readToken(config: Config) {
  if (!process.env.TOKEN_JWT_SIGNING_KEY) {
    throw newEnvironmentVariableNotSetError('TOKEN_JWT_SIGNING_KEY')
  }
  if (!process.env.TOKEN_AUDIENCE) {
    throw newEnvironmentVariableNotSetError('TOKEN_AUDIENCE')
  }
  if (!process.env.TOKEN_ISSUER) {
    throw newEnvironmentVariableNotSetError('TOKEN_ISSUER')
  }
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

export function readPassword(config: Config) {
  if (!process.env.PASSWORD_MIN_LENGTH) {
    throw newEnvironmentVariableNotSetError('PASSWORD_MIN_LENGTH')
  }
  if (!process.env.PASSWORD_MIN_LOWERCASE) {
    throw newEnvironmentVariableNotSetError('PASSWORD_MIN_LOWERCASE')
  }
  if (!process.env.PASSWORD_MIN_UPPERCASE) {
    throw newEnvironmentVariableNotSetError('PASSWORD_MIN_UPPERCASE')
  }
  if (!process.env.PASSWORD_MIN_NUMBERS) {
    throw newEnvironmentVariableNotSetError('PASSWORD_MIN_NUMBERS')
  }
  if (!process.env.PASSWORD_MIN_SYMBOLS) {
    throw newEnvironmentVariableNotSetError('PASSWORD_MIN_SYMBOLS')
  }
  config.password.minLength = parseInt(process.env.PASSWORD_MIN_LENGTH)
  config.password.minLowercase = parseInt(process.env.PASSWORD_MIN_LOWERCASE)
  config.password.minUppercase = parseInt(process.env.PASSWORD_MIN_UPPERCASE)
  config.password.minNumbers = parseInt(process.env.PASSWORD_MIN_NUMBERS)
  config.password.minSymbols = parseInt(process.env.PASSWORD_MIN_SYMBOLS)
}

export function readCORS(config: Config) {
  if (process.env.CORS_ORIGINS) {
    config.corsOrigins = process.env.CORS_ORIGINS.split(',')
    config.corsOrigins.forEach((e) => e.trim())
  }
}

export function readSearch(config: Config) {
  if (!process.env.SEARCH_URL) {
    throw newEnvironmentVariableNotSetError('SEARCH_URL')
  }
  config.search.url = process.env.SEARCH_URL
}

export function readSMTP(config: Config) {
  if (!process.env.SMTP_HOST) {
    throw newEnvironmentVariableNotSetError('SMTP_HOST')
  }
  if (!process.env.SMTP_SENDER_ADDRESS) {
    throw newEnvironmentVariableNotSetError('SMTP_SENDER_ADDRESS')
  }
  if (!process.env.SMTP_SENDER_NAME) {
    throw newEnvironmentVariableNotSetError('SMTP_SENDER_NAME')
  }
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

export function readSecurity(config: Config) {
  if (!process.env.SECURITY_API_KEY) {
    throw newEnvironmentVariableNotSetError('SECURITY_API_KEY')
  }
  if (process.env.SECURITY_MAX_FAILED_ATTEMPTS) {
    config.security.maxFailedAttempts = parseInt(
      process.env.SECURITY_MAX_FAILED_ATTEMPTS,
    )
  }
  if (process.env.SECURITY_LOCKOUT_PERIOD) {
    config.security.lockoutPeriod = parseInt(
      process.env.SECURITY_LOCKOUT_PERIOD,
    )
  }
  config.security.apiKey = process.env.SECURITY_API_KEY
}

export function readWebhooks(config: Config) {
  if (process.env.USER_WEBHOOKS) {
    config.userWebhooks = process.env.USER_WEBHOOKS.split(',')
    config.userWebhooks.forEach((e) => e.trim())
  }
}

function newEnvironmentVariableNotSetError(variable: string) {
  return new Error(`${variable} environment variable is not set.`)
}
