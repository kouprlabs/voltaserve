// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { FileWithPath } from 'react-dropzone'
import FileAPI, { FileType } from '@/client/api/file'
import { errorToString } from '@/client/error'
import store from '@/store/configure-store'
import {
  Upload,
  uploadCompleted,
  uploadUpdated,
} from '@/store/entities/uploads'

export const queue: Upload[] = []
let working = false

setInterval(async () => {
  if (queue.length === 0 || working) {
    return
  }
  working = true
  const upload = queue.at(0) as Upload
  try {
    const request = new XMLHttpRequest()
    store.dispatch(uploadUpdated({ id: upload.id, request }))
    if (upload.fileId) {
      await FileAPI.patch({
        id: upload.fileId,
        request,
        blob: upload.blob,
        onProgress: (progress) => {
          store.dispatch(uploadUpdated({ id: upload.id, progress }))
        },
      })
    } else if (upload.workspaceId && upload.parentId) {
      await FileAPI.create({
        type: FileType.File,
        workspaceId: upload.workspaceId,
        parentId: upload.parentId,
        name:
          (upload.blob as FileWithPath).path ||
          upload.blob.webkitRelativePath ||
          upload.blob.name,
        request,
        blob: upload.blob,
        onProgress: (progress) => {
          store.dispatch(uploadUpdated({ id: upload.id, progress }))
        },
      })
    }
    store.dispatch(uploadCompleted(upload.id))
    await store.getState().ui.files.mutate?.()
  } catch (error) {
    store.dispatch(
      uploadUpdated({
        id: upload.id,
        completed: true,
        error: errorToString(error),
      }),
    )
  } finally {
    queue.shift()
    working = false
  }
}, 1000)
