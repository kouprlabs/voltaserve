export const permissionViewer: string = 'viewer'
export const permissionEditor: string = 'editor'
export const permissionOwner: string = 'owner'

export type PermissionType = 'viewer' | 'editor' | 'owner'

export function geViewerPermission(permission: string): boolean {
  return (
    getPermissionWeight(permission) >= getPermissionWeight(permissionViewer)
  )
}

export function geEditorPermission(permission: string) {
  return (
    getPermissionWeight(permission) >= getPermissionWeight(permissionEditor)
  )
}

export function geOwnerPermission(permission: string) {
  return getPermissionWeight(permission) >= getPermissionWeight(permissionOwner)
}

export function ltViewerPermission(permission: string): boolean {
  return getPermissionWeight(permission) < getPermissionWeight(permissionViewer)
}

export function ltEditorPermission(permission: string) {
  return getPermissionWeight(permission) < getPermissionWeight(permissionEditor)
}

export function ltOwnerPermission(permission: string) {
  return getPermissionWeight(permission) < getPermissionWeight(permissionOwner)
}

export function getPermissionWeight(permission: string) {
  switch (permission) {
    case permissionViewer:
      return 1
    case permissionEditor:
      return 2
    case permissionOwner:
      return 3
    default:
      return 0
  }
}
