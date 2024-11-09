// Copyright 2024 Mateusz Kaźmierczak.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import useSWR, { SWRConfiguration } from 'swr'
import { PermissionType } from '@/client/api/permission'
import { consoleFetcher } from '@/client/fetcher'
import { getConfig } from '@/config/config'
import { getAccessToken } from '@/infra/token'

export type ConsoleObject =
  | 'workspace'
  | 'organization'
  | 'group'
  | 'user'
  | 'invitation'
  | 'index'

export interface ListResponse<T> {
  totalElements: number
  totalPages: number | null
  page: number
  size: number
  data: T[]
}

export interface CommonFields {
  id: string
  createTime: Date
  updateTime: Date
}

export interface ConsoleOrganization extends CommonFields {
  name: string
}

export interface ConsoleOrganizationUser extends CommonFields {
  permission: PermissionType
  organizationId: string
  organizationName: string
  createTime: Date
}

export interface ConsoleGroup extends CommonFields {
  name: string
  organization: ConsoleOrganization
}

export interface ConsoleGroupUser extends CommonFields {
  permission: PermissionType
  groupId: string
  groupName: string
  createTime: Date
}

export interface ConsoleWorkspace extends CommonFields {
  name: string
  organization: ConsoleOrganization
  storageCapacity: number
  rootId: string
  bucket: string
}

export interface ConsoleWorkspaceUser extends CommonFields {
  permission: PermissionType
  workspaceId: string
  workspaceName: string
  createTime: Date
}

export interface ConsoleInvitation extends CommonFields {
  organization: ConsoleOrganization
  ownerId: string
  email: string
  status: string
}

export interface ConsoleUserOrganization extends CommonFields {
  username: string
  permission: PermissionType
  picture: string
}

export type ListOptions = {
  id?: string
  query?: string
  size?: number
  page?: number
}

type ListQueryParams = {
  id?: string
  query?: string
  page?: string
  size?: string
}

export interface BaseIDRequest {
  id?: string
}

export interface BaseNameRequest extends BaseIDRequest {
  name: string
}

export interface InvitationStatusRequest extends BaseIDRequest {
  accept: boolean
}

export interface CountResponse {
  count: number
}

export interface ComponentVersion {
  name: string
  currentVersion: string
  latestVersion: string
  updateAvailable: boolean
  location: string
}

export default class ConsoleAPI {
  static async checkIndexesAvailability() {
    const response = await fetch(`${getConfig().consoleURL}/index/all`, {
      method: 'GET',
      headers: {
        Authorization: `Bearer ${getAccessToken()}`,
      },
    })
    if (response) {
      return response.ok
    } else {
      return false
    }
  }

  static useListUsersByOrganization(
    options: ListOptions,
    swrOptions?: SWRConfiguration,
  ) {
    const url = `/organization/users?${this.paramsFromListOptions(options)}`
    return useSWR<ListResponse<ConsoleUserOrganization>>(
      url,
      () =>
        consoleFetcher({
          url,
          method: 'GET',
        }) as Promise<ListResponse<ConsoleUserOrganization>>,
      swrOptions,
    )
  }

  static listUsersByOrganization(options: ListOptions) {
    return consoleFetcher({
      url: `/organization/users?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<ListResponse<ConsoleUserOrganization>>
  }

  static useListGroupsByUser(
    options: ListOptions,
    swrOptions?: SWRConfiguration,
  ) {
    const url = `/user/groups?${this.paramsFromListOptions(options)}`
    return useSWR<ListResponse<ConsoleGroupUser>>(
      url,
      () =>
        consoleFetcher({
          url,
          method: 'GET',
        }) as Promise<ListResponse<ConsoleGroupUser>>,
      swrOptions,
    )
  }

  static listGroupsByUser(options: ListOptions) {
    return consoleFetcher({
      url: `/user/groups?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<ListResponse<ConsoleGroupUser>>
  }

  static useListGroupsByOrganization(
    options: ListOptions,
    swrOptions?: SWRConfiguration,
  ) {
    const url = `/organization/groups?${this.paramsFromListOptions(options)}`
    return useSWR<ListResponse<ConsoleGroup>>(
      url,
      () =>
        consoleFetcher({
          url,
          method: 'GET',
        }) as Promise<ListResponse<ConsoleGroup>>,
      swrOptions,
    )
  }

  static listGroupsByOrganization(options: ListOptions) {
    return consoleFetcher({
      url: `/organization/groups?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<ListResponse<ConsoleGroup>>
  }

  static useListOrganizationsByUser(
    options: ListOptions,
    swrOptions?: SWRConfiguration,
  ) {
    const url = `/user/organizations?${this.paramsFromListOptions(options)}`
    return useSWR<ListResponse<ConsoleOrganizationUser>>(
      url,
      () =>
        consoleFetcher({
          url,
          method: 'GET',
        }) as Promise<ListResponse<ConsoleOrganizationUser>>,
      swrOptions,
    )
  }

  static listOrganizationsByUser(options: ListOptions) {
    return consoleFetcher({
      url: `/user/organizations?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<ListResponse<ConsoleOrganizationUser>>
  }

  static useListWorkspacesByUser(
    options: ListOptions,
    swrOptions?: SWRConfiguration,
  ) {
    const url = `/user/workspaces?${this.paramsFromListOptions(options)}`
    return useSWR<ListResponse<ConsoleWorkspaceUser>>(
      url,
      () =>
        consoleFetcher({
          url,
          method: 'GET',
        }) as Promise<ListResponse<ConsoleWorkspaceUser>>,
      swrOptions,
    )
  }

  static listWorkspacesByUser(options: ListOptions) {
    return consoleFetcher({
      url: `/user/workspaces?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<ListResponse<ConsoleWorkspaceUser>>
  }

  static useListWorkspacesByOrganization(
    options: ListOptions,
    swrOptions?: SWRConfiguration,
  ) {
    const url = `/organization/workspaces?${this.paramsFromListOptions(options)}`
    return useSWR<ListResponse<ConsoleWorkspace>>(
      url,
      () =>
        consoleFetcher({
          url,
          method: 'GET',
        }) as Promise<ListResponse<ConsoleWorkspace>>,
      swrOptions,
    )
  }

  static listWorkspacesByOrganization(options: ListOptions) {
    return consoleFetcher({
      url: `/organization/workspaces?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<ListResponse<ConsoleWorkspace>>
  }

  static useGetOrganizationById(
    options: BaseIDRequest,
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

  static getOrganizationById(options: BaseIDRequest) {
    return consoleFetcher({
      url: `/organization?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<ConsoleOrganization>
  }

  static countObject(object: ConsoleObject) {
    return consoleFetcher({
      url: `/${object}/count`,
      method: 'GET',
    }) as Promise<CountResponse>
  }

  static getComponentsVersions(options: ListOptions) {
    return consoleFetcher({
      url: `/overview/version/internal?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<ComponentVersion>
  }

  static renameObject(options: BaseNameRequest, object: ConsoleObject) {
    return consoleFetcher({
      url: `/${object}`,
      method: 'PATCH',
      body: JSON.stringify(options),
    }) as Promise<void>
  }

  static invitationChangeStatus(options: InvitationStatusRequest) {
    return consoleFetcher({
      url: `/invitation`,
      method: 'PATCH',
      body: JSON.stringify(options),
    }) as Promise<void>
  }

  static useListObject<T>(
    object: ConsoleObject,
    options: ListOptions,
    swrOptions?: SWRConfiguration,
  ) {
    const url = `/${object}/all?${this.paramsFromListOptions(options)}`
    return useSWR<ListResponse<T>>(
      url,
      () => consoleFetcher({ url, method: 'GET' }) as Promise<ListResponse<T>>,
      swrOptions,
    )
  }

  static listObject<T>(object: ConsoleObject, options: ListOptions) {
    return consoleFetcher({
      url: `/${object}/all?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<ListResponse<T>>
  }

  static useSearchObject<T>(
    object: ConsoleObject,
    options: ListOptions,
    swrOptions?: SWRConfiguration,
  ) {
    const url = `/${object}/search?${this.paramsFromListOptions(options)}`
    return useSWR<ListResponse<T>>(
      url,
      () => consoleFetcher({ url, method: 'GET' }) as Promise<ListResponse<T>>,
      swrOptions,
    )
  }

  static searchObject<T>(object: ConsoleObject, options: ListOptions) {
    return consoleFetcher({
      url: `/${object}/search?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<ListResponse<T>>
  }

  static useListOrSearchObject<T>(
    object: ConsoleObject,
    options: ListOptions,
    swrOptions?: SWRConfiguration,
  ) {
    if (options.query) {
      return this.useSearchObject<T>(object, options, swrOptions)
    } else {
      return this.useListObject<T>(object, options, swrOptions)
    }
  }

  static paramsFromListOptions = (options: ListOptions): URLSearchParams => {
    const params: ListQueryParams = {}
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
