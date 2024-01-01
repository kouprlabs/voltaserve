import useSWR from 'swr'
import { apiFetcher } from '@/client/fetcher'
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
  static create(options: CreateOptions): Promise<Group> {
    return apiFetcher({
      url: `/groups`,
      method: 'POST',
      body: JSON.stringify(options),
    })
  }

  static updateName(id: string, options: UpdateNameOptions): Promise<Group> {
    return apiFetcher({
      url: `/groups/${id}/update_name`,
      method: 'POST',
      body: JSON.stringify(options),
    })
  }

  static async updateImage(id: string, file: any): Promise<Group> {
    const formData = new FormData()
    formData.append('file', file)
    return apiFetcher({
      url: `/groups/${id}/update_image`,
      method: 'POST',
      body: formData,
    })
  }

  static async getById(id: string): Promise<Group> {
    return apiFetcher({
      url: `/groups/${id}`,
      method: 'GET',
    })
  }

  static useGetById(id: string | null | undefined, swrOptions?: any) {
    const url = `/groups/${id}`
    return useSWR<Group>(
      id ? url : null,
      () => apiFetcher({ url, method: 'GET' }),
      swrOptions,
    )
  }

  static async list(options?: ListOptions): Promise<List> {
    return apiFetcher({
      url: `/groups?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    })
  }

  static useList(options?: ListOptions, swrOptions?: any) {
    const url = `/groups?${this.paramsFromListOptions(options)}`
    return useSWR<List>(
      url,
      () => apiFetcher({ url, method: 'GET' }),
      swrOptions,
    )
  }

  static delete(id: string) {
    return apiFetcher({
      url: `/groups/${id}`,
      method: 'DELETE',
    })
  }

  static addMember(id: string, options: AddMemberOptions) {
    return apiFetcher({
      url: `/groups/${id}/add_member`,
      method: 'POST',
      body: JSON.stringify(options),
    })
  }

  static removeMember(id: string, options: RemoveMemberOptions) {
    return apiFetcher({
      url: `/groups/${id}/remove_member`,
      method: 'POST',
      body: JSON.stringify(options),
    })
  }

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
