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
import path from 'path'
import { FileAPI } from '@/client/api'
import { Token } from '@/client/idp'
import { getTargetPath } from '@/helper/path'
import { handleError } from '@/infra/error'

/*
  This method moves or renames a resource from a source URL to a destination URL.

  Example implementation:

  - Extract the source and destination paths from the headers or request body.
  - Use fs.rename() to move or rename the file from the source to the destination.
  - Set the response status code to 204 if successful or an appropriate error code if the source file is not found or encountered an error.
  - Return the response.
 */
async function handleMove(
  req: IncomingMessage,
  res: ServerResponse,
  token: Token,
) {
  try {
    const sourcePath = decodeURIComponent(req.url)
    const targetPath = decodeURIComponent(getTargetPath(req))
    const api = new FileAPI(token)
    const sourceFile = await api.getByPath(decodeURIComponent(req.url))
    const targetFile = await api.getByPath(
      decodeURIComponent(path.dirname(getTargetPath(req))),
    )
    if (sourceFile.workspaceId !== targetFile.workspaceId) {
      res.statusCode = 400
      res.end()
    } else {
      if (
        sourcePath.split('/').length === targetPath.split('/').length &&
        path.dirname(sourcePath) === path.dirname(targetPath)
      ) {
        await api.patchName(sourceFile.id, {
          name: decodeURIComponent(path.basename(targetPath)),
        })
      } else {
        await api.move(targetFile.id, { ids: [sourceFile.id] })
      }
      res.statusCode = 204
      res.end()
    }
  } catch (err) {
    handleError(err, res)
  }
}

export default handleMove
