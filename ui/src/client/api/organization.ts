/* eslint-disable react-hooks/rules-of-hooks */
import useSWR from 'swr'
import { apiFetcher } from '@/client/fetcher'
import { PermissionType } from './permission'

export enum SortBy {
  Name = 'name',
  DateCreated = 'date_created',
  DateModified = 'date_modified',
}

export enum SortOrder {
  Asc = 'asc',
  Desc = 'desc',
}

export type Organization = {
  id: string
  name: string
  permission: PermissionType
  createTime: string
  updateTime?: string
}

export type List = {
  data: Organization[]
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

export type CreateOptions = {
  name: string
  image?: string
}

export type UpdateNameOptions = {
  name: string
}

export type RemoveMemberOptions = {
  userId: string
}

export default class OrganizationAPI {
  static async getById(id: string): Promise<Organization> {
    return apiFetcher({
      url: `/organizations/${id}`,
      method: 'GET',
    })
  }

  static useGetById(id: string | null | undefined, swrOptions?: any) {
    const url = `/organizations/${id}`
    return useSWR<Organization>(
      id ? url : null,
      () => apiFetcher({ url, method: 'GET' }),
      swrOptions,
    )
  }

  static async list(options?: ListOptions): Promise<List> {
    return apiFetcher({
      url: `/organizations?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    })
  }

  static useList(options?: ListOptions, swrOptions?: any) {
    const url = `/organizations?${this.paramsFromListOptions(options)}`
    return useSWR<List>(
      url,
      () => apiFetcher({ url, method: 'GET' }),
      swrOptions,
    )
  }

  static async create(options: CreateOptions): Promise<Organization> {
    return apiFetcher({
      url: `/organizations`,
      method: 'POST',
      body: JSON.stringify(options),
    })
  }

  static async updateName(
    id: string,
    options: UpdateNameOptions,
  ): Promise<Organization> {
    return apiFetcher({
      url: `/organizations/${id}/update_name`,
      method: 'POST',
      body: JSON.stringify(options),
    })
  }

  static async updateImage(id: string, file: any): Promise<Organization> {
    const formData = new FormData()
    formData.append('file', file)
    return apiFetcher({
      url: `/organizations/${id}/update_image`,
      method: 'POST',
      body: formData,
    })
  }

  static async delete(id: string) {
    return apiFetcher({
      url: `/organizations/${id}`,
      method: 'DELETE',
    })
  }

  static async leave(id: string) {
    return apiFetcher({
      url: `/organizations/${id}/leave`,
      method: 'POST',
    })
  }

  static async removeMember(id: string, options: RemoveMemberOptions) {
    return apiFetcher({
      url: `/organizations/${id}/remove_member`,
      method: 'POST',
      body: JSON.stringify(options),
    })
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
