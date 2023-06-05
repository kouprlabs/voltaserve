import { API_URL } from '@/config'
import { Token } from './idp'
import { IncomingMessage, get } from 'http'

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

export type PermissionType = 'viewer' | 'editor' | 'owner'

export type Snapshot = {
  version: number
  original: Download
  preview?: Download
  ocr?: Download
  text?: Download
  thumbnail?: Thumbnail
}

export type Download = {
  extension: string
  size: number
  image: ImageProps | undefined
}

export type ImageProps = {
  width: number
  height: number
}

export type Thumbnail = {
  base64: string
  width: number
  height: number
}

export type FileCopyOptions = {
  ids: string[]
}

export type FileRenameOptions = {
  name: string
}

export type FileCreateFolderOptions = {
  workspaceId: string
  name: string
  parentId: string
}

export type FileUploadOptions = {
  workspaceId: string
  parentId: string | null
  blob: Blob
  name: string
}

export type FileMoveOptions = {
  ids: string[]
}

export class FileAPI {
  constructor(private token: Token) {}

  async upload(options: FileUploadOptions) {
    const params = new URLSearchParams({
      workspace_id: options.workspaceId,
    })
    if (options.parentId) {
      params.append('parent_id', options.parentId)
    }
    const formData = new FormData()
    formData.set('file', options.blob, options.name)
    await fetch(`${API_URL}/v1/files?${params}`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${this.token.access_token}`,
      },
      body: formData,
    })
  }

  async getByPath(path: string): Promise<File> {
    const result = await fetch(`${API_URL}/v1/files/get?path=${path}`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${this.token.access_token}`,
        'Content-Type': 'application/json',
      },
    })
    return result.json()
  }

  async listByPath(path: string): Promise<File[]> {
    const result = await fetch(`${API_URL}/v1/files/list?path=${path}`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${this.token.access_token}`,
        'Content-Type': 'application/json',
      },
    })
    return result.json()
  }

  async createFolder(options: FileCreateFolderOptions): Promise<void> {
    await fetch(`${API_URL}/v1/files/create_folder`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${this.token.access_token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        workspaceId: options.workspaceId,
        parentId: options.parentId,
        name: options.name,
      }),
    })
  }

  async copy(id: string, options: FileCopyOptions): Promise<File[]> {
    const result = await fetch(`${API_URL}/v1/files/${id}/copy`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${this.token.access_token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        ids: options.ids,
      }),
    })
    return result.json()
  }

  async move(id: string, options: FileMoveOptions) {
    await fetch(`${API_URL}/v1/files/${id}/move`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${this.token.access_token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        ids: options.ids,
      }),
    })
  }

  async rename(id: string, options: FileRenameOptions): Promise<File> {
    const result = await fetch(`${API_URL}/v1/files/${id}/rename`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${this.token.access_token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        name: options.name,
      }),
    })
    return result.json()
  }

  async delete(id: string): Promise<void> {
    await fetch(`${API_URL}/v1/files/${id}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${this.token.access_token}`,
        'Content-Type': 'application/json',
      },
    })
  }

  downloadOriginal(file: File, callback: (response: IncomingMessage) => void) {
    get(
      `${API_URL}/v1/files/${file.id}/original${file.original.extension}?access_token=${this.token.access_token}`,
      callback
    )
  }
}
