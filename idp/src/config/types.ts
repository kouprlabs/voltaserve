// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

export class Config {
  port: number
  publicUIURL: string
  databaseURL: string
  token: TokenConfig
  password: PasswordConfig
  corsOrigins: string[]
  search: SearchConfig
  smtp: SMTPConfig

  constructor() {
    this.token = new TokenConfig()
    this.password = new PasswordConfig()
    this.search = new SearchConfig()
    this.smtp = new SMTPConfig()
  }
}

export class TokenConfig {
  jwtSigningKey: string
  audience: string
  issuer: string
  accessTokenLifetime: number
  refreshTokenLifetime: number
}

export class PasswordConfig {
  minLength: number
  minLowercase: number
  minUppercase: number
  minNumbers: number
  minSymbols: number
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
