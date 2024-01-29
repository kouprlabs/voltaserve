import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import { FileWithPath } from 'react-dropzone'
import { newHashId } from '@/infra/id'

export type Upload = {
  id: string
  workspaceId: string
  parentId: string
  file: File | FileWithPath
  request?: any
  progress?: number
  error?: string
  completed: boolean
}

export class UploadDecorator {
  public value: Upload

  constructor(options: {
    id?: string
    workspaceId: string
    parentId: string
    file: File | FileWithPath
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

  get workspaceId(): string {
    return this.value.workspaceId
  }

  get parentId(): string {
    return this.value.parentId
  }

  get file(): File {
    return this.value.file
  }

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
