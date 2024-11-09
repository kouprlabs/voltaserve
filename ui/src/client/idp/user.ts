// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import useSWR, { SWRConfiguration } from 'swr'
import { idpFetcher } from '@/client/fetcher'
import { Picture } from '@/client/types'

export type User = {
  id: string
  username: string
  email: string
  fullName: string
  picture?: Picture
  pendingEmail?: string
}

export interface ConsoleUser extends User {
  isEmailConfirmed: boolean
  createTime: Date
  updateTime: Date
  isAdmin: boolean
  isActive: boolean
}

export interface List {
  data: ConsoleUser[]
  totalElements: number
  page: number
  size: number
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

export interface SuspendOptions {
  suspend: boolean
}

export interface MakeAdminOptions {
  makeAdmin: boolean
}

export type UpdatePasswordOptions = {
  currentPassword: string
  newPassword: string
}

export type DeleteOptions = {
  password: string
}

type ListOptions = {
  query?: string
  id?: string
  size?: number
  page?: number
}

type ListQueryParams = {
  query?: string
  id?: string
  page?: string
  size?: string
}

export default class UserAPI {
  static useGet(swrOptions?: SWRConfiguration) {
    const url = `/users/me`
    return useSWR<User>(
      url,
      () => idpFetcher({ url, method: 'GET' }) as Promise<User>,
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

  static useList(options: ListOptions, swrOptions?: SWRConfiguration) {
    const url = `/users?${this.paramsFromListOptions(options)}`
    return useSWR<List>(
      url,
      () => idpFetcher({ url, method: 'GET' }) as Promise<List>,
      swrOptions,
    )
  }

  static list(options: ListOptions) {
    return idpFetcher({
      url: `/users?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<List>
  }

  static updateFullName(options: UpdateFullNameOptions) {
    return idpFetcher({
      url: `/users/me/update_full_name`,
      method: 'POST',
      body: JSON.stringify(options),
    }) as Promise<User>
  }

  static suspend(id: string, options: SuspendOptions) {
    return idpFetcher({
      url: `/users/${id}/suspend`,
      method: 'POST',
      body: JSON.stringify(options),
    }) as Promise<User>
  }

  static makeAdmin(id: string, options: MakeAdminOptions) {
    return idpFetcher({
      url: `/users/${id}/make_admin`,
      method: 'POST',
      body: JSON.stringify(options),
    }) as Promise<User>
  }

  static updateEmailRequest(options: UpdateEmailRequestOptions) {
    return idpFetcher({
      url: `/users/me/update_email_request`,
      method: 'POST',
      body: JSON.stringify(options),
    }) as Promise<User>
  }

  static updateEmailConfirmation(options: UpdateEmailConfirmationOptions) {
    return idpFetcher({
      url: `/users/me/update_email_confirmation`,
      method: 'POST',
      body: JSON.stringify(options),
    }) as Promise<User>
  }

  static updatePassword(options: UpdatePasswordOptions) {
    return idpFetcher({
      url: `/users/me/update_password`,
      method: 'POST',
      body: JSON.stringify(options),
    }) as Promise<User>
  }

  static async delete(options: DeleteOptions) {
    return idpFetcher({
      url: `/users/me`,
      method: 'DELETE',
      body: JSON.stringify(options),
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
    }) as Promise<User>
  }

  static deletePicture() {
    return idpFetcher({
      url: `/users/me/delete_picture`,
      method: 'POST',
    }) as Promise<User>
  }

  static paramsFromListOptions(options: ListOptions): URLSearchParams {
    const params: ListQueryParams = {}
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
