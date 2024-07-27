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
      url: `/tasks?${paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<List>
  }

  static useList(options?: ListOptions, swrOptions?: SWRConfiguration) {
    const url = `/tasks?${paramsFromListOptions(options)}`
    return useSWR<List>(
      url,
      () => apiFetcher({ url, method: 'GET' }) as Promise<List>,
      swrOptions,
    )
  }

  static useGetCount(swrOptions?: SWRConfiguration) {
    const url = '/tasks/count'
    return useSWR<number>(
      url,
      () => apiFetcher({ url, method: 'GET' }) as Promise<number>,
      swrOptions,
    )
  }

  static async getCount() {
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
