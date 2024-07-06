// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

import useSWR, { SWRConfiguration } from 'swr'
import { apiFetcher } from '@/client/fetcher'
import { Snapshot } from './snapshot'

export type Info = {
  isAvailable: boolean
  isOutdated: boolean
  snapshot?: Snapshot
}

export default class WatermarkAPI {
  static create(id: string, showError = true) {
    return apiFetcher({
      url: `/watermarks/${id}`,
      method: 'POST',
      showError,
    })
  }

  static delete(id: string) {
    return apiFetcher({
      url: `/watermarks/${id}`,
      method: 'DELETE',
    })
  }

  static useGetInfo(
    id: string | null | undefined,
    swrOptions?: SWRConfiguration,
  ) {
    const url = `/watermarks/${id}/info`
    return useSWR<Info>(
      id ? url : null,
      () =>
        apiFetcher({
          url,
          method: 'GET',
          showError: false,
        }) as Promise<Info>,
      swrOptions,
    )
  }

  static async getInfo(id: string) {
    return apiFetcher({
      url: `/mosaics/${id}/info`,
      method: 'GET',
    }) as Promise<Info>
  }
}
