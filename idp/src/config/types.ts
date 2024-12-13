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
  port: number = 0
  publicUIURL: string = ''
  databaseURL: string = ''
  token: TokenConfig
  password: PasswordConfig
  security: SecurityConfig
  corsOrigins: string[] = []
  search: SearchConfig
  smtp: SMTPConfig

  constructor() {
    this.token = new TokenConfig()
    this.password = new PasswordConfig()
    this.search = new SearchConfig()
    this.smtp = new SMTPConfig()
    this.security = new SecurityConfig()
  }
}

export class TokenConfig {
  jwtSigningKey: string = ''
  audience: string = ''
  issuer: string = ''
  accessTokenLifetime: number = 0
  refreshTokenLifetime: number = 0
}

export class PasswordConfig {
  minLength: number = 0
  minLowercase: number = 0
  minUppercase: number = 0
  minNumbers: number = 0
  minSymbols: number = 0
}

export class SecurityConfig {
  maxFailedAttempts: number = 0
  lockoutPeriod: number = 0
}

export class SearchConfig {
  url: string = ''
}

export class SMTPConfig {
  host: string = ''
  port: number = 0
  secure: boolean = false
  username?: string
  password?: string
  senderAddress: string = ''
  senderName: string = ''
}
