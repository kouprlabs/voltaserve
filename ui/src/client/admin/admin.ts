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
  updateTime: Date
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

export interface OrganizationExtendedManagement extends OrganizationManagement {
  workspaces: WorkspaceManagement[]
  groups: GroupManagement[]
}

export interface OrganizationManagementList extends ListResponse {
  data: OrganizationManagement[]
}

export interface GroupManagement extends CommonFields {
  name: string
  organization: OrganizationManagement
  OrganizationName: string
}

export interface GroupManagementList extends ListResponse {
  data: GroupManagement[]
}

export interface WorkspaceManagement extends CommonFields {
  name: string
  organization: OrganizationManagement
  storageCapacity: number
  rootId: string
  bucket: string
}

export interface WorkspaceManagementList extends ListResponse {
  data: WorkspaceManagement[]
}

export interface InvitationsManagement extends CommonFields {
  organization: OrganizationManagement
  ownerId: string
  email: string
  status: string
}

export interface InvitationsManagementList extends ListResponse {
  data: InvitationsManagement[]
}

export type ListOptions = {
  size?: number
  page?: number
}

type ListQueryParams = {
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

  static async listGroups(options: ListOptions) {
    return adminFetcher({
      url: `/group/all?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<GroupManagementList>
  }

  static async listOrganizations(options: ListOptions) {
    return adminFetcher({
      url: `/organization/all?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<OrganizationManagementList>
  }

  static async listWorkspaces(options: ListOptions) {
    return adminFetcher({
      url: `/workspace/all?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<WorkspaceManagementList>
  }

  static async listInvitations(options: ListOptions) {
    return adminFetcher({
      url: `/invitation/all?${this.paramsFromListOptions(options)}`,
      method: 'GET',
    }) as Promise<InvitationsManagementList>
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
    if (options?.page) {
      params.page = options.page.toString()
    }
    if (options?.size) {
      params.size = options.size.toString()
    }
    return new URLSearchParams(params)
  }
}
