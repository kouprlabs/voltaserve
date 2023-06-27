/* eslint-disable react-hooks/rules-of-hooks */
import useSWR from 'swr'
import { apiFetch, apiFetcher } from '@/client/fetch'
import { User } from '@/client/idp/user'
import { getAccessTokenOrRedirect } from '@/infra/token'
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
  size?: number
  page?: number
  sortBy?: SortBy
  sortOrder?: SortOrder
}

export default class InvitationAPI {
  static async create(options: CreateOptions) {
    return apiFetch(`/invitations`, {
      method: 'POST',
      body: JSON.stringify(options),
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    })
  }

  static useGetIncoming(options?: ListOptions, swrOptions?: any) {
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
    return useSWR<List>(
      `/invitations/get_incoming?${new URLSearchParams(params)}`,
      apiFetcher,
      swrOptions
    )
  }

  static useGetOutgoing(
    organizationId: string,
    options?: ListOptions,
    swrOptions?: any
  ) {
    const params: any = { org: organizationId }
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
    return useSWR<List>(
      organizationId
        ? `/invitations/get_outgoing?${new URLSearchParams(params)}`
        : null,
      apiFetcher,
      swrOptions
    )
  }

  static async delete(id: string) {
    return apiFetch(`/invitations/${id}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    })
  }

  static async resend(id: string) {
    return apiFetch(`/invitations/${id}/resend`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    })
  }

  static async accept(id: string) {
    return apiFetch(`/invitations/${id}/accept`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    })
  }

  static async decline(id: string) {
    return apiFetch(`/invitations/${id}/decline`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    })
  }
}
