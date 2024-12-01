// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

export const NONE_PERMISSION = 'none'
export const VIEWER_PERMISSION = 'viewer'
export const EDITOR_PERMISSION = 'editor'
export const OWNER_PERMISSION = 'owner'

export type PermissionType = 'viewer' | 'editor' | 'owner' | 'none'

export function gtViewerPermission(permission: PermissionType): boolean {
  return getPermissionWeight(permission) > getPermissionWeight(VIEWER_PERMISSION)
}

export function gtEditorPermission(permission: PermissionType) {
  return getPermissionWeight(permission) > getPermissionWeight(EDITOR_PERMISSION)
}

export function gtOwnerPermission(permission: PermissionType) {
  return getPermissionWeight(permission) > getPermissionWeight(OWNER_PERMISSION)
}

export function geViewerPermission(permission: PermissionType): boolean {
  return getPermissionWeight(permission) >= getPermissionWeight(VIEWER_PERMISSION)
}

export function geEditorPermission(permission: PermissionType) {
  return getPermissionWeight(permission) >= getPermissionWeight(EDITOR_PERMISSION)
}

export function geOwnerPermission(permission: PermissionType) {
  return getPermissionWeight(permission) >= getPermissionWeight(OWNER_PERMISSION)
}

export function ltViewerPermission(permission: PermissionType): boolean {
  return getPermissionWeight(permission) < getPermissionWeight(VIEWER_PERMISSION)
}

export function ltEditorPermission(permission: PermissionType) {
  return getPermissionWeight(permission) < getPermissionWeight(EDITOR_PERMISSION)
}

export function ltOwnerPermission(permission: PermissionType) {
  return getPermissionWeight(permission) < getPermissionWeight(OWNER_PERMISSION)
}

export function leViewerPermission(permission: PermissionType): boolean {
  return getPermissionWeight(permission) <= getPermissionWeight(VIEWER_PERMISSION)
}

export function leEditorPermission(permission: PermissionType) {
  return getPermissionWeight(permission) <= getPermissionWeight(EDITOR_PERMISSION)
}

export function leOwnerPermission(permission: PermissionType) {
  return getPermissionWeight(permission) <= getPermissionWeight(OWNER_PERMISSION)
}

export function getPermissionWeight(permission: PermissionType) {
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
