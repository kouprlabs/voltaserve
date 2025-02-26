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

export type EntityCreateOptions = {
  language: string
}

export type Entity = {
  text: string
  label: string
  frequency: number
}

export type EntityList = {
  data: Entity[]
  totalPages: number
  totalElements: number
  page: number
  size: number
}

export type EntityListOptions = {
  query?: string
  size?: number
  page?: number
  sortBy?: EntitySortBy
  sortOrder?: EntitySortOrder
}

export enum EntitySortBy {
  Name = 'name',
  Frequency = 'frequency',
}

export enum EntitySortOrder {
  Asc = 'asc',
  Desc = 'desc',
}

type EntityListQueryParams = {
  page?: string
  size?: string
  sort_by?: string
  sort_order?: string
  query?: string
}

export class EntityAPI {
  static create(id: string, options: EntityCreateOptions, showError = true) {
    return apiFetcher({
      url: `/entities/${id}`,
      method: 'POST',
      body: JSON.stringify(options),
      showError,
    })
  }

  static delete(id: string) {
    return apiFetcher({
      url: `/entities/${id}`,
      method: 'DELETE',
    })
  }

  static useList(
    id: string | null | undefined,
    options?: EntityListOptions,
    swrOptions?: SWRConfiguration,
  ) {
    const url = `/entities/${id}?${this.paramsFromListOptions(options)}`
    return useSWR<EntityList>(
      id ? url : null,
      () => apiFetcher({ url, method: 'GET' }) as Promise<EntityList>,
      swrOptions,
    )
  }

  static paramsFromListOptions(options?: EntityListOptions): URLSearchParams {
    const params: EntityListQueryParams = {}
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
}
