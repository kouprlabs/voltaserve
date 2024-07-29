// Copyright 2024 Mateusz Kaźmierczak.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import {
  loadAdminAccessToken,
  removeAdminAccessToken,
  saveAdminAccessToken,
} from '@/local-storage'
import store from '@/store/configure-store'
import { errorOccurred } from '@/store/ui/error'

export const COOKIE = 'voltaserve_admin_access_token'

export function saveAdminToken(token: string) {
  try {
    document.cookie = `${COOKIE}=${token}; Path=/; Max-Age=${new Date((JSON.parse(atob(token.split('.')[1])).exp * 1000).toString()).getSeconds()}`
    saveAdminAccessToken(`${token}`)
  } catch (e) {
    clearAdminToken()
  }
}

export function clearAdminToken() {
  document.cookie = `${COOKIE}=; Max-Age=-99999999;`
  removeAdminAccessToken()
}

// export function getAdminAccessTokenOrRedirect(): string {
//   const accessToken = getAccessToken()
//   if (accessToken) {
//     return accessToken
//   } else {
//     window.location.href = '/sign-in'
//     return ''
//   }
// }

export function getAdminAccessToken() {
  const accessToken = loadAdminAccessToken()
  let tokenExpiry
  try {
    if (accessToken) {
      tokenExpiry = JSON.parse(atob(accessToken.split('.')[1])).exp * 1000
    }
    if (accessToken && tokenExpiry && new Date() < new Date(tokenExpiry)) {
      return accessToken
    } else {
      store.dispatch(errorOccurred('Admin token has expired or is invalid'))
      clearAdminToken()
      return undefined
    }
  } catch (e) {
    store.dispatch(errorOccurred('Admin token has expired or is invalid'))
    clearAdminToken()
  }
}

// setInterval(async () => {
//   const resp = await AdminApi.adminAuthenticate()
//   if (resp.detail !== 'Hello, it is root of admin microservice!') {
//     clearAdminToken()
//     window.location.href = '/'
//   }
// }, 60000)
