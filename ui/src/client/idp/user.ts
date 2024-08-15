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

export type User = {
  id: string
  username: string
  email: string
  fullName: string
  picture?: string
  pendingEmail?: string
}

export interface AdminUser extends User {
  isEmailConfirmed: boolean
  createTime: Date
  updateTime: Date
  isAdmin: boolean
  isActive: boolean
}

export interface AdminUsersResponse {
  data: AdminUser[]
  totalElements: number
  page: number
  size: number
}

export type UpdateFullNameOptions = {
  fullName: string
}

export interface baseUserIdRequest {
  id: string
}

export type UpdateEmailRequestOptions = {
  email: string
}

export type UpdateEmailConfirmationOptions = {
  token: string
}

export interface suspendUserOptions extends baseUserIdRequest {
  suspend: boolean
}

export interface makeAdminOptions extends baseUserIdRequest {
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
  id?: string
  size?: number
  page?: number
}

type ListQueryParams = {
  id?: string
  page?: string
  size?: string
}

export default class UserAPI {
  static useGet(swrOptions?: SWRConfiguration) {
    const url = `/user`
    return useSWR<User>(
      url,
      () => idpFetcher({ url, method: 'GET' }) as Promise<User>,
      swrOptions,
    )
  }

  static async getUserById(options: baseUserIdRequest) {
    return idpFetcher({
      url: `/user/by_id?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<AdminUser>
  }

  static async getAllUsers(options: ListOptions) {
    return idpFetcher({
      url: `/user/all?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<AdminUsersResponse>
  }

  static async updateFullName(options: UpdateFullNameOptions) {
    return idpFetcher({
      url: `/user/update_full_name`,
      method: 'POST',
      body: JSON.stringify(options),
    }) as Promise<User>
  }

  static async suspendUser(options: suspendUserOptions) {
    return idpFetcher({
      url: `/user/suspend`,
      method: 'PATCH',
      body: JSON.stringify(options),
    }) as Promise<User>
  }

  static async makeAdmin(options: makeAdminOptions) {
    return idpFetcher({
      url: `/user/admin`,
      method: 'PATCH',
      body: JSON.stringify(options),
    }) as Promise<User>
  }

  static async updateEmailRequest(options: UpdateEmailRequestOptions) {
    return idpFetcher({
      url: `/user/update_email_request`,
      method: 'POST',
      body: JSON.stringify(options),
    }) as Promise<User>
  }

  static async updateEmailConfirmation(
    options: UpdateEmailConfirmationOptions,
  ) {
    return idpFetcher({
      url: `/user/update_email_confirmation`,
      method: 'POST',
      body: JSON.stringify(options),
    }) as Promise<User>
  }

  static async updatePassword(options: UpdatePasswordOptions) {
    return idpFetcher({
      url: `/user/update_password`,
      method: 'POST',
      body: JSON.stringify(options),
    }) as Promise<User>
  }

  static async delete(options: DeleteOptions) {
    return idpFetcher({
      url: `/user`,
      method: 'DELETE',
      body: JSON.stringify(options),
    })
  }

  static async updatePicture(file: File) {
    const body = new FormData()
    body.append('file', file)
    return idpFetcher({
      url: `/user/update_picture`,
      method: 'POST',
      body,
      contentType: 'multipart/form-data',
    }) as Promise<User>
  }

  static async deletePicture() {
    return idpFetcher({
      url: `/user/delete_picture`,
      method: 'POST',
    }) as Promise<User>
  }

  static paramsFromListOptions(options?: ListOptions): URLSearchParams {
    const params: ListQueryParams = {}
    if (options?.id) {
      params.id = options.id.toString()
    }
    if (options?.page) {
      params.page = options.page.toString()
    }
    if (options?.size) {
      params.size = options.size.toString()
    }
    return new URLSearchParams(params)
  }
}
