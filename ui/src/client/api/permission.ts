// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

export const NONE_PERMISSION = 'none'
export const VIEWER_PERMISSION = 'viewer'
export const EDITOR_PERMISSION = 'editor'
export const OWNER_PERMISSION = 'owner'

export type PermissionType = 'viewer' | 'editor' | 'owner'

export function gtViewerPermission(permission: string): boolean {
  return (
    getPermissionWeight(permission) > getPermissionWeight(VIEWER_PERMISSION)
  )
}

export function gtEditorPermission(permission: string) {
  return (
    getPermissionWeight(permission) > getPermissionWeight(EDITOR_PERMISSION)
  )
}

export function gtOwnerPermission(permission: string) {
  return getPermissionWeight(permission) > getPermissionWeight(OWNER_PERMISSION)
}

export function geViewerPermission(permission: string): boolean {
  return (
    getPermissionWeight(permission) >= getPermissionWeight(VIEWER_PERMISSION)
  )
}

export function geEditorPermission(permission: string) {
  return (
    getPermissionWeight(permission) >= getPermissionWeight(EDITOR_PERMISSION)
  )
}

export function geOwnerPermission(permission: string) {
  return (
    getPermissionWeight(permission) >= getPermissionWeight(OWNER_PERMISSION)
  )
}

export function ltViewerPermission(permission: string): boolean {
  return (
    getPermissionWeight(permission) < getPermissionWeight(VIEWER_PERMISSION)
  )
}

export function ltEditorPermission(permission: string) {
  return (
    getPermissionWeight(permission) < getPermissionWeight(EDITOR_PERMISSION)
  )
}

export function ltOwnerPermission(permission: string) {
  return getPermissionWeight(permission) < getPermissionWeight(OWNER_PERMISSION)
}

export function leViewerPermission(permission: string): boolean {
  return (
    getPermissionWeight(permission) <= getPermissionWeight(VIEWER_PERMISSION)
  )
}

export function leEditorPermission(permission: string) {
  return (
    getPermissionWeight(permission) <= getPermissionWeight(EDITOR_PERMISSION)
  )
}

export function leOwnerPermission(permission: string) {
  return (
    getPermissionWeight(permission) <= getPermissionWeight(OWNER_PERMISSION)
  )
}

export function getPermissionWeight(permission: string) {
  switch (permission) {
    case VIEWER_PERMISSION:
      return 1
    case EDITOR_PERMISSION:
      return 2
    case OWNER_PERMISSION:
      return 3
    default:
      return 0
  }
}
