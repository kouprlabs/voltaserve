import useSWR, { SWRConfiguration } from 'swr'
import { apiFetcher } from '@/client/fetcher'

export type Tile = {
  width: number
  height: number
  lastColWidth: number
  lastRowHeight: number
}

export type ZoomLevel = {
  index: number
  width: number
  height: number
  rows: number
  cols: number
  scaleDownPercentage: number
  tile: Tile
}

export type Metadata = {
  isOutdated: boolean
  width: number
  height: number
  extension: string
  zoomLevels: ZoomLevel[]
}

export default class MosaicAPI {
  static create(id: string) {
    return apiFetcher({
      url: `/mosaics/${id}`,
      method: 'POST',
    })
  }

  static delete(id: string) {
    return apiFetcher({
      url: `/mosaics/${id}`,
      method: 'DELETE',
    })
  }

  static useGetMetadata(
    id: string | null | undefined,
    swrOptions?: SWRConfiguration,
  ) {
    const url = `/mosaics/${id}/metadata`
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

  static async getMetadata(id: string) {
    return apiFetcher({
      url: `/mosaics/${id}/metadata`,
      method: 'GET',
    }) as Promise<Metadata>
  }
}
