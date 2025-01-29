// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { decodeJwt } from 'jose'
import { TokenAPI, Token } from '@/client/idp/token'
import {
  ACCESS_TOKEN_COOKIE,
  getAccessToken,
  loadRefreshToken,
  loadTokenExpiry,
  saveAccessToken,
  saveRefreshToken,
  saveTokenExpiry,
} from '@/client/token'

export async function saveToken(token: Token) {
  document.cookie = `${ACCESS_TOKEN_COOKIE}=${token.access_token}; Path=/; Max-Age=${token.expires_in}`
  saveAccessToken(token.access_token)
  saveRefreshToken(token.refresh_token)
  const tokenExpiry = new Date()
  tokenExpiry.setSeconds(tokenExpiry.getSeconds() + token.expires_in)
  saveTokenExpiry(tokenExpiry.toISOString())
}

export function getAdminStatus(): boolean {
  const accessToken = getAccessToken()
  if (accessToken) {
    return !!decodeJwt(accessToken).is_admin
  } else {
    return false
  }
}

export function getUserId(): string {
  const accessToken = getAccessToken()
  if (accessToken) {
    const userId = decodeJwt(accessToken).sub
    return userId ? userId : ''
  } else {
    return ''
  }
}

setInterval(async () => {
  const refreshToken = loadRefreshToken()
  const tokenExpiry = loadTokenExpiry()
  if (tokenExpiry && refreshToken) {
    const earlyExpiry = new Date(tokenExpiry)
    earlyExpiry.setMinutes(earlyExpiry.getMinutes() - 1)
    if (new Date() >= earlyExpiry) {
      const token = await TokenAPI.exchange({
        grant_type: 'refresh_token',
        refresh_token: refreshToken,
      })
      saveToken(token).then()
    }
  }
}, 5000)
