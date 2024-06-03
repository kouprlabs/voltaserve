import useSWR, { SWRConfiguration } from 'swr'
import { apiFetcher } from '@/client/fetcher'
import { PermissionType } from './permission'

export enum SortBy {
  Name = 'name',
  DateCreated = 'date_created',
  DateModified = 'date_modified',
}

export enum SortOrder {
  Asc = 'asc',
  Desc = 'desc',
}

export type Task = {
  id: string
  name: string
  error?: string
  percentage: number
  isComplete: boolean
  isIndeterminate: boolean
  userId: string
}

export type List = {
  data: Task[]
  totalPages: number
  totalElements: number
  page: number
  size: number
}

export type ListOptions = {
  query?: string
  size?: number
  page?: number
  sortBy?: SortBy
  sortOrder?: SortOrder
}

type ListQueryParams = {
  page?: string
  size?: string
  sort_by?: string
  sort_order?: string
  query?: string
}

export default class TaskAPI {
  static useGet(id: string | null | undefined, swrOptions?: SWRConfiguration) {
    const url = `/tasks/${id}`
    return useSWR<Task>(
      id ? url : null,
      () => apiFetcher({ url, method: 'GET' }) as Promise<Task>,
      swrOptions,
    )
  }

  static async list(options?: ListOptions) {
    return apiFetcher({
      url: `/tasks?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<List>
  }

  static useList(options?: ListOptions, swrOptions?: SWRConfiguration) {
    const url = `/tasks?${this.paramsFromListOptions(options)}`
    return useSWR<List>(
      url,
      () => apiFetcher({ url, method: 'GET' }) as Promise<List>,
      swrOptions,
    )
  }

  static paramsFromListOptions(options?: ListOptions): URLSearchParams {
    const params: ListQueryParams = {}
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
