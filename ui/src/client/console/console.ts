// Copyright (c) 2024 Mateusz Kaźmierczak.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import useSWR, { SWRConfiguration } from 'swr'
import { PermissionType } from '@/client/api/permission'
import { consoleFetcher } from '@/client/fetcher'

export type ConsoleObject = 'workspace' | 'organization' | 'group' | 'user'

export interface ConsoleListResponse<T> {
  totalElements: number
  totalPages: number
  page: number
  size: number
  data: T[]
}

export interface ConsoleCommonFields {
  id: string
  createTime: Date
  updateTime: Date
}

export interface ConsoleOrganization extends ConsoleCommonFields {
  name: string
  permission?: string
}

export interface ConsoleOrganizationUser extends ConsoleCommonFields {
  permission: PermissionType
  organizationId: string
  organizationName: string
  createTime: Date
}

export interface ConsoleGroup extends ConsoleCommonFields {
  name: string
  organization: ConsoleOrganization
  permission?: string
}

export interface ConsoleGroupUser extends ConsoleCommonFields {
  permission: PermissionType
  groupId: string
  groupName: string
  createTime: Date
}

export interface ConsoleWorkspace extends ConsoleCommonFields {
  name: string
  organization: ConsoleOrganization
  storageCapacity: number
  rootId: string
  bucket: string
  permission?: string
}

export interface ConsoleWorkspaceUser extends ConsoleCommonFields {
  permission: PermissionType
  workspaceId: string
  workspaceName: string
  createTime: Date
}

export interface ConsoleInvitation extends ConsoleCommonFields {
  organization: ConsoleOrganization
  ownerId: string
  email: string
  status: string
}

export interface ConsoleUserOrganization extends ConsoleCommonFields {
  username: string
  permission: PermissionType
  picture: string
}

export type ConsoleListOptions = {
  id?: string
  query?: string
  size?: number
  page?: number
}

type ConsoleListQueryParams = {
  id?: string
  query?: string
  page?: string
  size?: string
}

export interface ConsoleBaseIDRequest {
  id?: string
}

export interface ConsoleBaseNameRequest extends ConsoleBaseIDRequest {
  name: string
}

export interface ConsoleInvitationStatusRequest extends ConsoleBaseIDRequest {
  accept: boolean
}

export type ConsoleGrantUserPermissionOptions = {
  userId: string
  resourceId: string
  resourceType: 'file' | 'group' | 'organization' | 'workspace'
  permission: string
}

export type ConsoleRevokeUserPermissionOptions = {
  userId: string
  resourceId: string
  resourceType: 'file' | 'group' | 'organization' | 'workspace'
}

export interface ConsoleCountResponse {
  count: number
}

export interface ConsoleComponentVersion {
  name: string
  currentVersion: string
  latestVersion: string
  updateAvailable: boolean
  location: string
}

export class ConsoleAPI {
  static useListUsersByOrganization(
    options: ConsoleListOptions,
    swrOptions?: SWRConfiguration,
  ) {
    const url = `/organization/users?${this.paramsFromListOptions(options)}`
    return useSWR<ConsoleListResponse<ConsoleUserOrganization>>(
      url,
      () =>
        consoleFetcher({
          url,
          method: 'GET',
        }) as Promise<ConsoleListResponse<ConsoleUserOrganization>>,
      swrOptions,
    )
  }

  static listUsersByOrganization(options: ConsoleListOptions) {
    return consoleFetcher({
      url: `/organization/users?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<ConsoleListResponse<ConsoleUserOrganization>>
  }

  static useListGroupsByUser(
    options: ConsoleListOptions,
    swrOptions?: SWRConfiguration,
  ) {
    const url = `/user/groups?${this.paramsFromListOptions(options)}`
    return useSWR<ConsoleListResponse<ConsoleGroupUser>>(
      url,
      () =>
        consoleFetcher({
          url,
          method: 'GET',
        }) as Promise<ConsoleListResponse<ConsoleGroupUser>>,
      swrOptions,
    )
  }

  static listGroupsByUser(options: ConsoleListOptions) {
    return consoleFetcher({
      url: `/user/groups?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<ConsoleListResponse<ConsoleGroupUser>>
  }

  static useListGroupsByOrganization(
    options: ConsoleListOptions,
    swrOptions?: SWRConfiguration,
  ) {
    const url = `/organization/groups?${this.paramsFromListOptions(options)}`
    return useSWR<ConsoleListResponse<ConsoleGroup>>(
      url,
      () =>
        consoleFetcher({
          url,
          method: 'GET',
        }) as Promise<ConsoleListResponse<ConsoleGroup>>,
      swrOptions,
    )
  }

  static listGroupsByOrganization(options: ConsoleListOptions) {
    return consoleFetcher({
      url: `/organization/groups?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<ConsoleListResponse<ConsoleGroup>>
  }

  static useListOrganizationsByUser(
    options: ConsoleListOptions,
    swrOptions?: SWRConfiguration,
  ) {
    const url = `/user/organizations?${this.paramsFromListOptions(options)}`
    return useSWR<ConsoleListResponse<ConsoleOrganizationUser>>(
      url,
      () =>
        consoleFetcher({
          url,
          method: 'GET',
        }) as Promise<ConsoleListResponse<ConsoleOrganizationUser>>,
      swrOptions,
    )
  }

  static listOrganizationsByUser(options: ConsoleListOptions) {
    return consoleFetcher({
      url: `/user/organizations?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<ConsoleListResponse<ConsoleOrganizationUser>>
  }

  static useListWorkspacesByUser(
    options: ConsoleListOptions,
    swrOptions?: SWRConfiguration,
  ) {
    const url = `/user/workspaces?${this.paramsFromListOptions(options)}`
    return useSWR<ConsoleListResponse<ConsoleWorkspaceUser>>(
      url,
      () =>
        consoleFetcher({
          url,
          method: 'GET',
        }) as Promise<ConsoleListResponse<ConsoleWorkspaceUser>>,
      swrOptions,
    )
  }

  static listWorkspacesByUser(options: ConsoleListOptions) {
    return consoleFetcher({
      url: `/user/workspaces?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<ConsoleListResponse<ConsoleWorkspaceUser>>
  }

  static useListWorkspacesByOrganization(
    options: ConsoleListOptions,
    swrOptions?: SWRConfiguration,
  ) {
    const url = `/organization/workspaces?${this.paramsFromListOptions(options)}`
    return useSWR<ConsoleListResponse<ConsoleWorkspace>>(
      url,
      () =>
        consoleFetcher({
          url,
          method: 'GET',
        }) as Promise<ConsoleListResponse<ConsoleWorkspace>>,
      swrOptions,
    )
  }

  static listWorkspacesByOrganization(options: ConsoleListOptions) {
    return consoleFetcher({
      url: `/organization/workspaces?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<ConsoleListResponse<ConsoleWorkspace>>
  }

  static useGetOrganizationById(
    options: ConsoleBaseIDRequest,
    swrOptions?: SWRConfiguration,
  ) {
    const url = `/organization?${this.paramsFromListOptions(options)}`
    return useSWR<ConsoleOrganization>(
      options.id ? url : null,
      () =>
        consoleFetcher({
          url,
          method: 'GET',
        }) as Promise<ConsoleOrganization>,
      swrOptions,
    )
  }

  static getOrganizationById(options: ConsoleBaseIDRequest) {
    return consoleFetcher({
      url: `/organization?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<ConsoleOrganization>
  }

  static countObject(object: ConsoleObject) {
    return consoleFetcher({
      url: `/${object}/count`,
      method: 'GET',
    }) as Promise<ConsoleCountResponse>
  }

  static getComponentsVersions(options: ConsoleListOptions) {
    return consoleFetcher({
      url: `/overview/version/internal?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<ConsoleComponentVersion>
  }

  static useListObject<T>(
    object: ConsoleObject,
    options: ConsoleListOptions,
    swrOptions?: SWRConfiguration,
  ) {
    const url = `/${object}/all?${this.paramsFromListOptions(options)}`
    return useSWR<ConsoleListResponse<T>>(
      url,
      () =>
        consoleFetcher({ url, method: 'GET' }) as Promise<
          ConsoleListResponse<T>
        >,
      swrOptions,
    )
  }

  static listObject<T>(object: ConsoleObject, options: ConsoleListOptions) {
    return consoleFetcher({
      url: `/${object}/all?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<ConsoleListResponse<T>>
  }

  static useSearchObject<T>(
    object: ConsoleObject,
    options: ConsoleListOptions,
    swrOptions?: SWRConfiguration,
  ) {
    const url = `/${object}/search?${this.paramsFromListOptions(options)}`
    return useSWR<ConsoleListResponse<T>>(
      url,
      () =>
        consoleFetcher({ url, method: 'GET' }) as Promise<
          ConsoleListResponse<T>
        >,
      swrOptions,
    )
  }

  static searchObject<T>(object: ConsoleObject, options: ConsoleListOptions) {
    return consoleFetcher({
      url: `/${object}/search?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<ConsoleListResponse<T>>
  }

  static useListOrSearchObject<T>(
    object: ConsoleObject,
    options: ConsoleListOptions,
    swrOptions?: SWRConfiguration,
  ) {
    if (options.query) {
      return this.useSearchObject<T>(object, options, swrOptions)
    } else {
      return this.useListObject<T>(object, options, swrOptions)
    }
  }

  static grantUserPermission(options: ConsoleGrantUserPermissionOptions) {
    return consoleFetcher({
      url: '/user_permission/grant',
      method: 'POST',
      body: JSON.stringify(options),
    }) as Promise<void>
  }

  static revokeUserPermission(options: ConsoleRevokeUserPermissionOptions) {
    return consoleFetcher({
      url: '/user_permission/revoke',
      method: 'POST',
      body: JSON.stringify(options),
    }) as Promise<void>
  }

  static paramsFromListOptions = (
    options: ConsoleListOptions,
  ): URLSearchParams => {
    const params: ConsoleListQueryParams = {}
    if (options.id) {
      params.id = options.id.toString()
    }
    if (options.page) {
      params.page = options.page.toString()
    }
    if (options.size) {
      params.size = options.size.toString()
    }
    if (options.query) {
      params.query = options.query.toString()
    }
    if (options.page) {
      params.page = options.page.toString()
    }
    if (options.size) {
      params.size = options.size.toString()
    }
    return new URLSearchParams(params)
  }
}
