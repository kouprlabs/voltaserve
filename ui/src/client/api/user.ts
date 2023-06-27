import useSWR from 'swr'
import { apiFetch } from '@/client/fetch'
import { getAccessTokenOrRedirect } from '@/infra/token'

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
  size?: number
  page?: number
  sortBy?: SortBy
  sortOrder?: SortOrder
}

export default class UserAPI {
  static useList(options?: ListOptions, swrOptions?: any) {
    return useSWR<List>('/users', () => this.list(options), swrOptions)
  }

  static async list(options?: ListOptions): Promise<List> {
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
    return apiFetch(`/users?${new URLSearchParams(params)}`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    }).then((result) => result.json())
  }
}
