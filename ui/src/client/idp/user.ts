// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import useSWR, { SWRConfiguration } from 'swr'
import { idpFetcher } from '@/client/fetcher'
import { Picture } from '@/client/types'

export type AuthUser = {
  id: string
  username: string
  email: string
  fullName: string
  picture?: Picture
  pendingEmail?: string
}

export interface ConsoleUser extends AuthUser {
  isEmailConfirmed: boolean
  createTime: Date
  updateTime: Date
  isAdmin: boolean
  isActive: boolean
}

export interface ConsoleUserList {
  data: ConsoleUser[]
  totalElements: number
  totalPages: number
  page: number
  size: number
}

export type AuthUserUpdateFullNameOptions = {
  fullName: string
}

export type AuthUserUpdateEmailRequestOptions = {
  email: string
}

export type AuthUserUpdateEmailConfirmationOptions = {
  token: string
}

export interface AuthUserSuspendOptions {
  suspend: boolean
}

export interface AuthUserMakeAdminOptions {
  makeAdmin: boolean
}

export type AuthUserUpdatePasswordOptions = {
  currentPassword: string
  newPassword: string
}

type AuthUserListOptions = {
  query?: string
  id?: string
  size?: number
  page?: number
}

type AuthUserListQueryParams = {
  query?: string
  id?: string
  page?: string
  size?: string
}

export class AuthUserAPI {
  static useGet(swrOptions?: SWRConfiguration) {
    const url = `/users/me`
    return useSWR<AuthUser>(
      url,
      () => idpFetcher({ url, method: 'GET' }) as Promise<AuthUser>,
      swrOptions,
    )
  }

  static useGetById(id?: string, swrOptions?: SWRConfiguration) {
    const url = `/users/${id}`
    return useSWR<ConsoleUser>(
      id ? url : null,
      () =>
        idpFetcher({
          url,
          method: 'GET',
        }) as Promise<ConsoleUser>,
      swrOptions,
    )
  }

  static getById(id: string) {
    return idpFetcher({
      url: `/users/${id}`,
      method: 'GET',
    }) as Promise<ConsoleUser>
  }

  static useList(options: AuthUserListOptions, swrOptions?: SWRConfiguration) {
    const url = `/users?${this.paramsFromListOptions(options)}`
    return useSWR<ConsoleUserList>(
      url,
      () => idpFetcher({ url, method: 'GET' }) as Promise<ConsoleUserList>,
      swrOptions,
    )
  }

  static list(options: AuthUserListOptions) {
    return idpFetcher({
      url: `/users?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<ConsoleUserList>
  }

  static updateFullName(options: AuthUserUpdateFullNameOptions) {
    return idpFetcher({
      url: `/users/me/update_full_name`,
      method: 'POST',
      body: JSON.stringify(options),
    }) as Promise<AuthUser>
  }

  static suspend(id: string, options: AuthUserSuspendOptions) {
    return idpFetcher({
      url: `/users/${id}/suspend`,
      method: 'POST',
      body: JSON.stringify(options),
    }) as Promise<AuthUser>
  }

  static makeAdmin(id: string, options: AuthUserMakeAdminOptions) {
    return idpFetcher({
      url: `/users/${id}/make_admin`,
      method: 'POST',
      body: JSON.stringify(options),
    }) as Promise<AuthUser>
  }

  static updateEmailRequest(options: AuthUserUpdateEmailRequestOptions) {
    return idpFetcher({
      url: `/users/me/update_email_request`,
      method: 'POST',
      body: JSON.stringify(options),
    }) as Promise<AuthUser>
  }

  static updateEmailConfirmation(
    options: AuthUserUpdateEmailConfirmationOptions,
  ) {
    return idpFetcher({
      url: `/users/me/update_email_confirmation`,
      method: 'POST',
      body: JSON.stringify(options),
    }) as Promise<AuthUser>
  }

  static updatePassword(options: AuthUserUpdatePasswordOptions) {
    return idpFetcher({
      url: `/users/me/update_password`,
      method: 'POST',
      body: JSON.stringify(options),
    }) as Promise<AuthUser>
  }

  static async delete() {
    return idpFetcher({
      url: `/users/me`,
      method: 'DELETE',
    })
  }

  static updatePicture(file: File) {
    const body = new FormData()
    body.append('file', file)
    return idpFetcher({
      url: `/users/me/update_picture`,
      method: 'POST',
      body,
      contentType: 'multipart/form-data',
    }) as Promise<AuthUser>
  }

  static deletePicture() {
    return idpFetcher({
      url: `/users/me/delete_picture`,
      method: 'POST',
    }) as Promise<AuthUser>
  }

  static paramsFromListOptions(options: AuthUserListOptions): URLSearchParams {
    const params: AuthUserListQueryParams = {}
    if (options?.id) {
      params.id = options.id.toString()
    }
    if (options.query) {
      params.query = options.query.toString()
    }
    if (options?.page) {
      params.page = options.page.toString()
    }
    if (options?.size) {
      params.size = options.size.toString()
    }
    if (options?.query) {
      params.query = options.query.toString()
    }
    return new URLSearchParams(params)
  }
}
