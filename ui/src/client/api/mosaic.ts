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

export type MosaicTile = {
  width: number
  height: number
  lastColWidth: number
  lastRowHeight: number
}

export type MosaicZoomLevel = {
  index: number
  width: number
  height: number
  rows: number
  cols: number
  scaleDownPercentage: number
  tile: MosaicTile
}

export type MosaicMetadata = {
  width: number
  height: number
  extension: string
  zoomLevels: MosaicZoomLevel[]
}

export class MosaicAPI {
  static create(fileId: string, showError = true) {
    return apiFetcher({
      url: `/mosaics/${fileId}`,
      method: 'POST',
      showError,
    })
  }

  static delete(fileId: string) {
    return apiFetcher({
      url: `/mosaics/${fileId}`,
      method: 'DELETE',
    })
  }

  static useGetMetadata(
    fileId: string | null | undefined,
    swrOptions?: SWRConfiguration,
  ) {
    const url = `/mosaics/${fileId}/metadata`
    return useSWR<MosaicMetadata>(
      fileId ? url : null,
      () =>
        apiFetcher({
          url,
          method: 'GET',
          showError: false,
        }) as Promise<MosaicMetadata>,
      swrOptions,
    )
  }

  static getMetadata(fileId: string) {
    return apiFetcher({
      url: `/mosaics/${fileId}/metadata`,
      method: 'GET',
    }) as Promise<MosaicMetadata>
  }
}
