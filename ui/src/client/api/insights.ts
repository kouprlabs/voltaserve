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
import { Snapshot } from './snapshot'

export type InsightsCreateOptions = {
  languageId: string
}

export type InsightsLanguage = {
  id: string
  iso6393: string
  name: string
}

export type InsightsInfo = {
  isAvailable: boolean
  isOutdated: boolean
  snapshot?: Snapshot
}

export type InsightsEntity = {
  text: string
  label: string
  frequency: number
}

export type InsightsEntityList = {
  data: InsightsEntity[]
  totalPages: number
  totalElements: number
  page: number
  size: number
}

export type InsightsListEntitiesOptions = {
  query?: string
  size?: number
  page?: number
  sortBy?: InsightsSortBy
  sortOrder?: InsightsSortOrder
}

export enum InsightsSortBy {
  Name = 'name',
  Frequency = 'frequency',
}

export enum InsightsSortOrder {
  Asc = 'asc',
  Desc = 'desc',
}

type InsightsListEntitiesQueryParams = {
  page?: string
  size?: string
  sort_by?: string
  sort_order?: string
  query?: string
}

export class InsightsAPI {
  static create(id: string, options: InsightsCreateOptions, showError = true) {
    return apiFetcher({
      url: `/insights/${id}`,
      method: 'POST',
      body: JSON.stringify(options),
      showError,
    })
  }

  static patch(id: string) {
    return apiFetcher({
      url: `/insights/${id}`,
      method: 'PATCH',
    })
  }

  static delete(id: string) {
    return apiFetcher({
      url: `/insights/${id}`,
      method: 'DELETE',
    })
  }

  static useGetInfo(
    id: string | null | undefined,
    swrOptions?: SWRConfiguration,
  ) {
    const url = `/insights/${id}/info`
    return useSWR<InsightsInfo>(
      id ? url : null,
      () =>
        apiFetcher({
          url,
          method: 'GET',
          showError: false,
        }) as Promise<InsightsInfo>,
      swrOptions,
    )
  }

  static getInfo(id: string) {
    return apiFetcher({
      url: `/insights/${id}/info`,
      method: 'GET',
    }) as Promise<InsightsInfo>
  }

  static useGetLanguages(swrOptions?: SWRConfiguration) {
    const url = `/insights/languages`
    return useSWR<InsightsLanguage[]>(
      url,
      () => apiFetcher({ url, method: 'GET' }) as Promise<InsightsLanguage[]>,
      swrOptions,
    )
  }

  static useListEntities(
    id: string | null | undefined,
    options?: InsightsListEntitiesOptions,
    swrOptions?: SWRConfiguration,
  ) {
    const url = `/insights/${id}/entities?${this.paramsFromListOptions(options)}`
    return useSWR<InsightsEntityList>(
      id ? url : null,
      () => apiFetcher({ url, method: 'GET' }) as Promise<InsightsEntityList>,
      swrOptions,
    )
  }

  static paramsFromListOptions(
    options?: InsightsListEntitiesOptions,
  ): URLSearchParams {
    const params: InsightsListEntitiesQueryParams = {}
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
