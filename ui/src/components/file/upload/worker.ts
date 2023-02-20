import { useParams } from 'react-router-dom'
import { errorToString } from '@/api/error'
import FileAPI from '@/api/file'
import store from '@/store/configure-store'
import { filesAdded } from '@/store/entities/files'
import {
  Upload,
  uploadCompleted,
  uploadUpdated,
} from '@/store/entities/uploads'

export const queue: Upload[] = []
let working = false

function getFileIdFromPath(): string {
  const segments = window.location.pathname.split('/')
  return segments[segments.length - 1]
}

setInterval(async () => {
  if (queue.length === 0 || working) {
    return
  }
  working = true
  const upload = queue.at(0) as Upload
  try {
    const request = new XMLHttpRequest()
    store.dispatch(uploadUpdated({ id: upload.id, request }))
    const result = await FileAPI.upload(
      upload.workspaceId,
      upload.parentId,
      request,
      upload.file,
      (progress) => {
        store.dispatch(uploadUpdated({ id: upload.id, progress }))
      }
    )
    const fileId = getFileIdFromPath()
    if (upload.parentId === fileId) {
      store.dispatch(filesAdded({ id: fileId, files: [result] }))
    }
    store.dispatch(uploadCompleted(upload.id))
  } catch (error) {
    store.dispatch(
      uploadUpdated({
        id: upload.id,
        completed: true,
        error: errorToString(error),
      })
    )
  } finally {
    queue.shift()
    working = false
  }
}, 1000)
