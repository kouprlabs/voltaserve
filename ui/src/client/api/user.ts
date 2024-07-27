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

export type User = {
  id: string
  username: string
  email: string
  fullName: string
  picture?: string
}

export type List = {
  data: User[]
  totalPages: number
  totalElements: number
  page: number
  size: number
}

export default class UserAPI {
  static async list(options?: ListOptions) {
    return apiFetcher({
      url: `/users?${paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<List>
  }

  static useList(options?: ListOptions, swrOptions?: SWRConfiguration) {
    const url = `/users?${paramsFromListOptions(options)}`
    return useSWR<List>(
      url,
      () => apiFetcher({ url, method: 'GET' }) as Promise<List>,
      swrOptions,
    )
  }
}
