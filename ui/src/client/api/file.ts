// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import useSWR, { SWRConfiguration } from 'swr'
import { apiFetcher } from '@/client/fetcher'
import { AuthUser } from '@/client/idp/user'
import { getAccessTokenOrRedirect } from '@/client/token'
import { getConfig } from '@/config/config'
import { encodeQuery } from '@/lib/helpers/query'
import { Group } from './group'
import { PermissionType } from './permission'
import { Snapshot } from './snapshot'

export enum FileType {
  File = 'file',
  Folder = 'folder',
}

export enum FileSortBy {
  Name = 'name',
  Kind = 'kind',
  Size = 'size',
  DateCreated = 'date_created',
  DateModified = 'date_modified',
}

export enum FileSortOrder {
  Asc = 'asc',
  Desc = 'desc',
}

export type File = {
  id: string
  workspaceId: string
  name: string
  type: FileType
  parentId: string
  permission: PermissionType
  isShared?: boolean
  snapshot?: Snapshot
  createTime: string
  updateTime?: string
}

export type FileList = {
  data: File[]
  totalPages: number
  totalElements: number
  page: number
  size: number
  query?: FileQuery
}

export type FileUserPermission = {
  id: string
  user: AuthUser
  permission: PermissionType
}

export type FileGroupPermission = {
  id: string
  group: Group
  permission: PermissionType
}

export type FileQuery = {
  text?: string
  type?: FileType
  createTimeAfter?: number
  createTimeBefore?: number
  updateTimeAfter?: number
  updateTimeBefore?: number
}

export type FileListOptions = {
  size?: number
  page?: number
  sortBy?: FileSortBy
  sortOrder?: FileSortOrder
  query?: FileQuery
}

export type FileMoveManyOptions = {
  sourceIds: string[]
  targetId: string
}

export type FileMoveManyResult = {
  succeeded: string[]
  failed: string[]
}

export type FileCopyManyOptions = {
  sourceIds: string[]
  targetId: string
}

export type FileCopyManyResult = {
  new: string[]
  succeeded: string[]
  failed: string[]
}

export type FileDeleteManyOptions = {
  ids: string
}

export type FileDeleteManyResult = {
  succeeded: string[]
  failed: string[]
}

export type FileFilePatchNameOptions = {
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

export type FileCreateOptions = {
  type: FileType
  workspaceId: string
  parentId?: string
  name?: string
  blob?: Blob
  request?: XMLHttpRequest
  onProgress?: (value: number) => void
}

export type FilePatchOptions = {
  id: string
  request: XMLHttpRequest
  blob: Blob
  onProgress?: (value: number) => void
}

type FileListQueryParams = {
  page?: string
  size?: string
  sort_by?: string
  sort_order?: string
  type?: string
  query?: string
}

export class FileAPI {
  static create({
    type,
    workspaceId,
    parentId,
    name,
    request,
    blob,
    onProgress,
  }: FileCreateOptions): Promise<File> {
    const params = new URLSearchParams({ type, workspace_id: workspaceId })
    if (parentId) {
      params.append('parent_id', parentId)
    }
    if (name) {
      params.append('name', name)
    }
    if (type === FileType.File && request && blob) {
      return this.upload(
        `${getConfig().apiURL}/files?${params}`,
        'POST',
        request,
        blob,
        onProgress,
      )
    } else if (type === FileType.Folder) {
      return apiFetcher({
        url: `/files?${params}`,
        method: 'POST',
      }) as Promise<File>
    }
    throw new Error('Invalid parameters')
  }

  static async patch({
    id,
    request,
    blob,
    onProgress,
  }: FilePatchOptions): Promise<File> {
    return this.upload(
      `${getConfig().apiURL}/files/${id}`,
      'PATCH',
      request,
      blob,
      onProgress,
    )
  }

  private static async upload(
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
          } catch (error) {
            reject(error)
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

  static list(id: string, options: FileListOptions) {
    return apiFetcher({
      url: `/files/${id}/list?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<FileList>
  }

  static useList(
    id: string | undefined,
    options: FileListOptions,
    swrOptions?: SWRConfiguration,
  ) {
    const url = `/files/${id}/list?${this.paramsFromListOptions(options)}`
    return useSWR<FileList | undefined>(
      id ? url : null,
      () => apiFetcher({ url, method: 'GET' }),
      swrOptions,
    )
  }

  static paramsFromListOptions(options?: FileListOptions): URLSearchParams {
    const params: FileListQueryParams = {}
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
    if (options?.query) {
      params.query = encodeQuery(JSON.stringify(options.query))
    }
    return new URLSearchParams(params)
  }

  static useGetPath(
    id: string | null | undefined,
    swrOptions?: SWRConfiguration,
  ) {
    const url = `/files/${id}/path`
    return useSWR<File[]>(
      id ? url : null,
      () => apiFetcher({ url, method: 'GET' }) as Promise<File[]>,
      swrOptions,
    )
  }

  static patchName(id: string, options: FileFilePatchNameOptions) {
    return apiFetcher({
      url: `/files/${id}/name`,
      method: 'PATCH',
      body: JSON.stringify(options),
    }) as Promise<File>
  }

  static async deleteOne(id: string) {
    return apiFetcher({
      url: `/files/${id}`,
      method: 'DELETE',
    })
  }

  static deleteMany(options: FileDeleteManyOptions) {
    return apiFetcher({
      url: `/files`,
      method: 'DELETE',
      body: JSON.stringify(options),
    }) as Promise<FileDeleteManyResult>
  }

  static async moveOne(id: string, targetId: string) {
    return apiFetcher({
      url: `/files/${id}/move/${targetId}`,
      method: 'POST',
    })
  }

  static moveMany(options: FileMoveManyOptions) {
    return apiFetcher({
      url: `/files/move`,
      method: 'POST',
      body: JSON.stringify(options),
    }) as Promise<FileMoveManyResult>
  }

  static copyOne(id: string, targetId: string) {
    return apiFetcher({
      url: `/files/${id}/copy/${targetId}`,
      method: 'POST',
    }) as Promise<File>
  }

  static copyMany(options: FileCopyManyOptions) {
    return apiFetcher({
      url: `/files/copy`,
      method: 'POST',
      body: JSON.stringify(options),
    }) as Promise<FileCopyManyResult>
  }

  static useGet(
    id: string | null | undefined,
    swrOptions?: SWRConfiguration,
    showError = true,
  ) {
    const url = `/files/${id}`
    return useSWR(
      id ? url : null,
      () => apiFetcher({ url, method: 'GET', showError }) as Promise<File>,
      swrOptions,
    )
  }

  static get(id: string) {
    return apiFetcher({
      url: `/files/${id}`,
      method: 'GET',
    }) as Promise<File>
  }

  static useGetCount(
    id: string | null | undefined,
    swrOptions?: SWRConfiguration,
  ) {
    const url = `/files/${id}/count`
    return useSWR<number>(
      id ? url : null,
      () => apiFetcher({ url, method: 'GET' }) as Promise<number>,
      swrOptions,
    )
  }

  static async grantUserPermission(options: FileGrantUserPermissionOptions) {
    return apiFetcher({
      url: `/files/grant_user_permission`,
      method: 'POST',
      body: JSON.stringify(options),
    })
  }

  static async revokeUserPermission(options: FileRevokeUserPermissionOptions) {
    return apiFetcher({
      url: `/files/revoke_user_permission`,
      method: 'POST',
      body: JSON.stringify(options),
    })
  }

  static async grantGroupPermission(options: FileGrantGroupPermissionOptions) {
    return apiFetcher({
      url: `/files/grant_group_permission`,
      method: 'POST',
      body: JSON.stringify(options),
    })
  }

  static async revokeGroupPermission(
    options: FileRevokeGroupPermissionOptions,
  ) {
    return apiFetcher({
      url: `/files/revoke_group_permission`,
      method: 'POST',
      body: JSON.stringify(options),
    })
  }

  static useGetUserPermissions(
    id: string | null | undefined,
    swrOptions?: SWRConfiguration,
  ) {
    const url = `/files/${id}/user_permissions`
    return useSWR<FileUserPermission[]>(
      id ? url : null,
      () => apiFetcher({ url, method: 'GET' }) as Promise<FileUserPermission[]>,
      swrOptions,
    )
  }

  static useGetGroupPermissions(
    id: string | null | undefined,
    swrOptions?: SWRConfiguration,
  ) {
    const url = `/files/${id}/group_permissions`
    return useSWR<FileGroupPermission[]>(
      id ? url : null,
      () =>
        apiFetcher({ url, method: 'GET' }) as Promise<FileGroupPermission[]>,
      swrOptions,
    )
  }
}
