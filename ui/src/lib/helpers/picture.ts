// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { Picture } from '@/client/types'
import { getAccessTokenOrRedirect } from '@/infra/token'

export function getPictureUrl(picture: Picture) {
  return `/proxy/idp/v3/users/me/picture${picture.extension}?${new URLSearchParams({
    access_token: getAccessTokenOrRedirect(),
  })}`
}

type PictureUrlByIdOptions = {
  organizationId?: string
  groupId?: string
  invitationId?: string
}

export function getPictureUrlById(id: string, picture: Picture, options?: PictureUrlByIdOptions) {
  return `/proxy/api/v3/users/${id}/picture${picture.extension}?${new URLSearchParams({
    access_token: getAccessTokenOrRedirect(),
    organization_id: options?.organizationId ?? '',
    group_id: options?.groupId ?? '',
    invitation_id: options?.invitationId ?? '',
  })}`
}
