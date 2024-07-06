// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

import { IncomingMessage, ServerResponse } from 'http'

/*
  This method should respond with the allowed methods and capabilities of the server.

  Example implementation:

  - Set the response status code to 200.
  - Set the Allow header to specify the supported methods, such as OPTIONS, GET, PUT, DELETE, etc.
  - Return the response.
 */
async function handleOptions(_: IncomingMessage, res: ServerResponse) {
  res.statusCode = 200
  res.setHeader(
    'Allow',
    'OPTIONS, GET, HEAD, PUT, DELETE, MKCOL, COPY, MOVE, PROPFIND, PROPPATCH',
  )
  res.end()
}

export default handleOptions
