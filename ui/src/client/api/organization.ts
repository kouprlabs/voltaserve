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

export type Organization = {
  id: string
  name: string
  permission: PermissionType
  createTime: string
  updateTime?: string
}

export type List = {
  data: Organization[]
  totalPages: number
  totalElements: number
  page: number
  size: number
}

export type ListOptions = {
  query?: string
  size?: number
  page?: number
  sortBy?: SortBy
  sortOrder?: SortOrder
}

export type CreateOptions = {
  name: string
  image?: string
}

export type PatchNameOptions = {
  name: string
}

export type RemoveMemberOptions = {
  userId: string
}

type ListQueryParams = {
  page?: string
  size?: string
  sort_by?: string
  sort_order?: string
  query?: string
}

export default class OrganizationAPI {
  static useGet(id: string | null | undefined, swrOptions?: SWRConfiguration) {
    const url = `/organizations/${id}`
    return useSWR<Organization>(
      id ? url : null,
      () => apiFetcher({ url, method: 'GET' }) as Promise<Organization>,
      swrOptions,
    )
  }

  static list(options?: ListOptions) {
    return apiFetcher({
      url: `/organizations?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<List>
  }

  static useList(options?: ListOptions, swrOptions?: SWRConfiguration) {
    const url = `/organizations?${this.paramsFromListOptions(options)}`
    return useSWR<List>(
      url,
      () => apiFetcher({ url, method: 'GET' }) as Promise<List>,
      swrOptions,
    )
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
      url: `/organizations`,
      method: 'POST',
      body: JSON.stringify(options),
    }) as Promise<Organization>
  }

  static patchName(id: string, options: PatchNameOptions) {
    return apiFetcher({
      url: `/organizations/${id}/name`,
      method: 'PATCH',
      body: JSON.stringify(options),
    }) as Promise<Organization>
  }

  static async delete(id: string) {
    return apiFetcher({
      url: `/organizations/${id}`,
      method: 'DELETE',
    })
  }

  static async leave(id: string) {
    return apiFetcher({
      url: `/organizations/${id}/leave`,
      method: 'POST',
    })
  }

  static async removeMember(id: string, options: RemoveMemberOptions) {
    return apiFetcher({
      url: `/organizations/${id}/members`,
      method: 'DELETE',
      body: JSON.stringify(options),
    })
  }
}
