/* eslint-disable react-hooks/rules-of-hooks */
import useSWR from 'swr'
import { idpFetcher } from '@/client/fetcher'

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
    const url = `/user`
    return useSWR<User>(
      url,
      () => idpFetcher({ url, method: 'GET' }) as Promise<User>,
      swrOptions,
    )
  }

  static async updateFullName(options: UpdateFullNameOptions) {
    return idpFetcher({
      url: `/user/update_full_name`,
      method: 'POST',
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
}
