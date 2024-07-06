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
import { FileAPI } from '@/client/api'
import { Token } from '@/client/idp'
import { handleError } from '@/infra/error'

/*
  This method deletes a resource identified by the URL.

  Example implementation:

  - Extract the file path from the URL.
  - Use fs.unlink() to delete the file.
  - Set the response status code to 204 if successful or an appropriate error code if the file is not found.
  - Return the response.
 */
async function handleDelete(
  req: IncomingMessage,
  res: ServerResponse,
  token: Token,
) {
  try {
    const api = new FileAPI(token)
    const file = await api.getByPath(decodeURIComponent(req.url))
    await api.delete(file.id)
    res.statusCode = 204
    res.end()
  } catch (err) {
    handleError(err, res)
  }
}

export default handleDelete
