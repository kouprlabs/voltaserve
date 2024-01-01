/* eslint-disable react-hooks/rules-of-hooks */
import useSWR from 'swr'
import { apiFetcher } from '@/client/fetcher'
import { User } from '@/client/idp/user'
import { Organization } from './organization'

export enum SortBy {
  Email = 'email',
  DateCreated = 'date_created',
  DateModified = 'date_modified',
}

export enum SortOrder {
  Asc = 'asc',
  Desc = 'desc',
}

export type InvitationStatus = 'pending' | 'accepted' | 'declined'

export type Invitation = {
  id: string
  owner: User
  email: string[]
  organization: Organization
  status: InvitationStatus
  createTime: string
  updateTime?: string
}

export type List = {
  data: Invitation[]
  totalPages: number
  totalElements: number
  page: number
  size: number
}

export type CreateOptions = {
  organizationId: string
  emails: string[]
}

export type ListOptions = {
  organizationId?: string
  size?: number
  page?: number
  sortBy?: SortBy
  sortOrder?: SortOrder
}

export default class InvitationAPI {
  static async create(options: CreateOptions) {
    return apiFetcher({
      url: `/invitations`,
      method: 'POST',
      body: JSON.stringify(options),
    })
  }

  static useGetIncoming(options?: ListOptions, swrOptions?: any) {
    const url = `/invitations/get_incoming?${this.paramsFromListOptions(
      options,
    )}`
    return useSWR<List>(
      url,
      () => apiFetcher({ url, method: 'GET' }),
      swrOptions,
    )
  }

  static useGetOutgoing(options?: ListOptions, swrOptions?: any) {
    const url = `/invitations/get_outgoing?${this.paramsFromListOptions(
      options,
    )}`
    return useSWR<List>(
      options?.organizationId ? url : null,
      () => apiFetcher({ url, method: 'GET' }),
      swrOptions,
    )
  }

  static async delete(id: string) {
    return apiFetcher({
      url: `/invitations/${id}`,
      method: 'DELETE',
    })
  }

  static async resend(id: string) {
    return apiFetcher({
      url: `/invitations/${id}/resend`,
      method: 'POST',
    })
  }

  static async accept(id: string) {
    return apiFetcher({
      url: `/invitations/${id}/accept`,
      method: 'POST',
    })
  }

  static async decline(id: string) {
    return apiFetcher({
      url: `/invitations/${id}/decline`,
      method: 'POST',
    })
  }

  static paramsFromListOptions(options?: ListOptions): URLSearchParams {
    const params: any = {}
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
