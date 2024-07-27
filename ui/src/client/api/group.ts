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

export type Group = {
  id: string
  name: string
  organization: Organization
  permission: string
  createTime: string
  updateTime?: string
}

export type List = {
  data: Group[]
  totalPages: number
  totalElements: number
  page: number
  size: number
}

export type CreateOptions = {
  name: string
  image?: string
  organizationId: string
}

export type PatchNameOptions = {
  name: string
}

export type AddMemberOptions = {
  userId: string
}

export type RemoveMemberOptions = {
  userId: string
}

export default class GroupAPI {
  static create(options: CreateOptions) {
    return apiFetcher({
      url: `/groups`,
      method: 'POST',
      body: JSON.stringify(options),
    }) as Promise<Group>
  }

  static patchName(id: string, options: PatchNameOptions) {
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

  static useList(options?: ListOptions, swrOptions?: SWRConfiguration) {
    const url = `/groups?${paramsFromListOptions(options)}`
    return useSWR<List>(
      url,
      () => apiFetcher({ url, method: 'GET' }) as Promise<List>,
      swrOptions,
    )
  }

  static delete(id: string) {
    return apiFetcher({
      url: `/groups/${id}`,
      method: 'DELETE',
    })
  }

  static addMember(id: string, options: AddMemberOptions) {
    return apiFetcher({
      url: `/groups/${id}/members`,
      method: 'POST',
      body: JSON.stringify(options),
    })
  }

  static removeMember(id: string, options: RemoveMemberOptions) {
    return apiFetcher({
      url: `/groups/${id}/members`,
      method: 'DELETE',
      body: JSON.stringify(options),
    })
  }
}
