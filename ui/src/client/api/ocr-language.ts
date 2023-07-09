import useSWR from 'swr'
import { apiFetch } from '@/client/fetch'
import { getAccessTokenOrRedirect } from '@/infra/token'

export enum SortBy {
  Id = 'id',
  Iso639Pt3 = 'iso639_3',
}

export enum SortOrder {
  Asc = 'asc',
  Desc = 'desc',
}

export type OcrLanguage = {
  id: string
  iso639pt3: string
}

export type List = {
  data: OcrLanguage[]
  totalPages: number
  totalElements: number
  page: number
  size: number
}

export type ListOptions = {
  query?: string
  size?: number
  page?: number
  sortBy?: SortBy
  sortOrder?: SortOrder
}

export default class OcrLanguageAPI {
  static useList(options?: ListOptions, swrOptions?: any) {
    return useSWR<List>(
      `/users?${this.paramsFromListOptions(options)}`,
      () => this.list(options),
      swrOptions
    )
  }

  static async list(options?: ListOptions): Promise<List> {
    return apiFetch(`/ocr_languages?${this.paramsFromListOptions(options)}`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    }).then((result) => result.json())
  }

  static paramsFromListOptions(options?: ListOptions): URLSearchParams {
    const params: any = {}
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
