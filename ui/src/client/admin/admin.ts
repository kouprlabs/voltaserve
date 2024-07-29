// Copyright 2024 Mateusz Kaźmierczak.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import useSWR, { SWRConfiguration } from 'swr'
import { adminFetcher, baseFetcher } from '@/client/fetcher'
import { getConfig } from '@/config/config'
import { getAdminAccessToken } from '@/infra/admin-token'

export type UserManagement = {
  id: string
  fullName: string
  username: string
  email: string
  picture?: string
  isEmailConfirmed: boolean
  createTime: Date
  updateTime: Date
}

export type UserManagementList = {
  data: UserManagement[]
  totalElements: number
  page: number
  size: number
}

export type IndexManagement = {
  tablename: string
  indexname: string
  indexdef: string
}

export type IndexManagementList = {
  data: IndexManagement[]
  totalElements: number
  page: number
  size: number
}

export type ListOptions = {
  size?: number
  page?: number
}

type ListQueryParams = {
  page?: string
  size?: string
}

export default class AdminApi {
  static async adminAuthenticate() {
    const adminToken = getAdminAccessToken()
    if (!adminToken) {
      return false
    }
    const response = await baseFetcher(
      `${getConfig().adminURL}/`,
      {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${adminToken}`,
          'Access-Control-Allow-Origin': `${getConfig().adminURL}`, // TODO: To be deleted after local tests
        },
      },
      false,
    )
    if (response) {
      return response.ok
    } else {
      return false
    }
  }

  static useListUsers(options?: ListOptions, swrOptions?: SWRConfiguration) {
    const url = `/user/all?${this.paramsFromListOptions(options)}`
    return useSWR<UserManagementList>(
      url,
      () =>
        adminFetcher({
          url,
          method: 'GET',
          adminAuthenticate: true,
        }) as Promise<UserManagementList>,
      swrOptions,
    )
  }

  static useListIndexes(options: ListOptions, swrOptions?: SWRConfiguration) {
    const url = `/index/all?${this.paramsFromListOptions(options)}`
    return useSWR<IndexManagementList>(
      url,
      () =>
        adminFetcher({
          url,
          method: 'GET',
          adminAuthenticate: true,
        }) as Promise<IndexManagementList>,
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
    return new URLSearchParams(params)
  }
}
