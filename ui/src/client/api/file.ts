/* eslint-disable react-hooks/rules-of-hooks */
import useSWR from 'swr'
import { apiFetch, apiFetcher } from '@/client/fetch'
import { User } from '@/client/idp/user'
import { getConfig } from '@/config/config'
import { getAccessTokenOrRedirect } from '@/infra/token'
import { Group } from './group'
import { PermissionType } from './permission'
import { Download, Snapshot, Thumbnail } from './snapshot'

export enum FileType {
  File = 'file',
  Folder = 'folder',
}

export enum SortBy {
  Name = 'name',
  Kind = 'kind',
  Size = 'size',
  DateCreated = 'date_created',
  DateModified = 'date_modified',
}

export enum SortOrder {
  Asc = 'asc',
  Desc = 'desc',
}

export enum SnapshotStatus {
  New = 'new',
  Processing = 'processing',
  Ready = 'ready',
  Error = 'error',
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
  ocr?: Download
  thumbnail?: Thumbnail
  status: SnapshotStatus
  snapshots: Snapshot[]
  permission: PermissionType
  isShared: boolean
  createTime: string
  updateTime?: string
}

export type CreateFolderOptions = {
  workspaceId: string
  name: string
  parentId: string
}

export type List = {
  data: File[]
  totalPages: number
  totalElements: number
  page: number
  size: number
}

export type SearchResult = {
  request: SearchOptions
} & List

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

export type SearchOptions = {
  text: string
  workspaceId: string
  parentId?: string
  type?: string
  createTimeAfter?: number
  createTimeBefore?: number
  updateTimeAfter?: number
  updateTimeBefore?: number
}

export type ListOptions = {
  size?: number
  page?: number
  type?: FileType
  sortBy?: SortBy
  sortOrder?: SortOrder
}

export type MoveOptions = {
  ids: string[]
}

export type CopyOptions = {
  ids: string[]
}

export type BatchDeleteOptions = {
  ids: string[]
}

export type BatchGetOptions = {
  ids: string[]
}

export type RenameOptions = {
  name: string
}

export type UpdateOcrLanguageOptions = {
  id: string
}

export type GrantUserPermissionOptions = {
  ids: string[]
  userId: string
  permission: string
}

export type RevokeUserPermissionOptions = {
  ids: string[]
  userId: string
}

export type GrantGroupPermissionOptions = {
  ids: string[]
  groupId: string
  permission: string
}

export type RevokeGroupPermissionOptions = {
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
    onProgress?: (value: number) => void,
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
      onProgress,
    )
  }

  static async patch(
    id: string,
    request: XMLHttpRequest,
    file: Blob,
    onProgress?: (value: number) => void,
  ): Promise<File> {
    return this.doUpload(
      `${getConfig().apiURL}/files/${id}`,
      'PATCH',
      request,
      file,
      onProgress,
    )
  }

  private static async doUpload(
    url: string,
    method: string,
    request: XMLHttpRequest,
    file: Blob,
    onProgress?: (value: number) => void,
  ) {
    const formData = new FormData()
    formData.append('file', file)
    return new Promise<File>((resolve, reject) => {
      request.open(method, url)
      request.setRequestHeader(
        'Authorization',
        `Bearer ${getAccessTokenOrRedirect()}`,
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

  static async createFolder(options: CreateFolderOptions): Promise<File> {
    return apiFetch('/files/create_folder', {
      method: 'POST',
      body: JSON.stringify(options),
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    }).then((result) => result.json())
  }

  static async list(id: string, options: ListOptions): Promise<List> {
    return apiFetch(
      `/files/${id}/list?${this.paramsFromListOptions(options)}`,
      {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
          'Content-Type': 'application/json',
        },
      },
    ).then((result) => result.json())
  }

  static async search(
    options: SearchOptions,
    size: number,
    page: number,
  ): Promise<List> {
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
      },
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

  static async getIds(id: string): Promise<string[]> {
    return apiFetch(`/files/${id}/get_ids`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
      },
    }).then(async (result) => await result.json())
  }

  static async rename(id: string, options: RenameOptions): Promise<File> {
    return apiFetch(`/files/${id}/rename`, {
      method: 'POST',
      body: JSON.stringify(options),
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    }).then((result) => result.json())
  }

  static async updateOcrLanguage(
    id: string,
    options: UpdateOcrLanguageOptions,
  ): Promise<File> {
    return apiFetch(`/files/${id}/update_ocr_language`, {
      method: 'POST',
      body: JSON.stringify(options),
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    }).then((result) => result.json())
  }

  static async deleteOcr(id: string): Promise<File> {
    return apiFetch(`/files/${id}/delete_ocr`, {
      method: 'POST',
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

  static async batchDelete(options: BatchDeleteOptions) {
    return apiFetch(`/files/batch_delete`, {
      method: 'POST',
      body: JSON.stringify(options),
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    })
  }

  static async move(id: string, options: MoveOptions) {
    return apiFetch(`/files/${id}/move`, {
      method: 'POST',
      body: JSON.stringify(options),
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    })
  }

  static async copy(id: string, options: CopyOptions) {
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

  static async batchGet(options: BatchGetOptions): Promise<File[]> {
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
      swrOptions,
    )
  }

  static async grantUserPermission(options: GrantUserPermissionOptions) {
    return apiFetch(`/files/grant_user_permission`, {
      method: 'POST',
      body: JSON.stringify(options),
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    })
  }

  static async revokeUserPermission(options: RevokeUserPermissionOptions) {
    return apiFetch(`/files/revoke_user_permission`, {
      method: 'POST',
      body: JSON.stringify(options),
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    })
  }

  static async grantGroupPermission(options: GrantGroupPermissionOptions) {
    return apiFetch(`/files/grant_group_permission`, {
      method: 'POST',
      body: JSON.stringify(options),
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    })
  }

  static async revokeGroupPermission(options: RevokeGroupPermissionOptions) {
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

  static paramsFromListOptions(options?: ListOptions): URLSearchParams {
    const params: any = {}
    if (options?.page) {
      params.page = options.page.toString()
    }
    if (options?.size) {
      params.size = options.size.toString()
    }
    if (options?.sortBy) {
      params.sort_by = options.sortBy.toString()
    }
    if (options?.sortOrder) {
      params.sort_order = options.sortOrder.toString()
    }
    if (options?.type) {
      params.type = options.type
    }
    return new URLSearchParams(params)
  }
}
