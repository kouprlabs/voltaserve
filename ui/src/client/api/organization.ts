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
import { PermissionType } from './permission'

export enum OrganizationSortBy {
  Name = 'name',
  DateCreated = 'date_created',
  DateModified = 'date_modified',
}

export enum OrganizationSortOrder {
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

export type OrganizationList = {
  data: Organization[]
  totalPages: number
  totalElements: number
  page: number
  size: number
}

export type OrganizationListOptions = {
  query?: string
  size?: number
  page?: number
  sortBy?: OrganizationSortBy
  sortOrder?: OrganizationSortOrder
}

export type OrganizationCreateOptions = {
  name: string
  image?: string
}

export type OrganizationPatchNameOptions = {
  name: string
}

export type OrganizationRemoveMemberOptions = {
  userId: string
}

type OrganizationListQueryParams = {
  page?: string
  size?: string
  sort_by?: string
  sort_order?: string
  query?: string
}

export class OrganizationAPI {
  static useGet(id: string | null | undefined, swrOptions?: SWRConfiguration) {
    const url = `/organizations/${id}`
    return useSWR<Organization>(
      id ? url : null,
      () => apiFetcher({ url, method: 'GET' }) as Promise<Organization>,
      swrOptions,
    )
  }

  static list(options?: OrganizationListOptions) {
    return apiFetcher({
      url: `/organizations?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<OrganizationList>
  }

  static useList(
    options?: OrganizationListOptions,
    swrOptions?: SWRConfiguration,
  ) {
    const url = `/organizations?${this.paramsFromListOptions(options)}`
    return useSWR<OrganizationList>(
      url,
      () => apiFetcher({ url, method: 'GET' }) as Promise<OrganizationList>,
      swrOptions,
    )
  }

  static paramsFromListOptions(
    options?: OrganizationListOptions,
  ): URLSearchParams {
    const params: OrganizationListQueryParams = {}
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

  static create(options: OrganizationCreateOptions) {
    return apiFetcher({
      url: `/organizations`,
      method: 'POST',
      body: JSON.stringify(options),
    }) as Promise<Organization>
  }

  static patchName(id: string, options: OrganizationPatchNameOptions) {
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

  static async removeMember(
    id: string,
    options: OrganizationRemoveMemberOptions,
  ) {
    return apiFetcher({
      url: `/organizations/${id}/members`,
      method: 'DELETE',
      body: JSON.stringify(options),
    })
  }
}
