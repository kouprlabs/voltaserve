// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

import { createWriteStream, unlink } from 'fs'
import { get } from 'http'
import { API_URL } from '@/config'
import { Token } from './idp'

export type APIErrorResponse = {
  code: string
  status: number
  message: string
  userMessage: string
  moreInfo: string
}

export class APIError extends Error {
  constructor(readonly error: APIErrorResponse) {
    super(JSON.stringify(error, null, 2))
  }
}

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
  permission: PermissionType
  isShared: boolean
  snapshot?: Snapshot
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

export type FileCreateOptions = {
  type: FileType
  workspaceId: string
  parentId?: string
  blob?: Blob
  name?: string
}

export type FilePatchOptions = {
  id: string
  blob: Blob
  name: string
}

export type FileMoveOptions = {
  ids: string[]
}

export class HealthAPI {
  async get(): Promise<string> {
    const response = await fetch(`${API_URL}/v2/health`, { method: 'GET' })
    return response.text()
  }
}

export class FileAPI {
  constructor(private token: Token) {}

  private async jsonResponseOrThrow<T>(response: Response): Promise<T> {
    if (response.headers.get('content-type')?.includes('application/json')) {
      const json = await response.json()
      if (response.status > 299) {
        throw new APIError(json)
      }
      return json
    } else {
      if (response.status > 299) {
        throw new Error(response.statusText)
      }
    }
  }

  async create({
    type,
    workspaceId,
    parentId,
    name,
    blob,
  }: FileCreateOptions): Promise<File> {
    const params = new URLSearchParams({ type, workspace_id: workspaceId })
    if (parentId) {
      params.append('parent_id', parentId)
    }
    if (name) {
      params.append('name', name)
    }
    if (type === FileType.File && blob) {
      return this.upload(`${API_URL}/v2/files?${params}`, 'POST', blob, name)
    } else if (type === FileType.Folder) {
      const response = await fetch(`${API_URL}/v2/files?${params}`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${this.token.access_token}`,
        },
      })
      return this.jsonResponseOrThrow(response)
    }
  }

  async patch({ id, blob, name }: FilePatchOptions): Promise<File> {
    return this.upload(`${API_URL}/v2/files/${id}`, 'PATCH', blob, name)
  }

  private async upload<T>(
    url: string,
    method: string,
    blob: Blob,
    name: string,
  ) {
    const formData = new FormData()
    formData.set('file', blob, name)
    const response = await fetch(url, {
      method,
      headers: {
        'Authorization': `Bearer ${this.token.access_token}`,
      },
      body: formData,
    })
    return this.jsonResponseOrThrow<T>(response)
  }

  async getByPath(path: string): Promise<File> {
    const response = await fetch(
      `${API_URL}/v2/files?path=${encodeURIComponent(path)}`,
      {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${this.token.access_token}`,
          'Content-Type': 'application/json',
        },
      },
    )
    return this.jsonResponseOrThrow(response)
  }

  async listByPath(path: string): Promise<File[]> {
    const response = await fetch(
      `${API_URL}/v2/files/list?path=${encodeURIComponent(path)}`,
      {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${this.token.access_token}`,
          'Content-Type': 'application/json',
        },
      },
    )
    return this.jsonResponseOrThrow(response)
  }

  async copy(id: string, options: FileCopyOptions): Promise<File[]> {
    const response = await fetch(`${API_URL}/v2/files/${id}/copy`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${this.token.access_token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        ids: options.ids,
      }),
    })
    return this.jsonResponseOrThrow(response)
  }

  async move(id: string, options: FileMoveOptions): Promise<void> {
    const response = await fetch(`${API_URL}/v2/files/${id}/move`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${this.token.access_token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        ids: options.ids,
      }),
    })
    return this.jsonResponseOrThrow(response)
  }

  async patchName(id: string, options: FileRenameOptions): Promise<File> {
    const response = await fetch(`${API_URL}/v2/files/${id}/name`, {
      method: 'PATCH',
      headers: {
        'Authorization': `Bearer ${this.token.access_token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        name: options.name,
      }),
    })
    return this.jsonResponseOrThrow(response)
  }

  async delete(id: string): Promise<void> {
    const response = await fetch(`${API_URL}/v2/files`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${this.token.access_token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ ids: [id] }),
    })
    return this.jsonResponseOrThrow(response)
  }

  downloadOriginal(file: File, outputPath: string): Promise<void> {
    return new Promise<void>((resolve, reject) => {
      const ws = createWriteStream(outputPath)
      const request = get(
        `${API_URL}/v2/files/${file.id}/original${file.snapshot.original.extension}?access_token=${this.token.access_token}`,
        (response) => {
          response.pipe(ws)
          ws.on('finish', () => {
            ws.close()
            resolve()
          })
        },
      )
      request.on('error', (error) => {
        unlink(outputPath, () => {
          reject(error)
        })
      })
    })
  }
}

export const VIEWER_PERMISSION = 'viewer'
export const EDITOR_PERMISSION = 'editor'
export const OWNER_PERMISSION = 'owner'

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
