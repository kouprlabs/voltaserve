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
import { adminFetcher } from '@/client/fetcher'
import { getConfig } from '@/config/config'

export interface ListResponse {
  totalElements: number
  totalPages: number | null
  page: number
  size: number
}

export interface CommonFields {
  id: string
  createTime: Date
  updateTime?: Date
}

export type IndexManagement = {
  tablename: string
  indexname: string
  indexdef: string
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
  size?: number
  page?: number
}

type ListQueryParams = {
  id?: string
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

export default class AdminApi {
  static async checkIndexesAvailability() {
    const response = await fetch(`${getConfig().adminURL}/index/all`, {
      method: 'GET',
      headers: {
        'Access-Control-Allow-Origin': `${getConfig().adminURL}`, // TODO: To be deleted after local tests
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
        adminFetcher({
          url,
          method: 'GET',
        }) as Promise<IndexManagementList>,
      swrOptions,
    )
  }

  static async getUsersByOrganization(options: ListOptions) {
    return adminFetcher({
      url: `/organization/users?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<UserOrganizationManagementList>
  }

  static async listGroups(options: ListOptions) {
    return adminFetcher({
      url: `/group/all?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<GroupManagementList>
  }

  static async getGroupsByUser(options: ListOptions) {
    return adminFetcher({
      url: `/user/groups?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<GroupUserManagementList>
  }

  static async getGroupsByOrganization(options: ListOptions) {
    return adminFetcher({
      url: `/organization/groups?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<GroupManagementList>
  }

  static async listOrganizations(options: ListOptions) {
    return adminFetcher({
      url: `/organization/all?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<OrganizationManagementList>
  }

  static async getOrganizationsByUser(options: ListOptions) {
    return adminFetcher({
      url: `/user/organizations?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<OrganizationUserManagementList>
  }

  static async getOrganizationById(options: baseIdRequest) {
    return adminFetcher({
      url: `/organization?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<OrganizationManagement>
  }

  static async listWorkspaces(options: ListOptions) {
    return adminFetcher({
      url: `/workspace/all?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<WorkspaceManagementList>
  }

  static async getWorkspacesByUser(options: ListOptions) {
    return adminFetcher({
      url: `/user/workspaces?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<WorkspaceUserManagementList>
  }

  static async getWorkspacesByOrganization(options: ListOptions) {
    return adminFetcher({
      url: `/organization/workspaces?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<WorkspaceManagementList>
  }

  static async listInvitations(options: ListOptions) {
    return adminFetcher({
      url: `/invitation/all?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<InvitationManagementList>
  }

  static async renameObject(options: baseIdNameRequest, object: string) {
    return adminFetcher({
      url: `/${object}`,
      method: 'PATCH',
      body: JSON.stringify(options),
    }) as Promise<void>
  }

  static async invitationChangeStatus(options: invitationStatusRequest) {
    return adminFetcher({
      url: `/invitation`,
      method: 'PATCH',
      body: JSON.stringify(options),
    }) as Promise<void>
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
