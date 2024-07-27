// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { SortBy, SortOrder } from '@/client/api/types/queries'
import { FileViewType } from '@/types/file'

const ICON_SCALE_KEY = 'voltaserve_file_icon_scale'

export function loadIconScale(): number | null {
  const value = localStorage.getItem(ICON_SCALE_KEY)
  if (value) {
    return JSON.parse(value)
  } else {
    return null
  }
}

export function saveIconScale(scale: number) {
  return localStorage.setItem(ICON_SCALE_KEY, JSON.stringify(scale))
}

export const SORT_BY_KEY = 'voltaserve_file_sort_by'

export function loadFileSortBy(): SortBy | null {
  const value = localStorage.getItem(SORT_BY_KEY)
  if (value) {
    return JSON.parse(value)
  } else {
    return null
  }
}

export function saveFileSortBy(sortBy: SortBy) {
  return localStorage.setItem(SORT_BY_KEY, JSON.stringify(sortBy))
}

export const SORT_ORDER_KEY = 'voltaserve_file_sort_order'

export function loadFileSortOrder(): SortOrder | null {
  const value = localStorage.getItem(SORT_ORDER_KEY)
  if (value) {
    return JSON.parse(value)
  } else {
    return null
  }
}

export function saveFileSortOrder(sortOrder: SortOrder) {
  return localStorage.setItem(SORT_ORDER_KEY, JSON.stringify(sortOrder))
}

export const VIEW_TYPE_KEY = 'voltaserve_file_view_type'

export function loadFileViewType(): FileViewType | null {
  const value = localStorage.getItem(VIEW_TYPE_KEY)
  if (value) {
    return JSON.parse(value)
  } else {
    return null
  }
}

export function saveFileViewType(viewType: FileViewType) {
  return localStorage.setItem(VIEW_TYPE_KEY, JSON.stringify(viewType))
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
