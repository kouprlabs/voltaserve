import useSWR, { SWRConfiguration } from 'swr'
import { apiFetcher } from '@/client/fetcher'

export type StorageUsage = {
  bytes: number
  maxBytes: number
  percentage: number
}

export default class StorageAPI {
  static useGetAccountUsage(swrOptions?: SWRConfiguration) {
    const url = `/storage/get_account_usage`
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
      ? `/storage/get_workspace_usage?${new URLSearchParams({
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
    const url = `/storage/get_file_usage?${new URLSearchParams({
      id: id!,
    })}`
    return useSWR<StorageUsage>(
      id ? url : null,
      () => apiFetcher({ url, method: 'GET' }) as Promise<StorageUsage>,
      swrOptions,
    )
  }
}
