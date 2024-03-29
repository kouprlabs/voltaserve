export const FILES_PAGINATION_STEP = 21

export function filePaginationSteps() {
  return [
    FILES_PAGINATION_STEP,
    FILES_PAGINATION_STEP * 2,
    FILES_PAGINATION_STEP * 4,
    FILES_PAGINATION_STEP * 5,
  ]
}

export function filesPaginationStorage() {
  return {
    prefix: 'voltaserve',
    namespace: 'files',
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
