// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

export function getAccessToken() {
  const accessToken = loadAccessToken()
  const tokenExpiry = loadTokenExpiry()
  if (accessToken && tokenExpiry && new Date() < new Date(tokenExpiry)) {
    return accessToken
  } else {
    clearToken().then()
  }
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

export const ACCESS_TOKEN_COOKIE = 'voltaserve_access_token'

export async function clearToken() {
  document.cookie = `${ACCESS_TOKEN_COOKIE}=; Max-Age=-99999999;`
  removeAccessToken()
  removeRefreshToken()
  removeTokenExpiry()
}

export const ACCESS_TOKEN = 'voltaserve_access_token'

export function saveAccessToken(token: string) {
  return localStorage.setItem(ACCESS_TOKEN, token)
}

export function loadAccessToken(): string | null {
  return localStorage.getItem(ACCESS_TOKEN)
}

export function removeAccessToken() {
  return localStorage.removeItem(ACCESS_TOKEN)
}

export const REFRESH_TOKEN = 'voltaserve_refresh_token'

export function saveRefreshToken(token: string) {
  return localStorage.setItem(REFRESH_TOKEN, token)
}

export function loadRefreshToken(): string | null {
  return localStorage.getItem(REFRESH_TOKEN)
}

export function removeRefreshToken() {
  return localStorage.removeItem(REFRESH_TOKEN)
}

export const TOKEN_EXPIRY = 'voltaserve_token_expiry'

export function saveTokenExpiry(tokenExpiry: string) {
  return localStorage.setItem(TOKEN_EXPIRY, tokenExpiry)
}

export function loadTokenExpiry(): string | null {
  return localStorage.getItem(TOKEN_EXPIRY)
}

export function removeTokenExpiry() {
  return localStorage.removeItem(TOKEN_EXPIRY)
}
