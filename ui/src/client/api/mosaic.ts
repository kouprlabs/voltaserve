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
import { Snapshot } from './snapshot'

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
  isOutdated: boolean
  snapshot?: Snapshot
  metadata?: Metadata
}

export type Metadata = {
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

  static useGetInfo(id: string | null | undefined, swrOptions?: SWRConfiguration) {
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

  static getInfo(id: string) {
    return apiFetcher({
      url: `/mosaics/${id}/info`,
      method: 'GET',
    }) as Promise<Info>
  }
}
