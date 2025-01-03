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
import { apiFetcher } from '../fetcher'
import { File } from './file'

export type Snapshot = {
  id: string
  version: number
  status: Status
  original: Download
  preview?: Download
  ocr?: Download
  text?: Download
  entities?: Download
  mosaic?: Download
  thumbnail?: Download
  language?: string
  isActive: boolean
  task?: TaskInfo
  createTime: string
  updateTime?: string
}

export enum Status {
  Waiting = 'waiting',
  Processing = 'processing',
  Ready = 'ready',
  Error = 'error',
}

export type List = {
  data: Snapshot[]
  totalPages: number
  totalElements: number
  page: number
  size: number
}

export type ListOptions = {
  fileId: string
  query?: string
  organizationId?: string
  size?: number
  page?: number
  sortBy?: SortBy
  sortOrder?: SortOrder
}

export enum SortBy {
  Version = 'version',
  DateCreated = 'date_created',
  DateModified = 'date_modified',
}

export enum SortOrder {
  Asc = 'asc',
  Desc = 'desc',
}

export type TaskInfo = {
  id: string
  isPending: boolean
}

export type Download = {
  extension?: string
  size?: number
  image?: ImageProps
  document?: DocumentProps
}

export type ImageProps = {
  width: number
  height: number
  zoomLevels?: ZoomLevel[]
}

export type DocumentProps = {
  pages?: PagesProps
  thumbnails?: ThumbnailsProps
}

export type PagesProps = {
  count: number
  extension: string
}

export type ThumbnailsProps = {
  extension: string
}

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

export type ListQueryParams = {
  file_id: string
  page?: string
  size?: string
  sort_by?: string
  sort_order?: string
  query?: string
}

export default class SnapshotAPI {
  static list(options: ListOptions) {
    return apiFetcher({
      url: `/snapshots?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<List>
  }

  static useList(options: ListOptions, swrOptions?: SWRConfiguration) {
    const url = `/snapshots?${this.paramsFromListOptions(options)}`
    return useSWR<List | undefined>(
      url,
      () => apiFetcher({ url, method: 'GET' }),
      swrOptions,
    )
  }

  static paramsFromListOptions(options: ListOptions): URLSearchParams {
    const params: ListQueryParams = { file_id: options.fileId }
    if (options?.query) {
      params.query = encodeURIComponent(options.query.toString())
    }
    if (options?.page) {
      params.page = options.page.toString()
    }
    if (options?.size) {
      params.size = options.size.toString()
    }
    if (options?.sortBy) {
      params.sort_by = options.sortBy.toString()
    }
    if (options?.sortOrder) {
      params.sort_order = options.sortOrder.toString()
    }
    return new URLSearchParams(params)
  }

  static activate(id: string) {
    return apiFetcher({
      url: `/snapshots/${id}/activate`,
      method: 'POST',
    }) as Promise<File>
  }

  static async detach(id: string) {
    return apiFetcher({
      url: `/snapshots/${id}/detach`,
      method: 'POST',
    })
  }
}
