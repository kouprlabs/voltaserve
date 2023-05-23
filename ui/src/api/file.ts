/* eslint-disable react-hooks/rules-of-hooks */
import useSWR from 'swr'
import { getConfig } from '@/config/config'
import { getAccessTokenOrRedirect } from '@/infra/token'
import { apiFetch, apiFetcher } from './fetch'
import { Group } from './group'
import { PermissionType } from './permission'
import { Download, Snapshot, Thumbnail } from './snapshot'
import { User } from './user'

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

export type FileCreateFolderOptions = {
  workspaceId: string
  name: string
  parentId: string
}

export type FileList = {
  data: File[]
  totalPages: number
  totalElements: number
  page: number
  size: number
}

export type FileSearchResult = {
  request: FileSearchOptions
} & FileList

export type UserPermission = {
  id: string
  user: User
  permission: string
}

export type GroupPermission = {
  id: string
  group: Group
  permission: string
}

export type FileSearchOptions = {
  text: string
  workspaceId: string
  parentId?: string
  type?: string
  createTimeAfter?: number
  createTimeBefore?: number
  updateTimeAfter?: number
  updateTimeBefore?: number
}

export type FileMoveOptions = {
  ids: string[]
}

export type FileCopyOptions = {
  ids: string[]
}

export type FileBatchDeleteOptions = {
  ids: string[]
}

export type FileBatchGetOptions = {
  ids: string[]
}

export type FileRenameOptions = {
  name: string
}

export type FileGrantUserPermissionOptions = {
  ids: string[]
  userId: string
  permission: string
}

export type FileRevokeUserPermissionOptions = {
  ids: string[]
  userId: string
}

export type FileGrantGroupPermissionOptions = {
  ids: string[]
  groupId: string
  permission: string
}

export type FileRevokeGroupPermissionOptions = {
  ids: string[]
  groupId: string
}

export default class FileAPI {
  public static DEFAULT_PAGE_SIZE = 42

  static async upload(
    workspaceId: string,
    parentId: string | null,
    request: XMLHttpRequest,
    file: Blob,
    onProgress?: (value: number) => void
  ): Promise<File> {
    const params = new URLSearchParams({ workspace_id: workspaceId })
    if (parentId) {
      params.append('parent_id', parentId)
    }
    return this.doUpload(
      `${getConfig().apiURL}/files?${params}`,
      'POST',
      request,
      file,
      onProgress
    )
  }

  static async patch(
    id: string,
    request: XMLHttpRequest,
    file: Blob,
    onProgress?: (value: number) => void
  ): Promise<File> {
    return this.doUpload(
      `${getConfig().apiURL}/files/${id}`,
      'PATCH',
      request,
      file,
      onProgress
    )
  }

  private static async doUpload(
    url: string,
    method: string,
    request: XMLHttpRequest,
    file: Blob,
    onProgress?: (value: number) => void
  ) {
    const formData = new FormData()
    formData.append('file', file)
    return new Promise<File>((resolve, reject) => {
      request.open(method, url)
      request.setRequestHeader(
        'Authorization',
        `Bearer ${getAccessTokenOrRedirect()}`
      )
      request.onloadend = () => {
        if (request.status <= 299) {
          try {
            resolve(JSON.parse(request.responseText))
          } catch (e) {
            reject(e)
          }
        } else {
          try {
            reject(JSON.parse(request.responseText))
          } catch {
            reject(request.responseText)
          }
        }
      }
      request.upload.onprogress = (e) => {
        onProgress?.((e.loaded / e.total) * 100)
      }
      request.send(formData)
    })
  }

  static async createFolder(options: FileCreateFolderOptions): Promise<File> {
    return apiFetch('/files/create_folder', {
      method: 'POST',
      body: JSON.stringify(options),
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    }).then((result) => result.json())
  }

  static async list(
    id: string,
    size: number,
    page: number,
    type?: FileType
  ): Promise<FileList> {
    const params: any = {
      page: page.toString(),
      size: size.toString(),
    }
    if (type) {
      params.type = type
    }
    return apiFetch(`/files/${id}/list?${new URLSearchParams(params)}`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    }).then((result) => result.json())
  }

  static async search(
    options: FileSearchOptions,
    size: number,
    page: number
  ): Promise<FileList> {
    return apiFetch(
      `/files/search?${new URLSearchParams({
        page: page.toString(),
        size: size.toString(),
      })}`,
      {
        method: 'POST',
        body: JSON.stringify(options),
        headers: {
          'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
          'Content-Type': 'application/json',
        },
      }
    ).then((result) => result.json())
  }

  static async getPath(id: string): Promise<File[]> {
    return apiFetch(`/files/${id}/get_path`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
      },
    }).then(async (result) => await result.json())
  }

  static async rename(id: string, options: FileRenameOptions): Promise<File> {
    return apiFetch(`/files/${id}/rename`, {
      method: 'POST',
      body: JSON.stringify(options),
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    }).then((result) => result.json())
  }

  static async delete(id: string) {
    return apiFetch(`/files/${id}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    })
  }

  static async batchDelete(options: FileBatchDeleteOptions) {
    return apiFetch(`/files/batch_delete`, {
      method: 'POST',
      body: JSON.stringify(options),
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    })
  }

  static async move(id: string, options: FileMoveOptions) {
    return apiFetch(`/files/${id}/move`, {
      method: 'POST',
      body: JSON.stringify(options),
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    })
  }

  static async copy(id: string, options: FileCopyOptions) {
    return apiFetch(`/files/${id}/copy`, {
      method: 'POST',
      body: JSON.stringify(options),
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    })
  }

  static useGetById(id: string, swrOptions?: any) {
    return useSWR<File>(id ? `/files/${id}` : null, apiFetcher, swrOptions)
  }

  static async getById(id: string): Promise<File> {
    return apiFetch(`/files/${id}`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    }).then((result) => result.json())
  }

  static async batchGet(options: FileBatchGetOptions): Promise<File[]> {
    return apiFetch(`/files/batch_get`, {
      method: 'POST',
      body: JSON.stringify(options),
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    }).then((result) => result.json())
  }

  static async getItemCount(id: string): Promise<number> {
    return apiFetch(`/files/${id}/get_item_count`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    }).then((result) => result.json())
  }

  static useGetItemCount(id: string, swrOptions?: any) {
    return useSWR<number>(
      id ? `/files/${id}/get_item_count` : null,
      apiFetcher,
      swrOptions
    )
  }

  static async grantUserPermission(options: FileGrantUserPermissionOptions) {
    return apiFetch(`/files/grant_user_permission`, {
      method: 'POST',
      body: JSON.stringify(options),
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    })
  }

  static async revokeUserPermission(options: FileRevokeUserPermissionOptions) {
    return apiFetch(`/files/revoke_user_permission`, {
      method: 'POST',
      body: JSON.stringify(options),
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    })
  }

  static async grantGroupPermission(options: FileGrantGroupPermissionOptions) {
    return apiFetch(`/files/grant_group_permission`, {
      method: 'POST',
      body: JSON.stringify(options),
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    })
  }

  static async revokeGroupPermission(
    options: FileRevokeGroupPermissionOptions
  ) {
    return apiFetch(`/files/revoke_group_permission`, {
      method: 'POST',
      body: JSON.stringify(options),
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    })
  }

  static async getUserPermissions(id: string): Promise<UserPermission[]> {
    return apiFetch(`/files/${id}/get_user_permissions`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    }).then((result) => result.json())
  }

  static async getGroupPermissions(id: string): Promise<GroupPermission[]> {
    return apiFetch(`/files/${id}/get_group_permissions`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    }).then((result) => result.json())
  }
}
