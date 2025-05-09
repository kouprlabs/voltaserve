// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import { FileWithPath } from 'react-dropzone'
import { newHashId } from '@/infra/id'

export type Upload = {
  id: string
  fileId?: string
  workspaceId?: string
  parentId?: string
  blob: File | FileWithPath
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  request?: any
  progress?: number
  error?: string
  completed: boolean
}

export class UploadDecorator {
  public value: Upload

  constructor(options: {
    id?: string
    fileId?: string
    workspaceId?: string
    parentId?: string
    blob: File | FileWithPath
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    request?: any
    progress?: number
    error?: string
    completed?: boolean
  }) {
    this.value = {
      id: newHashId(),
      completed: false,
      ...options,
    }
  }

  get id(): string {
    return this.value.id
  }

  get fileId(): string | undefined {
    return this.value.fileId
  }

  get workspaceId(): string | undefined {
    return this.value.workspaceId
  }

  get parentId(): string | undefined {
    return this.value.parentId
  }

  get blob(): File {
    return this.value.blob
  }

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  get request(): any {
    return this.value.request
  }

  get progress(): number | undefined {
    return this.value.progress
  }

  get error(): string | undefined {
    return this.value.error
  }

  get completed(): boolean {
    return this.value.completed
  }

  get isProgressing() {
    return !this.value.completed && this.value.request
  }

  get isPending() {
    return !this.value.completed && !this.value.request
  }

  get isSucceeded() {
    return this.value.completed && !this.value.error
  }

  get isFailed() {
    return this.value.completed && this.value.error
  }
}

export type UploadsState = {
  items: Upload[]
}

type UploadUpdateOptions = {
  id: string
  workspaceId?: string
  parentId?: string
  file?: File
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  request?: any
  progress?: number
  error?: string
  completed?: boolean
}

const initialState: UploadsState = {
  items: [],
}

const slice = createSlice({
  name: 'uploads',
  initialState,
  reducers: {
    uploadAdded: (state, action: PayloadAction<Upload>) => {
      state.items.unshift(action.payload)
    },
    uploadUpdated: (state, action: PayloadAction<UploadUpdateOptions>) => {
      const upload = state.items.find((e) => e.id === action.payload.id)
      if (upload) {
        Object.assign(upload, action.payload)
      }
    },
    uploadCompleted: (state, action: PayloadAction<string>) => {
      const index = state.items.findIndex((e) => e.id === action.payload)
      if (index !== -1) {
        state.items[index].completed = true
        state.items.push(state.items.splice(index, 1)[0])
      }
    },
    uploadRemoved: (state, action: PayloadAction<string>) => {
      state.items = state.items.filter((e) => e.id !== action.payload)
    },
    completedUploadsCleared: (state) => {
      state.items = state.items.filter((e) => !e.completed)
    },
  },
})

export const {
  uploadAdded,
  uploadUpdated,
  uploadRemoved,
  uploadCompleted,
  completedUploadsCleared,
} = slice.actions

export default slice.reducer
