/* eslint-disable react-hooks/rules-of-hooks */
import useSWR from 'swr'
import { getAccessTokenOrRedirect } from '@/infra/token'
import { apiFetch } from './fetch'
import { Organization } from './organization'
import { PermissionType } from './permission'

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

export interface WorkspaceSearchOptions {
  text: string
}

export interface WorkspaceCreateOptions {
  name: string
  image?: string
  organizationId: string
  storageCapacity: number
}

export interface WorkspaceUpdateNameOptions {
  name: string
}

export interface WorkspaceStorageCapacityOptions {
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

  static useGetAllOrSearch(
    options?: { search?: WorkspaceSearchOptions },
    swrOptions?: any
  ) {
    if (options?.search) {
      return this.useSearch(options?.search, swrOptions)
    } else {
      return this.useGetAll(swrOptions)
    }
  }

  static useGetAll(swrOptions?: any) {
    return useSWR<Workspace[]>('/workspaces', () => this.getAll(), swrOptions)
  }

  static async getAll(): Promise<Workspace[]> {
    return apiFetch('/workspaces', {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    }).then((result) => result.json())
  }

  static useSearch(options: WorkspaceSearchOptions, swrOptions?: any) {
    return useSWR<Workspace[]>(
      '/workspaces/search',
      () => this.search(options),
      swrOptions
    )
  }

  static async search(options: WorkspaceSearchOptions): Promise<Workspace[]> {
    return apiFetch('/workspaces/search', {
      method: 'POST',
      body: JSON.stringify(options),
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    }).then((result) => result.json())
  }

  static async create(options: WorkspaceCreateOptions): Promise<Workspace> {
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
    options: WorkspaceUpdateNameOptions
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
    options: WorkspaceStorageCapacityOptions
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
