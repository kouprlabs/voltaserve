import { FileWithPath } from 'react-dropzone'
import FileAPI from '@/client/api/file'
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
    await FileAPI.upload({
      workspaceId: upload.workspaceId,
      parentId: upload.parentId,
      name: (upload.file as FileWithPath).path || upload.file.name,
      request,
      file: upload.file,
      onProgress: (progress) => {
        store.dispatch(uploadUpdated({ id: upload.id, progress }))
      },
    })
    store.dispatch(uploadCompleted(upload.id))
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
