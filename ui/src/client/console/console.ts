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

export type IndexManagement = {
  tableName: string
  indexName: string
  indexDef: string
}

export interface OrganizationManagement extends CommonFields {
  name: string
}

export interface OrganizationUserManagement extends CommonFields {
  permission: PermissionType
  organizationId: string
  organizationName: string
  createTime: Date
}

export interface GroupManagement extends CommonFields {
  name: string
  organization: OrganizationManagement
}

export interface GroupUserManagement extends CommonFields {
  permission: PermissionType
  groupId: string
  groupName: string
  createTime: Date
}

export interface WorkspaceManagement extends CommonFields {
  name: string
  organization: OrganizationManagement
  storageCapacity: number
  rootId: string
  bucket: string
}

export interface WorkspaceUserManagement extends CommonFields {
  permission: PermissionType
  workspaceId: string
  workspaceName: string
  createTime: Date
}

export interface InvitationManagement extends CommonFields {
  organization: OrganizationManagement
  ownerId: string
  email: string
  status: string
}

export interface UserOrganizationManagement extends CommonFields {
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

export interface baseIdRequest {
  id: string
}

export interface baseIdNameRequest extends baseIdRequest {
  name: string
}

export interface invitationStatusRequest extends baseIdRequest {
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

  static useListIndexes(options: ListOptions, swrOptions?: SWRConfiguration) {
    const url = `/index/all?${this.paramsFromListOptions(options)}`
    return useSWR<ListResponse<IndexManagement>>(
      url,
      () =>
        consoleFetcher({
          url,
          method: 'GET',
        }) as Promise<ListResponse<IndexManagement>>,
      swrOptions,
    )
  }

  static getUsersByOrganization(options: ListOptions) {
    return consoleFetcher({
      url: `/organization/users?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<ListResponse<UserOrganizationManagement>>
  }

  static getGroupsByUser(options: ListOptions) {
    return consoleFetcher({
      url: `/user/groups?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<ListResponse<GroupUserManagement>>
  }

  static getGroupsByOrganization(options: ListOptions) {
    return consoleFetcher({
      url: `/organization/groups?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<ListResponse<GroupManagement>>
  }

  static getOrganizationsByUser(options: ListOptions) {
    return consoleFetcher({
      url: `/user/organizations?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<ListResponse<OrganizationUserManagement>>
  }

  static getOrganizationById(options: baseIdRequest) {
    return consoleFetcher({
      url: `/organization?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<OrganizationManagement>
  }

  static getWorkspacesByUser(options: ListOptions) {
    return consoleFetcher({
      url: `/user/workspaces?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<ListResponse<WorkspaceUserManagement>>
  }

  static getWorkspacesByOrganization(options: ListOptions) {
    return consoleFetcher({
      url: `/organization/workspaces?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<ListResponse<WorkspaceManagement>>
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

  static renameObject(options: baseIdNameRequest, object: ConsoleObject) {
    return consoleFetcher({
      url: `/${object}`,
      method: 'PATCH',
      body: JSON.stringify(options),
    }) as Promise<void>
  }

  static invitationChangeStatus(options: invitationStatusRequest) {
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
