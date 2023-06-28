import useSWR from 'swr'
import { apiFetch } from '@/client/fetch'
import { getAccessTokenOrRedirect } from '@/infra/token'
import { Organization } from './organization'

export enum SortBy {
  Name = 'name',
  DateCreated = 'date_created',
  DateModified = 'date_modified',
}

export enum SortOrder {
  Asc = 'asc',
  Desc = 'desc',
}

export type Group = {
  id: string
  name: string
  organization: Organization
  permission: string
  createTime: string
  updateTime?: string
}

export type List = {
  data: Group[]
  totalPages: number
  totalElements: number
  page: number
  size: number
}

export type ListOptions = {
  query?: string
  organizationId?: string
  size?: number
  page?: number
  sortBy?: SortBy
  sortOrder?: SortOrder
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

  static useList(options?: ListOptions, swrOptions?: any) {
    return useSWR<List>(
      `/groups?${this.paramsFromListOptions(options)}`,
      () => this.list(options),
      swrOptions
    )
  }

  static async list(options?: ListOptions): Promise<List> {
    return apiFetch(`/groups?${this.paramsFromListOptions(options)}`, {
      method: 'GET',
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

  static paramsFromListOptions(options?: ListOptions): URLSearchParams {
    const params: any = {}
    if (options?.query) {
      params.query = encodeURIComponent(options.query.toString())
    }
    if (options?.organizationId) {
      params.organization_id = options.organizationId.toString()
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
