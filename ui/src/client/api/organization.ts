/* eslint-disable react-hooks/rules-of-hooks */
import useSWR from 'swr'
import { apiFetch } from '@/client/fetch'
import { User } from '@/client/idp/user'
import { getAccessTokenOrRedirect } from '@/infra/token'
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
  static useGetById(id: string, swrOptions?: any) {
    return useSWR<Organization>(
      id ? `/organizations/${id}` : null,
      () => this.getById(id),
      swrOptions
    )
  }

  static async getById(id: string): Promise<Organization> {
    return apiFetch(`/organizations/${id}`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    }).then((result) => result.json())
  }

  static useList(options?: ListOptions, swrOptions?: any) {
    return useSWR<List>('/organizations', () => this.list(options), swrOptions)
  }

  static async list(options?: ListOptions): Promise<List> {
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
    return apiFetch(`/organizations?${new URLSearchParams(params)}`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    }).then((result) => result.json())
  }

  static async create(options: CreateOptions): Promise<Organization> {
    return apiFetch(`/organizations`, {
      method: 'POST',
      body: JSON.stringify(options),
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    }).then((result) => result.json())
  }

  static async updateName(
    id: string,
    options: UpdateNameOptions
  ): Promise<Organization> {
    return apiFetch(`/organizations/${id}/update_name`, {
      method: 'POST',
      body: JSON.stringify(options),
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    }).then((result) => result.json())
  }

  static async updateImage(id: string, file: any): Promise<Organization> {
    const formData = new FormData()
    formData.append('file', file)
    return apiFetch(`/organizations/${id}/update_image`, {
      method: 'POST',
      body: formData,
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    }).then((result) => result.json())
  }

  static async delete(id: string) {
    return apiFetch(`/organizations/${id}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    })
  }

  static async leave(id: string) {
    return apiFetch(`/organizations/${id}/leave`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    })
  }

  static async removeMember(id: string, options: RemoveMemberOptions) {
    return apiFetch(`/organizations/${id}/remove_member`, {
      method: 'POST',
      body: JSON.stringify(options),
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    })
  }
}
