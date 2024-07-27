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
import { paramsFromListOptions } from '@/client/api/query-helpers'
import { ListOptions } from '@/client/api/types/queries'
import { apiFetcher } from '@/client/fetcher'
import { Organization } from './organization'
import { PermissionType } from './permission'

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

export interface PatchNameOptions {
  name: string
}

export interface PatchStorageCapacityOptions {
  storageCapacity: number
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
    const url = `/workspaces?${paramsFromListOptions(options)}`
    return useSWR<List>(
      url,
      () => apiFetcher({ url, method: 'GET' }) as Promise<List>,
      swrOptions,
    )
  }

  static async list(options?: ListOptions) {
    return apiFetcher({
      url: `/workspaces?${paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<List>
  }

  static async create(options: CreateOptions) {
    return apiFetcher({
      url: '/workspaces',
      method: 'POST',
      body: JSON.stringify(options),
    }) as Promise<Workspace>
  }

  static async patchName(id: string, options: PatchNameOptions) {
    return apiFetcher({
      url: `/workspaces/${id}/name`,
      method: 'PATCH',
      body: JSON.stringify(options),
    }) as Promise<Workspace>
  }

  static async patchStorageCapacity(
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
