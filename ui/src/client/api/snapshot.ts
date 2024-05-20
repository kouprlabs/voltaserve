import useSWR, { SWRConfiguration } from 'swr'
import { apiFetcher } from '../fetcher'

export enum SortBy {
  Version = 'version',
  DateCreated = 'date_created',
  DateModified = 'date_modified',
}

export enum SortOrder {
  Asc = 'asc',
  Desc = 'desc',
}

export type Snapshot = {
  id: string
  version: number
  original: Download
  preview?: Download
  text?: Download
  thumbnail?: Thumbnail
  isActive: boolean
  createTime: string
  updateTime?: string
}

export type List = {
  data: Snapshot[]
  totalPages: number
  totalElements: number
  page: number
  size: number
}

export type ListOptions = {
  query?: string
  organizationId?: string
  size?: number
  page?: number
  sortBy?: SortBy
  sortOrder?: SortOrder
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

type ListQueryParams = {
  page?: string
  size?: string
  sort_by?: string
  sort_order?: string
  query?: string
}

export default class SnapshotAPI {
  static async list(fileId: string, options?: ListOptions) {
    return apiFetcher({
      url: `/files/${fileId}/snapshots?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<List>
  }

  static useList(
    fileId: string | undefined,
    options?: ListOptions,
    swrOptions?: SWRConfiguration,
  ) {
    const url = `/files/${fileId}/snapshots?${this.paramsFromListOptions(options)}`
    return useSWR<List | undefined>(
      fileId ? url : null,
      () => apiFetcher({ url, method: 'GET' }),
      swrOptions,
    )
  }

  static async activate(fileId: string, snapshotId: string) {
    return apiFetcher({
      url: `/files/${fileId}/snapshots/${snapshotId}/activate`,
      method: 'POST',
    }) as Promise<File>
  }

  static async delete(fileId: string, snapshotId: string) {
    return apiFetcher({
      url: `/files/${fileId}/snapshots/${snapshotId}`,
      method: 'DELETE',
    })
  }

  static paramsFromListOptions(options?: ListOptions): URLSearchParams {
    const params: ListQueryParams = {}
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
}
