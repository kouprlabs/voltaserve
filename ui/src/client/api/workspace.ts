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
import { Organization } from './organization'
import { PermissionType } from './permission'

export enum SortBy {
  Name = 'name',
  DateCreated = 'date_created',
  DateModified = 'date_modified',
}

export enum SortOrder {
  Asc = 'asc',
  Desc = 'desc',
}

export type Workspace = {
  id: string
  name: string
  permission: PermissionType
  storageCapacity: number
  rootId: string
  organization: Organization
  createTime: string
  updateTime?: string
}

export type List = {
  data: Workspace[]
  totalPages: number
  totalElements: number
  page: number
  size: number
}

export interface CreateOptions {
  name: string
  image?: string
  organizationId: string
  storageCapacity: number
}

export type ListOptions = {
  query?: string
  size?: number
  page?: number
  sortBy?: SortBy
  sortOrder?: SortOrder
}

export interface PatchNameOptions {
  name: string
}

export interface PatchStorageCapacityOptions {
  storageCapacity: number
}

type ListQueryParams = {
  page?: string
  size?: string
  sort_by?: string
  sort_order?: string
  query?: string
}

export default class WorkspaceAPI {
  static useGet(id: string | null | undefined, swrOptions?: SWRConfiguration) {
    const url = `/workspaces/${id}`
    return useSWR<Workspace>(
      id ? url : null,
      () => apiFetcher({ url, method: 'GET' }) as Promise<Workspace>,
      swrOptions,
    )
  }

  static useList(options?: ListOptions, swrOptions?: SWRConfiguration) {
    const url = `/workspaces?${this.paramsFromListOptions(options)}`
    return useSWR<List>(
      url,
      () => apiFetcher({ url, method: 'GET' }) as Promise<List>,
      swrOptions,
    )
  }

  static list(options?: ListOptions) {
    return apiFetcher({
      url: `/workspaces?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<List>
  }

  static paramsFromListOptions(options?: ListOptions): URLSearchParams {
    const params: ListQueryParams = {}
    if (options?.query) {
      params.query = encodeURIComponent(options.query.toString())
    }
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
    return new URLSearchParams(params)
  }

  static create(options: CreateOptions) {
    return apiFetcher({
      url: '/workspaces',
      method: 'POST',
      body: JSON.stringify(options),
    }) as Promise<Workspace>
  }

  static patchName(id: string, options: PatchNameOptions) {
    return apiFetcher({
      url: `/workspaces/${id}/name`,
      method: 'PATCH',
      body: JSON.stringify(options),
    }) as Promise<Workspace>
  }

  static patchStorageCapacity(
    id: string,
    options: PatchStorageCapacityOptions,
  ) {
    return apiFetcher({
      url: `/workspaces/${id}/storage_capacity`,
      method: 'PATCH',
      body: JSON.stringify(options),
    }) as Promise<Workspace>
  }

  static async delete(id: string) {
    return apiFetcher({
      url: `/workspaces/${id}`,
      method: 'DELETE',
    })
  }
}
