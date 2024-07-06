// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

import { IncomingMessage } from 'http'

export function getTargetPath(req: IncomingMessage) {
  const destination = req.headers.destination as string
  if (!destination) {
    return null
  }
  /* Check if the destination header is a full URL */
  if (destination.startsWith('http://') || destination.startsWith('https://')) {
    return new URL(destination).pathname
  } else {
    /* Extract the path from the destination header */
    const startIndex =
      destination.indexOf(req.headers.host) + req.headers.host.length
    return destination.substring(startIndex)
  }
}
