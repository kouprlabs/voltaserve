/* eslint-disable react-hooks/rules-of-hooks */
import useSWR from 'swr'
import { apiFetcher } from '@/client/fetcher'
import { User } from '@/client/idp/user'
import { getConfig } from '@/config/config'
import { encodeQuery } from '@/helpers/query'
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
  thumbnail?: Thumbnail
  status: SnapshotStatus
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
  query?: Query
}

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

export type Query = {
  text: string
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
  query?: Query
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

export type UploadOptions = {
  workspaceId: string
  parentId?: string
  name?: string
  blob: Blob
  request: XMLHttpRequest
  onProgress?: (value: number) => void
}

export type PatchOptions = {
  id: string
  request: XMLHttpRequest
  blob: Blob
  onProgress?: (value: number) => void
}

export default class FileAPI {
  static async upload({
    workspaceId,
    parentId,
    name,
    request,
    blob,
    onProgress,
  }: UploadOptions): Promise<File> {
    const params = new URLSearchParams({ workspace_id: workspaceId })
    if (parentId) {
      params.append('parent_id', parentId)
    }
    if (name) {
      params.append('name', name)
    }
    return this.doUpload(
      `${getConfig().apiURL}/files?${params}`,
      'POST',
      request,
      blob,
      onProgress,
    )
  }

  static async patch({
    id,
    request,
    blob,
    onProgress,
  }: PatchOptions): Promise<File> {
    return this.doUpload(
      `${getConfig().apiURL}/files/${id}`,
      'PATCH',
      request,
      blob,
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

  static async createFolder(options: CreateFolderOptions) {
    return apiFetcher({
      url: '/files/create_folder',
      method: 'POST',
      body: JSON.stringify(options),
    }) as Promise<File>
  }

  static async list(id: string, options: ListOptions) {
    return apiFetcher({
      url: `/files/${id}/list?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<List>
  }

  static useList(
    id: string | undefined,
    options: ListOptions,
    swrOptions?: any,
  ) {
    const url = `/files/${id}/list?${this.paramsFromListOptions(options)}`
    return useSWR<List | undefined>(
      id ? url : null,
      () => apiFetcher({ url, method: 'GET' }),
      swrOptions,
    )
  }

  static async getPath(id: string) {
    return apiFetcher({
      url: `/files/${id}/get_path`,
      method: 'GET',
    }) as Promise<File[]>
  }

  static useGetPath(id: string | null | undefined, swrOptions?: any) {
    const url = `/files/${id}/get_path`
    return useSWR<File[]>(
      id ? url : null,
      () => apiFetcher({ url, method: 'GET' }) as Promise<File[]>,
      swrOptions,
    )
  }

  static async getIds(id: string) {
    return apiFetcher({
      url: `/files/${id}/get_ids`,
      method: 'GET',
    }) as Promise<string[]>
  }

  static async rename(id: string, options: RenameOptions) {
    return apiFetcher({
      url: `/files/${id}/rename`,
      method: 'POST',
      body: JSON.stringify(options),
    }) as Promise<File>
  }

  static async delete(id: string) {
    return apiFetcher({
      url: `/files/${id}`,
      method: 'DELETE',
    })
  }

  static async batchDelete(options: BatchDeleteOptions) {
    return apiFetcher({
      url: `/files/batch_delete`,
      method: 'POST',
      body: JSON.stringify(options),
    })
  }

  static async move(id: string, options: MoveOptions) {
    return apiFetcher({
      url: `/files/${id}/move`,
      method: 'POST',
      body: JSON.stringify(options),
    })
  }

  static async copy(id: string, options: CopyOptions) {
    return apiFetcher({
      url: `/files/${id}/copy`,
      method: 'POST',
      body: JSON.stringify(options),
    })
  }

  static useGetById(id: string | null | undefined, swrOptions?: any) {
    const url = `/files/${id}`
    return useSWR(
      id ? url : null,
      () => apiFetcher({ url, method: 'GET' }) as Promise<File>,
      swrOptions,
    )
  }

  static async getById(id: string) {
    return apiFetcher({
      url: `/files/${id}`,
      method: 'GET',
    }) as Promise<File>
  }

  static async batchGet(options: BatchGetOptions) {
    return apiFetcher({
      url: `/files/batch_get`,
      method: 'POST',
      body: JSON.stringify(options),
    }) as Promise<File[]>
  }

  static async getItemCount(id: string) {
    return apiFetcher({
      url: `/files/${id}/get_item_count`,
      method: 'GET',
    }) as Promise<number>
  }

  static useGetItemCount(id: string | null | undefined, swrOptions?: any) {
    const url = `/files/${id}/get_item_count`
    return useSWR<number>(
      id ? url : null,
      () => apiFetcher({ url, method: 'GET' }) as Promise<number>,
      swrOptions,
    )
  }

  static async grantUserPermission(options: GrantUserPermissionOptions) {
    return apiFetcher({
      url: `/files/grant_user_permission`,
      method: 'POST',
      body: JSON.stringify(options),
    })
  }

  static async revokeUserPermission(options: RevokeUserPermissionOptions) {
    return apiFetcher({
      url: `/files/revoke_user_permission`,
      method: 'POST',
      body: JSON.stringify(options),
    })
  }

  static async grantGroupPermission(options: GrantGroupPermissionOptions) {
    return apiFetcher({
      url: `/files/grant_group_permission`,
      method: 'POST',
      body: JSON.stringify(options),
    })
  }

  static async revokeGroupPermission(options: RevokeGroupPermissionOptions) {
    return apiFetcher({
      url: `/files/revoke_group_permission`,
      method: 'POST',
      body: JSON.stringify(options),
    })
  }

  static async getUserPermissions(id: string) {
    return apiFetcher({
      url: `/files/${id}/get_user_permissions`,
      method: 'GET',
    }) as Promise<UserPermission[]>
  }

  static useGetUserPermissions(
    id: string | null | undefined,
    swrOptions?: any,
  ) {
    const url = `/files/${id}/get_user_permissions`
    return useSWR<UserPermission[]>(
      id ? url : null,
      () => apiFetcher({ url, method: 'GET' }) as Promise<UserPermission[]>,
      swrOptions,
    )
  }

  static async getGroupPermissions(id: string) {
    return apiFetcher<GroupPermission[]>({
      url: `/files/${id}/get_group_permissions`,
      method: 'GET',
    })
  }

  static useGetGroupPermissions(
    id: string | null | undefined,
    swrOptions?: any,
  ) {
    const url = `/files/${id}/get_group_permissions`
    return useSWR<GroupPermission[]>(
      id ? url : null,
      () => apiFetcher({ url, method: 'GET' }) as Promise<GroupPermission[]>,
      swrOptions,
    )
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
    if (options?.query) {
      params.query = encodeQuery(JSON.stringify(options.query))
    }
    return new URLSearchParams(params)
  }
}
