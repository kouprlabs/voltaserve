// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { ErrorResponse } from '@/error/core.ts'
import { UserDTO } from '@/user/service.ts'
import { getConfig } from '@/config/config.ts'
import { newInternalServerError } from '@/error/creators.ts'

export enum UserWebhookEventType {
  Create = 'create',
  Delete = 'delete',
}

export type UserWebhookOptions = {
  eventType: UserWebhookEventType
  user: UserDTO
}

export async function call(
  url: string,
  opts: UserWebhookOptions,
): Promise<void> {
  const response = await fetch(
    `${url}?api_key=${getConfig().security.apiKey}`,
    {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json; charset=UTF-8',
      },
      body: JSON.stringify(opts),
    },
  )
  await successfulResponseOrError(response)
}

async function successfulResponseOrError(response: Response): Promise<void> {
  if (response.status > 299) {
    try {
      throw await response.json() as ErrorResponse
    } catch (error) {
      throw newInternalServerError(error)
    }
  }
}
