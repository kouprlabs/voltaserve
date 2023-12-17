import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import { File, List } from '@/client/api/file'

type FilesState = {
  current?: string
  folder?: File
  list?: List
}

const initialState: FilesState = {}

const slice = createSlice({
  name: 'files',
  initialState,
  reducers: {
    filesAdded: (
      state,
      action: PayloadAction<{ id: string; files: File[] }>,
    ) => {
      if (state.list && state.current === action.payload.id) {
        state.list.data.push(...action.payload.files)
      }
    },
    filesRemoved: (
      state,
      action: PayloadAction<{ id: string; files: string[] }>,
    ) => {
      if (state.list && state.current === action.payload.id) {
        state.list.data = state.list.data.filter(
          (e) => action.payload.files.findIndex((id) => e.id === id) === -1,
        )
      }
    },
    filesUpdated: (state, action: PayloadAction<File[]>) => {
      action.payload.forEach((updatedFile) => {
        if (!state.list) {
          return
        }
        const index = state.list.data.findIndex((f) => f.id === updatedFile.id)
        if (index !== -1) {
          state.list.data[index] = updatedFile
        }
      })
    },
    listUpdated: (state, action: PayloadAction<List>) => {
      state.list = action.payload
    },
    listExtended: (state, action: PayloadAction<List>) => {
      if (!state.list) {
        return
      }
      state.list.data.push(...action.payload.data)
      state.list.page = action.payload.page
      state.list.size = action.payload.size
      state.list.totalElements = action.payload.totalElements
      state.list.totalPages = action.payload.totalPages
    },
    folderUpdated: (state, action: PayloadAction<File>) => {
      state.folder = action.payload
    },
    currentUpdated: (state, action: PayloadAction<string>) => {
      state.current = action.payload
    },
  },
})

export const {
  filesAdded,
  filesRemoved,
  filesUpdated,
  listUpdated,
  listExtended,
  folderUpdated,
  currentUpdated,
} = slice.actions

export default slice.reducer
