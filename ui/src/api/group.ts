import useSWR from 'swr'
import { getAccessTokenOrRedirect } from '@/infra/token'
import { apiFetch } from './fetch'
import { Organization } from './organization'
import { User } from './user'

export type Group = {
  id: string
  name: string
  organization: Organization
  permission: string
  createTime: string
  updateTime?: string
}

export interface SearchOptions {
  text: string
}

export type CreateOptions = {
  name: string
  image?: string
  organizationId: string
}

export type UpdateNameOptions = {
  name: string
}

export type AddMemberOptions = {
  userId: string
}

export type RemoveMemberOptions = {
  userId: string
}

export default class GroupAPI {
  static create = (options: CreateOptions): Promise<Group> =>
    apiFetch(`/groups`, {
      method: 'POST',
      body: JSON.stringify(options),
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    }).then((result) => result.json())

  static updateName = (
    id: string,
    options: UpdateNameOptions
  ): Promise<Group> =>
    apiFetch(`/groups/${id}/update_name`, {
      method: 'POST',
      body: JSON.stringify(options),
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    }).then((result) => result.json())

  static async updateImage(id: string, file: any): Promise<Group> {
    const formData = new FormData()
    formData.append('file', file)
    return apiFetch(`/groups/${id}/update_image`, {
      method: 'POST',
      body: formData,
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    }).then((result) => result.json())
  }

  static useGetById = (id: string, swrOptions?: any) =>
    useSWR<Group>(
      id ? `/groups/${id}` : null,
      () => this.getById(id),
      swrOptions
    )

  static async getById(id: string): Promise<Group> {
    return apiFetch(`/groups/${id}`, {
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

  static useGetAll = (swrOptions?: any) =>
    useSWR<Group[]>(`/groups`, () => this.getAll(), swrOptions)

  static async getAll(): Promise<Group[]> {
    return apiFetch('/groups', {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    }).then((result) => result.json())
  }

  static useSearch(options: SearchOptions, swrOptions?: any) {
    return useSWR<Group[]>(
      '/groups/search',
      () => this.search(options),
      swrOptions
    )
  }

  static async search(options: SearchOptions): Promise<Group[]> {
    return apiFetch('/groups/search', {
      method: 'POST',
      body: JSON.stringify(options),
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    }).then((result) => result.json())
  }

  static delete = (id: string) =>
    apiFetch(`/groups/${id}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    })

  static useGetMembers = (id: string, swrOptions?: any) =>
    useSWR<User[]>(
      id ? `/groups/${id}/get_members` : null,
      () => this.getMembers(id),
      swrOptions
    )

  static async getMembers(id: string): Promise<User[]> {
    return apiFetch(`/groups/${id}/get_members`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    }).then((result) => result.json())
  }

  static useSearchMembers = (id: string, query: string, swrOptions?: any) =>
    useSWR<User[]>(
      id && query
        ? `/groups/${id}/search_members?${new URLSearchParams({
            query,
          })}`
        : null,
      () => this.searchMembers(id, query),
      swrOptions
    )

  static async searchMembers(id: string, query: string): Promise<User[]> {
    return apiFetch(
      `/groups/${id}/search_members?${new URLSearchParams({
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

  static addMember = (id: string, options: AddMemberOptions) =>
    apiFetch(`/groups/${id}/add_member`, {
      method: 'POST',
      body: JSON.stringify(options),
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    })

  static removeMember = (id: string, options: RemoveMemberOptions) =>
    apiFetch(`/groups/${id}/remove_member`, {
      method: 'POST',
      body: JSON.stringify(options),
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    })

  static useGetAvailableUsers = (id: string, swrOptions?: any) =>
    useSWR<User[]>(
      id ? `/groups/${id}/get_available_users` : null,
      () => this.getAvailableUsers(id),
      swrOptions
    )

  static async getAvailableUsers(id: string): Promise<User[]> {
    return apiFetch(`/groups/${id}/get_available_users`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    }).then((result) => result.json())
  }
}
