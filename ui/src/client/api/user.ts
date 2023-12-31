import useSWR from 'swr'
import { apiFetcher } from '@/client/fetcher'

export enum SortBy {
  Email = 'email',
  FullName = 'full_name',
}

export enum SortOrder {
  Asc = 'asc',
  Desc = 'desc',
}

export type User = {
  id: string
  username: string
  email: string
  fullName: string
  picture?: string
}

export type List = {
  data: User[]
  totalPages: number
  totalElements: number
  page: number
  size: number
}

export type ListOptions = {
  query?: string
  organizationId?: string
  groupId?: string
  nonGroupMembersOnly?: boolean
  size?: number
  page?: number
  sortBy?: SortBy
  sortOrder?: SortOrder
}

export default class UserAPI {
  static async list(options?: ListOptions): Promise<List> {
    return apiFetcher({
      url: `/users?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    })
  }

  static useList(options?: ListOptions, swrOptions?: any) {
    const url = `/users?${this.paramsFromListOptions(options)}`
    return useSWR<List>(
      url,
      () => apiFetcher({ url, method: 'GET' }),
      swrOptions,
    )
  }

  static paramsFromListOptions(options?: ListOptions): URLSearchParams {
    const params: any = {}
    if (options?.query) {
      params.query = encodeURIComponent(options.query.toString())
    }
    if (options?.organizationId) {
      params.organization_id = options.organizationId.toString()
    }
    if (options?.groupId) {
      params.group_id = options.groupId.toString()
    }
    if (options?.nonGroupMembersOnly) {
      params.non_group_members_only = options.nonGroupMembersOnly.toString()
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
