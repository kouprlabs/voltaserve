/* eslint-disable react-hooks/rules-of-hooks */
import useSWR from 'swr'
import { getAccessTokenOrRedirect } from '@/infra/token'
import { apiFetch } from './fetch'
import { Group } from './group'
import { PermissionType } from './permission'
import { User } from './user'

export type Organization = {
  id: string
  name: string
  permission: PermissionType
  createTime: string
  updateTime?: string
}

export interface SearchOptions {
  text: string
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

  static useGetAllOrSearch(
    options?: { search?: SearchOptions },
    swrOptions?: any
  ) {
    if (options?.search) {
      return this.useSearch(options?.search, swrOptions)
    } else {
      return this.useGetAll(swrOptions)
    }
  }

  static useSearch(options: SearchOptions, swrOptions?: any) {
    return useSWR<Organization[]>(
      '/organizations/search',
      () => this.search(options),
      swrOptions
    )
  }

  static async search(
    options: SearchOptions
  ): Promise<Organization[]> {
    return apiFetch('/organizations/search', {
      method: 'POST',
      body: JSON.stringify(options),
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    }).then((result) => result.json())
  }

  static useGetAll(swrOptions?: any) {
    return useSWR<Organization[]>(
      '/organizations',
      () => this.getAll(),
      swrOptions
    )
  }

  static async getAll(): Promise<Organization[]> {
    return apiFetch('/organizations', {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    }).then((result) => result.json())
  }

  static async create(
    options: CreateOptions
  ): Promise<Organization> {
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

  static useGetMembers(id?: string, swrOptions?: any) {
    return useSWR<User[]>(
      id ? `/organizations/${id}/get_members` : null,
      () => this.getMembers(id as string),
      swrOptions
    )
  }

  static async getMembers(id: string): Promise<User[]> {
    return apiFetch(`/organizations/${id}/get_members`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    }).then((result) => result.json())
  }

  static useGetGroups(id?: string, swrOptions?: any) {
    return useSWR<Group[]>(
      id ? `/organizations/${id}/get_groups` : null,
      () => this.getGroups(id as string),
      swrOptions
    )
  }

  static async getGroups(id: string): Promise<Group[]> {
    return apiFetch(`/organizations/${id}/get_groups`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    }).then((result) => result.json())
  }

  static useSearchMembers(id: string, query: string, swrOptions?: any) {
    return useSWR<User[]>(
      id && query
        ? `/organizations/${id}/search_members?${new URLSearchParams({
            query,
          })}`
        : null,
      () => this.searchMembers(id, query),
      swrOptions
    )
  }

  static async searchMembers(id: string, query: string): Promise<User[]> {
    return apiFetch(
      `/organizations/${id}/search_members?${new URLSearchParams({
        query,
      })}`,
      {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
          'Content-Type': 'application/json',
        },
      }
    ).then((result) => result.json())
  }

  static async removeMember(
    id: string,
    options: RemoveMemberOptions
  ) {
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
