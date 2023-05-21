import { PermissionType } from './permission'
import { Download, Thumbnail, Snapshot } from './snapshot'

export enum FileType {
  File = 'file',
  Folder = 'folder',
}

export type File = {
  id: string
  workspaceId: string
  name: string
  type: FileType
  parentId: string
  version: number
  original?: Download
  preview?: Download
  thumbnail?: Thumbnail
  snapshots: Snapshot[]
  permission: PermissionType
  isShared: boolean
  createTime: string
  updateTime?: string
}
