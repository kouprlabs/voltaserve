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
import TokenAPI, { Token } from '@/client/idp/token'
import {
  loadAccessToken,
  loadRefreshToken,
  loadTokenExpiry,
  removeAccessToken,
  removeRefreshToken,
  removeTokenExpiry,
  saveAccessToken,
  saveRefreshToken,
  saveTokenExpiry,
} from '@/local-storage'

export const COOKIE = 'voltaserve_access_token'

export async function saveToken(token: Token) {
  document.cookie = `${COOKIE}=${token.access_token}; Path=/; Max-Age=${token.expires_in}`
  saveAccessToken(token.access_token)
  saveRefreshToken(token.refresh_token)
  const tokenExpiry = new Date()
  tokenExpiry.setSeconds(tokenExpiry.getSeconds() + token.expires_in)
  saveTokenExpiry(tokenExpiry.toISOString())
}

export async function clearToken() {
  document.cookie = `${COOKIE}=; Max-Age=-99999999;`
  removeAccessToken()
  removeRefreshToken()
  removeTokenExpiry()
}

export function getAccessTokenOrRedirect(): string {
  const accessToken = getAccessToken()
  if (accessToken) {
    return accessToken
  } else {
    window.location.href = '/sign-in'
    return ''
  }
}

export function getAccessToken() {
  const accessToken = loadAccessToken()
  const tokenExpiry = loadTokenExpiry()
  if (accessToken && tokenExpiry && new Date() < new Date(tokenExpiry)) {
    return accessToken
  } else {
    clearToken().then()
  }
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
