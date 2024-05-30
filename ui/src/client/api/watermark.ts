import useSWR, { SWRConfiguration } from 'swr'
import { apiFetcher } from '@/client/fetcher'

export type Metadata = {
  isOutdated: boolean
}

export default class WatermarkAPI {
  static create(id: string) {
    return apiFetcher({
      url: `/watermarks/${id}`,
      method: 'POST',
    })
  }

  static delete(id: string) {
    return apiFetcher({
      url: `/watermarks/${id}`,
      method: 'DELETE',
    })
  }

  static useGetMetadata(
    id: string | null | undefined,
    swrOptions?: SWRConfiguration,
  ) {
    const url = `/watermarks/${id}/metadata`
    return useSWR<Metadata>(
      id ? url : null,
      () =>
        apiFetcher({
          url,
          method: 'GET',
          showError: false,
        }) as Promise<Metadata>,
      swrOptions,
    )
  }
}
