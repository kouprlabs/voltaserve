import useSWR, { SWRConfiguration } from 'swr'
import { apiFetcher } from '../fetcher'

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
  watermark?: Download
  thumbnail?: Thumbnail
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
  extension: string
  size: number
  image?: ImageProps
}

export type ImageProps = {
  width: number
  height: number
}

export type Thumbnail = {
  base64: string
  width: number
  height: number
}

export type ListQueryParams = {
  file_id: string
  page?: string
  size?: string
  sort_by?: string
  sort_order?: string
  query?: string
}

export type ActivateOptions = {
  fileId: string
}

export type DetachOptions = {
  fileId: string
}

export default class SnapshotAPI {
  static async list(options: ListOptions) {
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

  static async activate(id: string, options: ActivateOptions) {
    return apiFetcher({
      url: `/snapshots/${id}/activate`,
      method: 'POST',
      body: JSON.stringify(options),
    }) as Promise<File>
  }

  static async detach(id: string, options: DetachOptions) {
    return apiFetcher({
      url: `/snapshots/${id}/detach`,
      method: 'POST',
      body: JSON.stringify(options),
    })
  }
}
