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
import { PermissionType } from '@/client/api/permission'
import { apiFetcher } from '@/client/fetcher'
import { Organization } from './organization'

export enum GroupSortBy {
  Name = 'name',
  DateCreated = 'date_created',
  DateModified = 'date_modified',
}

export enum GroupSortOrder {
  Asc = 'asc',
  Desc = 'desc',
}

export type Group = {
  id: string
  name: string
  organization: Organization
  permission: PermissionType
  createTime: string
  updateTime?: string
}

export type GroupList = {
  data: Group[]
  totalPages: number
  totalElements: number
  page: number
  size: number
}

export type GroupListOptions = {
  query?: string
  organizationId?: string
  size?: number
  page?: number
  sortBy?: GroupSortBy
  sortOrder?: GroupSortOrder
}

export type GroupCreateOptions = {
  name: string
  image?: string
  organizationId: string
}

export type GroupPatchNameOptions = {
  name: string
}

export type GroupAddMemberOptions = {
  userId: string
}

export type GroupRemoveMemberOptions = {
  userId: string
}

type GroupListQueryParams = {
  page?: string
  size?: string
  sort_by?: string
  sort_order?: string
  query?: string
  organization_id?: string
}

export class GroupAPI {
  static create(options: GroupCreateOptions) {
    return apiFetcher({
      url: `/groups`,
      method: 'POST',
      body: JSON.stringify(options),
    }) as Promise<Group>
  }

  static patchName(id: string, options: GroupPatchNameOptions) {
    return apiFetcher({
      url: `/groups/${id}/name`,
      method: 'PATCH',
      body: JSON.stringify(options),
    }) as Promise<Group>
  }

  static useGet(id: string | null | undefined, swrOptions?: SWRConfiguration) {
    const url = `/groups/${id}`
    return useSWR<Group>(
      id ? url : null,
      () => apiFetcher({ url, method: 'GET' }) as Promise<Group>,
      swrOptions,
    )
  }

  static useList(options?: GroupListOptions, swrOptions?: SWRConfiguration) {
    const url = `/groups?${this.paramsFromListOptions(options)}`
    return useSWR<GroupList>(
      url,
      () => apiFetcher({ url, method: 'GET' }) as Promise<GroupList>,
      swrOptions,
    )
  }

  static paramsFromListOptions(options?: GroupListOptions): URLSearchParams {
    const params: GroupListQueryParams = {}
    if (options?.query) {
      params.query = encodeURIComponent(options.query.toString())
    }
    if (options?.organizationId) {
      params.organization_id = options.organizationId.toString()
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

  static delete(id: string) {
    return apiFetcher({
      url: `/groups/${id}`,
      method: 'DELETE',
    })
  }

  static addMember(id: string, options: GroupAddMemberOptions) {
    return apiFetcher({
      url: `/groups/${id}/members`,
      method: 'POST',
      body: JSON.stringify(options),
    })
  }

  static removeMember(id: string, options: GroupRemoveMemberOptions) {
    return apiFetcher({
      url: `/groups/${id}/members`,
      method: 'DELETE',
      body: JSON.stringify(options),
    })
  }
}
