/* eslint-disable react-hooks/rules-of-hooks */
import useSWR from 'swr'
import { getAccessTokenOrRedirect } from '@/infra/token'
import { idpFetch, idpFetcher } from './fetch'

export type User = {
  id: string
  username: string
  email: string
  fullName: string
  picture?: string
}

export type UserUpdateFullNameOptions = {
  fullName: string
}

export type UserUpdateEmailOptions = {
  email: string
}

export type UserUpdatePasswordOptions = {
  currentPassword: string
  newPassword: string
}

export type UserDeleteOptions = {
  password: string
}

export default class UserAPI {
  static useGet(swrOptions?: any) {
    return useSWR<User>(`/user`, idpFetcher, swrOptions)
  }

  static async updateFullName(
    options: UserUpdateFullNameOptions
  ): Promise<User> {
    return idpFetch(`/user/update_full_name`, {
      method: 'POST',
      body: JSON.stringify(options),
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    }).then((result) => result.json())
  }

  static async updateEmail(options: UserUpdateEmailOptions): Promise<User> {
    return idpFetch(`/user/update_email`, {
      method: 'POST',
      body: JSON.stringify(options),
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    }).then((result) => result.json())
  }

  static async updatePassword(
    options: UserUpdatePasswordOptions
  ): Promise<User> {
    return idpFetch(`/user/update_password`, {
      method: 'POST',
      body: JSON.stringify(options),
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
    }).then((result) => result.json())
  }

  static async delete(options: UserDeleteOptions) {
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
