/* eslint-disable react-hooks/rules-of-hooks */
import useSWR from 'swr'
import { apiFetcher } from '@/client/fetcher'

export type StorageUsage = {
  bytes: number
  maxBytes: number
  percentage: number
}

export default class StorageAPI {
  static useGetAccountUsage(swrOptions?: any) {
    const url = `/storage/get_account_usage`
    return useSWR<StorageUsage>(
      url,
      () => apiFetcher({ url, method: 'GET' }),
      swrOptions,
    )
  }

  static useGetWorkspaceUsage(id?: string, swrOptions?: any) {
    const url = id
      ? `/storage/get_workspace_usage?${new URLSearchParams({
          id,
        })}`
      : null
    return useSWR<StorageUsage>(
      url,
      () => apiFetcher({ url: url!, method: 'GET' }),
      swrOptions,
    )
  }

  static useGetFileUsage(id: string, swrOptions?: any) {
    const url = `/storage/get_file_usage?${new URLSearchParams({
      id,
    })}`
    return useSWR<StorageUsage>(
      id ? url : null,
      () => apiFetcher({ url, method: 'GET' }),
      swrOptions,
    )
  }
}
