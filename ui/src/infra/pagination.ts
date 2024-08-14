// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { StorageOptions } from '@/lib/hooks/page-pagination'

export const FILES_PAGINATION_STEP = 21

export function filePaginationSteps() {
  return [
    FILES_PAGINATION_STEP,
    FILES_PAGINATION_STEP * 2,
    FILES_PAGINATION_STEP * 4,
    FILES_PAGINATION_STEP * 5,
  ]
}

export function filesPaginationStorage(): StorageOptions {
  return {
    prefix: 'voltaserve',
    namespace: 'files',
    enabled: true,
  }
}

export function incomingInvitationPaginationStorage() {
  return {
    prefix: 'voltaserve',
    namespace: 'incoming_invitation',
  }
}

export function groupPaginationStorage() {
  return {
    prefix: 'voltaserve',
    namespace: 'group',
  }
}

export function groupMemberPaginationStorage() {
  return {
    prefix: 'voltaserve',
    namespace: 'group_member',
  }
}

export function outgoingInvitationPaginationStorage() {
  return {
    prefix: 'voltaserve',
    namespace: 'outgoing_invitation',
  }
}

export function organizationPaginationStorage() {
  return {
    prefix: 'voltaserve',
    namespace: 'organization',
  }
}

export function taskPaginationStorage() {
  return {
    prefix: 'voltaserve',
    namespace: 'task',
  }
}

export function organizationMemberPaginationStorage() {
  return {
    prefix: 'voltaserve',
    namespace: 'organization_member',
  }
}

export function workspacePaginationStorage() {
  return {
    prefix: 'voltaserve',
    namespace: 'workspace',
  }
}

export function adminUsersPaginationStorage() {
  return {
    prefix: 'voltaserve',
    namespace: 'admin_users',
  }
}

export function adminGroupsPaginationStorage() {
  return {
    prefix: 'voltaserve',
    namespace: 'admin_groups',
  }
}

export function adminInvitationsPaginationStorage() {
  return {
    prefix: 'voltaserve',
    namespace: 'admin_invitations',
  }
}

export function adminOrganizationsPaginationStorage() {
  return {
    prefix: 'voltaserve',
    namespace: 'admin_organizations',
  }
}

export function adminWorkspacesPaginationStorage() {
  return {
    prefix: 'voltaserve',
    namespace: 'admin_workspaces',
  }
}
