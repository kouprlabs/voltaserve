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
import { PermissionType } from './permission'

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

export default class OrganizationAPI {
  static useGet(id: string | null | undefined, swrOptions?: SWRConfiguration) {
    const url = `/organizations/${id}`
    return useSWR<Organization>(
      id ? url : null,
      () => apiFetcher({ url, method: 'GET' }) as Promise<Organization>,
      swrOptions,
    )
  }

  static async list(options?: ListOptions) {
    return apiFetcher({
      url: `/organizations?${paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<List>
  }

  static useList(options?: ListOptions, swrOptions?: SWRConfiguration) {
    const url = `/organizations?${paramsFromListOptions(options)}`
    return useSWR<List>(
      url,
      () => apiFetcher({ url, method: 'GET' }) as Promise<List>,
      swrOptions,
    )
  }

  static async create(options: CreateOptions) {
    return apiFetcher({
      url: `/organizations`,
      method: 'POST',
      body: JSON.stringify(options),
    }) as Promise<Organization>
  }

  static async patchName(id: string, options: PatchNameOptions) {
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
