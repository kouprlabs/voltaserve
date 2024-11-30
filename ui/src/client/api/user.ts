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
import { Picture } from '@/client/types'

export enum SortBy {
  Email = 'email',
  FullName = 'full_name',
}

export enum SortOrder {
  Asc = 'asc',
  Desc = 'desc',
}

export type User = {
  id: string
  username: string
  email: string
  fullName: string
  picture?: Picture
}

export type List = {
  data: User[]
  totalPages: number
  totalElements: number
  page: number
  size: number
}

export type ListOptions = {
  query?: string
  organizationId?: string
  groupId?: string
  excludeGroupMembers?: boolean
  excludeMe?: boolean
  size?: number
  page?: number
  sortBy?: SortBy
  sortOrder?: SortOrder
}

type ListQueryParams = {
  page?: string
  size?: string
  sort_by?: string
  sort_order?: string
  query?: string
  organization_id?: string
  group_id?: string
  exclude_group_members?: string
  exclude_me?: string
}

export default class UserAPI {
  static list(options?: ListOptions) {
    return apiFetcher({
      url: `/users?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<List>
  }

  static useList(options?: ListOptions, swrOptions?: SWRConfiguration) {
    const url = `/users?${this.paramsFromListOptions(options)}`
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
    if (options?.organizationId) {
      params.organization_id = options.organizationId.toString()
    }
    if (options?.groupId) {
      params.group_id = options.groupId.toString()
    }
    if (options?.excludeGroupMembers) {
      params.exclude_group_members = options.excludeGroupMembers.toString()
    }
    if (options?.excludeMe) {
      params.exclude_me = options.excludeMe.toString()
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
}
