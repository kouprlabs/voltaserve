// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

import { ServerResponse } from 'http'
import { APIError } from '@/client/api'
import { IdPError } from '@/client/idp'

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export function handleError(err: any, res: ServerResponse) {
  if (err instanceof APIError) {
    res.statusCode = err.error.status
    res.statusMessage = err.error.userMessage
    res.end()
  } else if (err instanceof IdPError) {
    res.statusCode = err.error.status
    res.statusMessage = err.error.userMessage
    res.end()
  } else {
    res.statusCode = 500
    res.end()
  }
  console.error(err)
}
