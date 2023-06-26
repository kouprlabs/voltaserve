/* eslint-disable react-hooks/rules-of-hooks */
import useSWR from 'swr'
import { idpFetch, idpFetcher } from '@/client/fetch'
import { getAccessTokenOrRedirect } from '@/infra/token'

export type User = {
  id: string
  username: string
  email: string
  fullName: string
  picture?: string
  pendingEmail?: string
}

export type UpdateFullNameOptions = {
  fullName: string
}

export type UpdateEmailRequestOptions = {
  email: string
}

export type UpdateEmailConfirmationOptions = {
  token: string
}

export type UpdatePasswordOptions = {
  currentPassword: string
  newPassword: string
}

export type DeleteOptions = {
  password: string
}

export default class UserAPI {
  static useGet(swrOptions?: any) {
    return useSWR<User>(`/user`, idpFetcher, swrOptions)
  }

  static async updateFullName(options: UpdateFullNameOptions): Promise<User> {
    return idpFetch(`/user/update_full_name`, {
      method: 'POST',
      body: JSON.stringify(options),
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    }).then((result) => result.json())
  }

  static async updateEmailRequest(
    options: UpdateEmailRequestOptions
  ): Promise<User> {
    return idpFetch(`/user/update_email_request`, {
      method: 'POST',
      body: JSON.stringify(options),
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    }).then((result) => result.json())
  }

  static async updateEmailConfirmation(
    options: UpdateEmailConfirmationOptions
  ): Promise<User> {
    return idpFetch(`/user/update_email_confirmation`, {
      method: 'POST',
      body: JSON.stringify(options),
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    }).then((result) => result.json())
  }

  static async updatePassword(options: UpdatePasswordOptions): Promise<User> {
    return idpFetch(`/user/update_password`, {
      method: 'POST',
      body: JSON.stringify(options),
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    }).then((result) => result.json())
  }

  static async delete(options: DeleteOptions) {
    return idpFetch(`/user`, {
      method: 'DELETE',
      body: JSON.stringify(options),
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    })
  }

  static async updatePicture(file: File): Promise<User> {
    const body = new FormData()
    body.append('file', file)
    return idpFetch(`/user/update_picture`, {
      method: 'POST',
      body,
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
      },
    }).then((result) => result.json())
  }

  static async deletePicture(): Promise<User> {
    return idpFetch(`/user/delete_picture`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
      },
    }).then((result) => result.json())
  }
}
