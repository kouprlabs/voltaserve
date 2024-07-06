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
import { File, FileAPI, FileType } from '@/client/api'
import { Token } from '@/client/idp'
import { handleError } from '@/infra/error'

/*
  This method creates a new collection (directory) at the specified URL.

  Example implementation:

  - Extract the directory path from the URL.
  - Use fs.mkdir() to create the directory.
  - Set the response status code to 201 if created or an appropriate error code if the directory already exists or encountered an error.
  - Return the response.
 */
async function handleMkcol(
  req: IncomingMessage,
  res: ServerResponse,
  token: Token,
) {
  let directory: File
  try {
    const api = new FileAPI(token)
    directory = await api.getByPath(decodeURIComponent(path.dirname(req.url)))
    await api.create({
      type: FileType.Folder,
      workspaceId: directory.workspaceId,
      parentId: directory.id,
      name: decodeURIComponent(path.basename(req.url)),
    })
    res.statusCode = 201
    res.end()
  } catch (err) {
    handleError(err, res)
  }
}

export default handleMkcol
