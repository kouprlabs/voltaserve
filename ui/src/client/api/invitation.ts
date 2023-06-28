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
  organizationId?: string
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
    return useSWR<List>(
      `/invitations/get_incoming?${this.paramsFromListOptions(options)}`,
      apiFetcher,
      swrOptions
    )
  }

  static useGetOutgoing(options?: ListOptions, swrOptions?: any) {
    return useSWR<List>(
      options?.organizationId
        ? `/invitations/get_outgoing?${this.paramsFromListOptions(options)}`
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
