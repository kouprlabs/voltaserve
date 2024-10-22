import { Picture } from '@/client/types'
import { getAccessTokenOrRedirect } from '@/infra/token'

export function getPictureUrl(picture: Picture) {
  return `/proxy/idp/v2/user/picture${picture.extension}?${new URLSearchParams({
    access_token: getAccessTokenOrRedirect(),
  })}`
}

type PictureUrlByIdOptions = {
  organizationId?: string
  groupId?: string
  invitationId?: string
}

export function getPictureUrlById(
  id: string,
  picture: Picture,
  options?: PictureUrlByIdOptions,
) {
  return `/proxy/api/v2/users/${id}/picture${picture.extension}?${new URLSearchParams(
    {
      access_token: getAccessTokenOrRedirect(),
      organization_id: options?.organizationId ?? '',
      group_id: options?.groupId ?? '',
      invitation_id: options?.invitationId ?? '',
    },
  )}`
}
