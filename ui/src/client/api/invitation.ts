/* eslint-disable react-hooks/rules-of-hooks */
import useSWR from 'swr'
import { apiFetch, apiFetcher } from '@/client/fetch'
import { User } from '@/client/idp/user'
import { getAccessTokenOrRedirect } from '@/infra/token'
import { Organization } from './organization'

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

export type CreateOptions = {
  organizationId: string
  emails: string[]
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

  static useGetIncoming(swrOptions?: any) {
    return useSWR<Invitation[]>(
      `/invitations/get_incoming`,
      apiFetcher,
      swrOptions
    )
  }

  static useGetOutgoing(organizationId: string, swrOptions?: any) {
    return useSWR<Invitation[]>(
      organizationId
        ? `/invitations/get_outgoing?${new URLSearchParams({
            organization_id: organizationId,
          })}`
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
