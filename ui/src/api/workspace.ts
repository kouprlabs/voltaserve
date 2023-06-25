/* eslint-disable react-hooks/rules-of-hooks */
import useSWR from 'swr'
import { getAccessTokenOrRedirect } from '@/infra/token'
import { apiFetch } from './fetch'
import { Organization } from './organization'
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

export type Workspace = {
  id: string
  name: string
  permission: PermissionType
  storageCapacity: number
  rootId: string
  organization: Organization
  createTime: string
  updateTime?: string
}

export type List = {
  data: Workspace[]
  totalPages: number
  totalElements: number
  page: number
  size: number
}

export interface SearchOptions {
  text: string
}

export type ListOptions = {
  size?: number
  page?: number
  sortBy?: SortBy
  sortOrder?: SortOrder
}

export interface CreateOptions {
  name: string
  image?: string
  organizationId: string
  storageCapacity: number
}

export interface UpdateNameOptions {
  name: string
}

export interface StorageCapacityOptions {
  storageCapacity: number
}

export default class WorkspaceAPI {
  static useGetById(id: string, swrOptions?: any) {
    return useSWR<Workspace>(
      id ? `/workspaces/${id}` : null,
      () => this.getById(id),
      swrOptions
    )
  }

  static async getById(id: string): Promise<Workspace> {
    return apiFetch(`/workspaces/${id}`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    }).then((result) => result.json())
  }

  static useListOrSearch(
    options?: { search?: SearchOptions; list?: ListOptions },
    swrOptions?: any
  ) {
    if (options?.search) {
      return this.useSearch(options?.search, swrOptions)
    } else {
      return this.useList(options?.list, swrOptions)
    }
  }

  static useList(options?: ListOptions, swrOptions?: any) {
    return useSWR<List>('/workspaces', () => this.list(options), swrOptions)
  }

  static async list(options?: ListOptions): Promise<List> {
    const params: any = {}
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
    return apiFetch(`/workspaces?${new URLSearchParams(params)}`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    }).then((result) => result.json())
  }

  static useSearch(options: SearchOptions, swrOptions?: any) {
    return useSWR<List>(
      '/workspaces/search',
      () => this.search(options),
      swrOptions
    )
  }

  static async search(options: SearchOptions): Promise<List> {
    return apiFetch('/workspaces/search', {
      method: 'POST',
      body: JSON.stringify(options),
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    }).then((result) => result.json())
  }

  static async create(options: CreateOptions): Promise<Workspace> {
    return apiFetch('/workspaces', {
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
  ): Promise<Workspace> {
    return apiFetch(`/workspaces/${id}/update_name`, {
      method: 'POST',
      body: JSON.stringify(options),
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    }).then((result) => result.json())
  }

  static async updateStorageCapacity(
    id: string,
    options: StorageCapacityOptions
  ): Promise<Workspace> {
    return apiFetch(`/workspaces/${id}/update_storage_capacity`, {
      method: 'POST',
      body: JSON.stringify(options),
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    }).then((result) => result.json())
  }

  static async updateImage(id: string, file: any): Promise<Workspace> {
    const formData = new FormData()
    formData.append('file', file)
    return apiFetch(`/workspaces/${id}/update_image`, {
      method: 'POST',
      body: formData,
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    }).then((result) => result.json())
  }

  static async delete(id: string) {
    return apiFetch(`/workspaces/${id}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    })
  }
}
