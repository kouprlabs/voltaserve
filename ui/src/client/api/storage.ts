/* eslint-disable react-hooks/rules-of-hooks */
import useSWR from 'swr'
import { apiFetcher } from '@/client/fetch'

export type StorageUsage = {
  bytes: number
  maxBytes: number
  percentage: number
}

export default class StorageAPI {
  static useGetAccountUsage(swrOptions?: any) {
    return useSWR<StorageUsage>(
      `/storage/get_account_usage`,
      apiFetcher,
      swrOptions
    )
  }

  static useGetWorkspaceUsage(id: string, swrOptions?: any) {
    return useSWR<StorageUsage>(
      id
        ? `/storage/get_workspace_usage?${new URLSearchParams({
            id,
          })}`
        : null,
      apiFetcher,
      swrOptions
    )
  }

  static useGetFileUsage(id: string, swrOptions?: any) {
    return useSWR<StorageUsage>(
      id
        ? `/storage/get_file_usage?${new URLSearchParams({
            id,
          })}`
        : null,
      apiFetcher,
      swrOptions
    )
  }
}
