import useSWR, { SWRConfiguration } from 'swr'
import { apiFetcher } from '@/client/fetcher'
import { Snapshot } from './snapshot'

export type CreateOptions = {
  languageId: string
}

export type Language = {
  id: string
  iso6393: string
  name: string
}

export type Info = {
  isAvailable: boolean
  isOutdated: boolean
  snapshot?: Snapshot
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

export type ListEntitiesOptions = {
  query?: string
  size?: number
  page?: number
  sortBy?: SortBy
  sortOrder?: SortOrder
}

export enum SortBy {
  Name = 'name',
  Frequency = 'frequency',
}

export enum SortOrder {
  Asc = 'asc',
  Desc = 'desc',
}

type ListEntitiesQueryParams = {
  page?: string
  size?: string
  sort_by?: string
  sort_order?: string
  query?: string
}

export default class InsightsAPI {
  static create(id: string, options: CreateOptions, showError = true) {
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
    return useSWR<Info>(
      id ? url : null,
      () =>
        apiFetcher({
          url,
          method: 'GET',
          showError: false,
        }) as Promise<Info>,
      swrOptions,
    )
  }

  static async getInfo(id: string) {
    return apiFetcher({
      url: `/insights/${id}/info`,
      method: 'GET',
    }) as Promise<Info>
  }

  static useGetLanguages(swrOptions?: SWRConfiguration) {
    const url = `/insights/languages`
    return useSWR<Language[]>(
      url,
      () => apiFetcher({ url, method: 'GET' }) as Promise<Language[]>,
      swrOptions,
    )
  }

  static useListEntities(
    id: string | null | undefined,
    options?: ListEntitiesOptions,
    swrOptions?: SWRConfiguration,
  ) {
    const url = `/insights/${id}/entities?${this.paramsFromListOptions(options)}`
    return useSWR<EntityList>(
      id ? url : null,
      () => apiFetcher({ url, method: 'GET' }) as Promise<EntityList>,
      swrOptions,
    )
  }

  static paramsFromListOptions(options?: ListEntitiesOptions): URLSearchParams {
    const params: ListEntitiesQueryParams = {}
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
