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
import { File } from './file'

export type Snapshot = {
  id: string
  version: number
  status: SnapshotStatus
  original: SnapshotDownloadable
  preview?: SnapshotDownloadable
  ocr?: SnapshotDownloadable
  text?: SnapshotDownloadable
  thumbnail?: SnapshotDownloadable
  summary?: string
  intent?: string
  language?: string
  capabilities: SnapshotCapabilities
  isActive: boolean
  task?: SnapshotTaskInfo
  createTime: string
  updateTime?: string
}

export enum SnapshotStatus {
  Waiting = 'waiting',
  Processing = 'processing',
  Ready = 'ready',
  Error = 'error',
}

export type SnapshotList = {
  data: Snapshot[]
  totalPages: number
  totalElements: number
  page: number
  size: number
}

export type SnapshotListOptions = {
  fileId: string
  query?: string
  organizationId?: string
  size?: number
  page?: number
  sortBy?: SnapshotSortBy
  sortOrder?: SnapshotSortOrder
}

export enum SnapshotSortBy {
  Version = 'version',
  DateCreated = 'date_created',
  DateModified = 'date_modified',
}

export enum SnapshotSortOrder {
  Asc = 'asc',
  Desc = 'desc',
}

export type SnapshotTaskInfo = {
  id: string
  isPending: boolean
}

export type SnapshotCapabilities = {
  original: boolean
  preview: boolean
  ocr: boolean
  text: boolean
  summary: boolean
  entities: boolean
  mosaic: boolean
  thumbnail: boolean
}

export type SnapshotDownloadable = {
  extension?: string
  size?: number
  image?: ImageProps
  document?: DocumentProps
}

export type ImageProps = {
  width: number
  height: number
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

export type SnapshotListQueryParams = {
  file_id: string
  page?: string
  size?: string
  sort_by?: string
  sort_order?: string
  query?: string
}

export type SnapshotLanguage = {
  id: string
  iso6393: string
  name: string
}

export class SnapshotAPI {
  static list(options: SnapshotListOptions) {
    return apiFetcher({
      url: `/snapshots?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<SnapshotList>
  }

  static useList(options: SnapshotListOptions, swrOptions?: SWRConfiguration) {
    const url = `/snapshots?${this.paramsFromListOptions(options)}`
    return useSWR<SnapshotList | undefined>(
      url,
      () => apiFetcher({ url, method: 'GET' }),
      swrOptions,
    )
  }

  static useGetLanguages(swrOptions?: SWRConfiguration) {
    const url = `/snapshots/languages`
    return useSWR<SnapshotLanguage[]>(
      url,
      () => apiFetcher({ url, method: 'GET' }) as Promise<SnapshotLanguage[]>,
      swrOptions,
    )
  }

  static paramsFromListOptions(options: SnapshotListOptions): URLSearchParams {
    const params: SnapshotListQueryParams = { file_id: options.fileId }
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
