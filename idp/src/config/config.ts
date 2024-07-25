// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import * as process from 'node:process'
import { Config } from './types'

let config: Config

export function getConfig(): Config {
  if (!config) {
    config = new Config()
    config.port = parseInt(process.env.PORT)
    readURLs(config)
    readToken(config)
    readPassword(config)
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

export function readPassword(config: Config) {
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
