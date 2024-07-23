// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import useSWR, { SWRConfiguration } from 'swr'
import { apiFetcher } from '@/client/fetcher'
import { User } from '@/client/idp/user'
import { getConfig } from '@/config/config'
import { getAccessTokenOrRedirect } from '@/infra/token'
import { encodeQuery } from '@/lib/helpers/query'
import { Group } from './group'
import { PermissionType } from './permission'
import { Snapshot } from './snapshot'

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
  type?: FileType
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

export type CopyManyOptions = {
  sourceIds: string[]
  targetId: string
}

export type CopyManyResult = {
  new: string[]
  succeeded: string[]
  failed: string[]
}

export type DeleteOptions = {
  ids: string[]
}

export type PatchNameOptions = {
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

export type CreateOptions = {
  type: FileType
  workspaceId: string
  parentId?: string
  name?: string
  blob?: Blob
  request?: XMLHttpRequest
  onProgress?: (value: number) => void
}

export type PatchOptions = {
  id: string
  request: XMLHttpRequest
  blob: Blob
  onProgress?: (value: number) => void
}

type ListQueryParams = {
  page?: string
  size?: string
  sort_by?: string
  sort_order?: string
  type?: string
  query?: string
}

export default class FileAPI {
  static async create({
    type,
    workspaceId,
    parentId,
    name,
    request,
    blob,
    onProgress,
  }: CreateOptions): Promise<File> {
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
  }: PatchOptions): Promise<File> {
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

  static async list(id: string, options: ListOptions) {
    return apiFetcher({
      url: `/files/${id}/list?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<List>
  }

  static useList(
    id: string | undefined,
    options: ListOptions,
    swrOptions?: SWRConfiguration,
  ) {
    const url = `/files/${id}/list?${this.paramsFromListOptions(options)}`
    return useSWR<List | undefined>(
      id ? url : null,
      () => apiFetcher({ url, method: 'GET' }),
      swrOptions,
    )
  }

  static paramsFromListOptions(options?: ListOptions): URLSearchParams {
    const params: ListQueryParams = {}
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

  static async patchName(id: string, options: PatchNameOptions) {
    return apiFetcher({
      url: `/files/${id}/name`,
      method: 'PATCH',
      body: JSON.stringify(options),
    }) as Promise<File>
  }

  static async delete(options: DeleteOptions) {
    return apiFetcher({
      url: `/files`,
      method: 'DELETE',
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

  static async copyOne(id: string, targetId: string) {
    return apiFetcher({
      url: `/files/${id}/copy/${targetId}`,
      method: 'POST',
    }) as Promise<File>
  }

  static async copyMany(options: CopyManyOptions) {
    return apiFetcher({
      url: `/files/copy`,
      method: 'POST',
      body: JSON.stringify(options),
    }) as Promise<CopyManyResult>
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

  static async get(id: string) {
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

  static useGetUserPermissions(
    id: string | null | undefined,
    swrOptions?: SWRConfiguration,
  ) {
    const url = `/files/${id}/user_permissions`
    return useSWR<UserPermission[]>(
      id ? url : null,
      () => apiFetcher({ url, method: 'GET' }) as Promise<UserPermission[]>,
      swrOptions,
    )
  }

  static useGetGroupPermissions(
    id: string | null | undefined,
    swrOptions?: SWRConfiguration,
  ) {
    const url = `/files/${id}/group_permissions`
    return useSWR<GroupPermission[]>(
      id ? url : null,
      () => apiFetcher({ url, method: 'GET' }) as Promise<GroupPermission[]>,
      swrOptions,
    )
  }
}
