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

export type Info = {
  isAvailable: boolean
  metadata?: Metadata
}

export type Metadata = {
  isOutdated: boolean
  width: number
  height: number
  extension: string
  zoomLevels: ZoomLevel[]
}

export default class MosaicAPI {
  static create(id: string, showError = true) {
    return apiFetcher({
      url: `/mosaics/${id}`,
      method: 'POST',
      showError,
    })
  }

  static delete(id: string) {
    return apiFetcher({
      url: `/mosaics/${id}`,
      method: 'DELETE',
    })
  }

  static useGetInfo(
    id: string | null | undefined,
    swrOptions?: SWRConfiguration,
  ) {
    const url = `/mosaics/${id}/info`
    return useSWR<Info>(
      id ? url : null,
      () =>
        apiFetcher({
          url,
          method: 'GET',
          showError: false,
        }) as Promise<Info>,
      swrOptions,
    )
  }

  static async getInfo(id: string) {
    return apiFetcher({
      url: `/mosaics/${id}/info`,
      method: 'GET',
    }) as Promise<Info>
  }
}
