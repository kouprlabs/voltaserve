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

export type StorageUsage = {
  bytes: number
  maxBytes: number
  percentage: number
}

export default class StorageAPI {
  static useGetAccountUsage(swrOptions?: SWRConfiguration) {
    const url = `/storage/account_usage`
    return useSWR<StorageUsage>(
      url,
      () => apiFetcher({ url, method: 'GET' }) as Promise<StorageUsage>,
      swrOptions,
    )
  }

  static useGetWorkspaceUsage(
    id: string | null | undefined,
    swrOptions?: SWRConfiguration,
  ) {
    const url = id
      ? `/storage/workspace_usage?${new URLSearchParams({
          id,
        })}`
      : null
    return useSWR<StorageUsage>(
      url,
      () => apiFetcher({ url: url!, method: 'GET' }) as Promise<StorageUsage>,
      swrOptions,
    )
  }

  static useGetFileUsage(
    id: string | null | undefined,
    swrOptions?: SWRConfiguration,
  ) {
    const url = `/storage/file_usage?${new URLSearchParams({
      id: id!,
    })}`
    return useSWR<StorageUsage>(
      id ? url : null,
      () => apiFetcher({ url, method: 'GET' }) as Promise<StorageUsage>,
      swrOptions,
    )
  }
}
