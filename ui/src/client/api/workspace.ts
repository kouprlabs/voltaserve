/* eslint-disable react-hooks/rules-of-hooks */
import useSWR from 'swr'
import { apiFetch } from '@/client/fetch'
import { getAccessTokenOrRedirect } from '@/infra/token'
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
  isAutomaticOcrEnabled: boolean
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

export interface CreateOptions {
  name: string
  image?: string
  organizationId: string
  storageCapacity: number
}

export type ListOptions = {
  query?: string
  size?: number
  page?: number
  sortBy?: SortBy
  sortOrder?: SortOrder
}

export interface UpdateNameOptions {
  name: string
}

export interface UpdateIsAutomaticOcrEnabledOptions {
  isEnabled: boolean
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

  static useList(options?: ListOptions, swrOptions?: any) {
    return useSWR<List>(
      `/workspaces?${this.paramsFromListOptions(options)}`,
      () => this.list(options),
      swrOptions
    )
  }

  static async list(options?: ListOptions): Promise<List> {
    return apiFetch(`/workspaces?${this.paramsFromListOptions(options)}`, {
      method: 'GET',
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

  static async updateIsAutomaticOcrEnabled(
    id: string,
    options: UpdateIsAutomaticOcrEnabledOptions
  ): Promise<Workspace> {
    return apiFetch(`/workspaces/${id}/update_is_automatic_ocr_enabled`, {
      method: 'POST',
      body: JSON.stringify(options),
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
