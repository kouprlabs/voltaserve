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

export enum TaskSortBy {
  Name = 'name',
  Status = 'status',
  DateCreated = 'date_created',
  DateModified = 'date_modified',
}

export enum TaskSortOrder {
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
  status: TaskStatus
  payload?: TaskPayload
  createTime: string
  updateTime?: string
}

export enum TaskStatus {
  Waiting = 'waiting',
  Running = 'running',
  Success = 'success',
  Error = 'error',
}

export type TaskPayload = {
  object?: string
}

export type TaskList = {
  data: Task[]
  totalPages: number
  totalElements: number
  page: number
  size: number
}

export type TaskListOptions = {
  query?: string
  size?: number
  page?: number
  sortBy?: TaskSortBy
  sortOrder?: TaskSortOrder
}

type TaskListQueryParams = {
  page?: string
  size?: string
  sort_by?: string
  sort_order?: string
  query?: string
}

export class TaskAPI {
  static useGet(id: string | null | undefined, swrOptions?: SWRConfiguration) {
    const url = `/tasks/${id}`
    return useSWR<Task>(
      id ? url : null,
      () => apiFetcher({ url, method: 'GET' }) as Promise<Task>,
      swrOptions,
    )
  }

  static list(options?: TaskListOptions) {
    return apiFetcher({
      url: `/tasks?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<TaskList>
  }

  static useList(options?: TaskListOptions, swrOptions?: SWRConfiguration) {
    const url = `/tasks?${this.paramsFromListOptions(options)}`
    return useSWR<TaskList>(
      url,
      () => apiFetcher({ url, method: 'GET' }) as Promise<TaskList>,
      swrOptions,
    )
  }

  static paramsFromListOptions(options?: TaskListOptions): URLSearchParams {
    const params: TaskListQueryParams = {}
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
