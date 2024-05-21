import useSWR, { SWRConfiguration } from 'swr'
import { apiFetcher } from '@/client/fetcher'

export type Language = {
  id: string
  iso6393: string
  name: string
}

export type Summary = {
  hasLanguage: boolean
  hasOcr: boolean
  hasText: boolean
  hasEntities: boolean
}

export type Entity = {
  text: string
  label: string
}

export type EntityList = {
  data: Entity[]
  totalPages: number
  totalElements: number
  page: number
  size: number
}

export type ListEntitiesOptions = {
  query: string
  page: number
  size: number
  sortBy: SortBy
  sortOrder: SortOrder
}

export enum SortBy {
  Name = 'name',
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

export default class AnalysisAPI {
  static async getSummary(id: string) {
    return apiFetcher({
      url: `/analysis/${id}/summary`,
      method: 'GET',
    }) as Promise<Summary>
  }

  static useGetSummary(
    id: string | null | undefined,
    swrOptions?: SWRConfiguration,
  ) {
    const url = `/analysis/${id}/summary`
    return useSWR<Summary>(
      id ? url : null,
      () => apiFetcher({ url, method: 'GET' }) as Promise<Summary>,
      swrOptions,
    )
  }

  static async getLanguages() {
    return apiFetcher({
      url: `/analysis/languages`,
      method: 'GET',
    }) as Promise<Language[]>
  }

  static useGetLanguages(swrOptions?: SWRConfiguration) {
    const url = `/analysis/languages`
    return useSWR<Language[]>(
      url,
      () => apiFetcher({ url, method: 'GET' }) as Promise<Language[]>,
      swrOptions,
    )
  }

  static async listEntities(id: string, options?: ListEntitiesOptions) {
    return apiFetcher({
      url: `/analysis/${id}/entities?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<EntityList>
  }

  static useListEntities(
    id: string | null | undefined,
    options?: ListEntitiesOptions,
    swrOptions?: SWRConfiguration,
  ) {
    const url = `/analysis/${id}/entities?${this.paramsFromListOptions(options)}`
    return useSWR<EntityList>(
      url,
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
