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
import { apiFetcher } from '@/client/fetcher'

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
  percentage?: number
  isIndeterminate: boolean
  userId: string
  status: Status
  payload?: TaskPayload
}

export enum Status {
  Waiting = 'waiting',
  Running = 'running',
  Success = 'success',
  Error = 'error',
}

export type TaskPayload = {
  object?: string
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

  static list(options?: ListOptions) {
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

  static useGetCount(swrOptions?: SWRConfiguration) {
    const url = '/tasks/count'
    return useSWR<number>(
      url,
      () => apiFetcher({ url, method: 'GET' }) as Promise<number>,
      swrOptions,
    )
  }

  static getCount() {
    return apiFetcher({
      url: `/tasks/count`,
      method: 'GET',
    }) as Promise<number>
  }

  static async dismiss(id: string) {
    return apiFetcher({
      url: `/tasks/${id}/dismiss`,
      method: 'POST',
    })
  }

  static async dismissAll() {
    return apiFetcher({
      url: `/tasks/dismiss`,
      method: 'POST',
    })
  }
}
