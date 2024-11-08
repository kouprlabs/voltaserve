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

export interface ListResponse {
  totalElements: number
  totalPages: number | null
  page: number
  size: number
}

export interface EmptyListResponse extends ListResponse {
  data: []
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

export interface IndexManagementList extends ListResponse {
  data: IndexManagement[]
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

export interface OrganizationManagementList extends ListResponse {
  data: OrganizationManagement[]
}

export interface OrganizationUserManagementList extends ListResponse {
  data: OrganizationUserManagement[]
}

export interface GroupManagement extends CommonFields {
  name: string
  organization: OrganizationManagement
  OrganizationName: string
}

export interface GroupUserManagement extends CommonFields {
  permission: PermissionType
  groupId: string
  groupName: string
  createTime: Date
}

export interface GroupManagementList extends ListResponse {
  data: GroupManagement[]
}

export interface GroupUserManagementList extends ListResponse {
  data: GroupUserManagement[]
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

export interface WorkspaceManagementList extends ListResponse {
  data: WorkspaceManagement[]
}

export interface WorkspaceUserManagementList extends ListResponse {
  data: WorkspaceUserManagement[]
}

export interface InvitationManagement extends CommonFields {
  organization: OrganizationManagement
  ownerId: string
  email: string
  status: string
}

export interface InvitationManagementList extends ListResponse {
  data: InvitationManagement[]
}

export interface UserOrganizationManagement extends CommonFields {
  username: string
  permission: PermissionType
  picture: string
}

export interface UserOrganizationManagementList extends ListResponse {
  data: UserOrganizationManagement[]
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
    return useSWR<IndexManagementList>(
      url,
      () =>
        consoleFetcher({
          url,
          method: 'GET',
        }) as Promise<IndexManagementList>,
      swrOptions,
    )
  }

  static getUsersByOrganization(options: ListOptions) {
    return consoleFetcher({
      url: `/organization/users?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<UserOrganizationManagementList>
  }

  static useListGroups(options: ListOptions, swrOptions?: SWRConfiguration) {
    const url = `/workspaces?${this.paramsFromListOptions(options)}`
    return useSWR<GroupManagementList>(
      url,
      () =>
        consoleFetcher({ url, method: 'GET' }) as Promise<GroupManagementList>,
      swrOptions,
    )
  }

  static listGroups(options: ListOptions) {
    return consoleFetcher({
      url: `/group/all?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<GroupManagementList>
  }

  static getGroupsByUser(options: ListOptions) {
    return consoleFetcher({
      url: `/user/groups?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<GroupUserManagementList>
  }

  static getGroupsByOrganization(options: ListOptions) {
    return consoleFetcher({
      url: `/organization/groups?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<GroupManagementList>
  }

  static listOrganizations(options: ListOptions) {
    return consoleFetcher({
      url: `/organization/all?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<OrganizationManagementList>
  }

  static getOrganizationsByUser(options: ListOptions) {
    return consoleFetcher({
      url: `/user/organizations?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<OrganizationUserManagementList>
  }

  static getOrganizationById(options: baseIdRequest) {
    return consoleFetcher({
      url: `/organization?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<OrganizationManagement>
  }

  static listWorkspaces(options: ListOptions) {
    return consoleFetcher({
      url: `/workspace/all?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<WorkspaceManagementList>
  }

  static getWorkspacesByUser(options: ListOptions) {
    return consoleFetcher({
      url: `/user/workspaces?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<WorkspaceUserManagementList>
  }

  static getWorkspacesByOrganization(options: ListOptions) {
    return consoleFetcher({
      url: `/organization/workspaces?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<WorkspaceManagementList>
  }

  static useListInvitations(
    options: ListOptions,
    swrOptions?: SWRConfiguration,
  ) {
    const url = `/invitation/all?${this.paramsFromListOptions(options)}`
    return useSWR<InvitationManagementList>(
      url,
      () =>
        consoleFetcher({
          url,
          method: 'GET',
        }) as Promise<InvitationManagementList>,
      swrOptions,
    )
  }

  static listInvitations(options: ListOptions) {
    return consoleFetcher({
      url: `/invitation/all?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<InvitationManagementList>
  }

  static countObject(object: string) {
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

  static renameObject(options: baseIdNameRequest, object: string) {
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

  static listObject(object: string, options: ListOptions) {
    return consoleFetcher({
      url: `/${object}/all?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<EmptyListResponse>
  }

  static searchObject(object: string, options: ListOptions) {
    return consoleFetcher({
      url: `/${object}/search?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<EmptyListResponse>
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
