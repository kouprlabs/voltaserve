/* eslint-disable react-hooks/rules-of-hooks */
import useSWR from 'swr'
import { apiFetcher } from '@/client/fetcher'
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

export interface StorageCapacityOptions {
  storageCapacity: number
}

export default class WorkspaceAPI {
  static async getById(id: string): Promise<Workspace> {
    return apiFetcher({
      url: `/workspaces/${id}`,
      method: 'GET',
    })
  }

  static useGetById(id: string | null | undefined, swrOptions?: any) {
    const url = `/workspaces/${id}`
    return useSWR<Workspace>(
      id ? url : null,
      () => apiFetcher({ url, method: 'GET' }),
      swrOptions,
    )
  }

  static useList(options?: ListOptions, swrOptions?: any) {
    const url = `/workspaces?${this.paramsFromListOptions(options)}`
    return useSWR<List>(
      url,
      () => apiFetcher({ url, method: 'GET' }),
      swrOptions,
    )
  }

  static async list(options?: ListOptions): Promise<List> {
    return apiFetcher({
      url: `/workspaces?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    })
  }

  static async create(options: CreateOptions): Promise<Workspace> {
    return apiFetcher({
      url: '/workspaces',
      method: 'POST',
      body: JSON.stringify(options),
    })
  }

  static async updateName(
    id: string,
    options: UpdateNameOptions,
  ): Promise<Workspace> {
    return apiFetcher({
      url: `/workspaces/${id}/update_name`,
      method: 'POST',
      body: JSON.stringify(options),
    })
  }

  static async updateStorageCapacity(
    id: string,
    options: StorageCapacityOptions,
  ): Promise<Workspace> {
    return apiFetcher({
      url: `/workspaces/${id}/update_storage_capacity`,
      method: 'POST',
      body: JSON.stringify(options),
    })
  }

  static async updateImage(id: string, file: any): Promise<Workspace> {
    const formData = new FormData()
    formData.append('file', file)
    return apiFetcher({
      url: `/workspaces/${id}/update_image`,
      method: 'POST',
      body: formData,
    })
  }

  static async delete(id: string) {
    return apiFetcher({
      url: `/workspaces/${id}`,
      method: 'DELETE',
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
