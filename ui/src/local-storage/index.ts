// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { FileSortBy, FileSortOrder } from '@/client/api/file'
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

export function loadFileSortBy(): FileSortBy | null {
  const value = localStorage.getItem(SORT_BY_KEY)
  if (value) {
    return JSON.parse(value)
  } else {
    return null
  }
}

export function saveFileSortBy(sortBy: FileSortBy) {
  return localStorage.setItem(SORT_BY_KEY, JSON.stringify(sortBy))
}

export const SORT_ORDER_KEY = 'voltaserve_file_sort_order'

export function loadFileSortOrder(): FileSortOrder | null {
  const value = localStorage.getItem(SORT_ORDER_KEY)
  if (value) {
    return JSON.parse(value)
  } else {
    return null
  }
}

export function saveFileSortOrder(sortOrder: FileSortOrder) {
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

const THEME = 'voltaserve_theme'

export type ThemeValue = 'light' | 'dark' | 'system'

export function loadTheme(): ThemeValue {
  return localStorage.getItem(THEME) as ThemeValue
}

export function saveTheme(theme: ThemeValue) {
  localStorage.setItem(THEME, theme)
}
