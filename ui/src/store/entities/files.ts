import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import { File, FileList } from '@/api/file'

type FilesState = {
  current?: string
  folder?: File
  list?: FileList
  sortType: SortType
  sortDirection: SortDirection
}

export enum SortType {
  ByName = 'by_name',
  ByDateCreated = 'by_date_created',
  ByDateModified = 'by_date_modified',
}

export enum SortDirection {
  Ascending = 'ascending',
  Descending = 'descending',
}

const initialState: FilesState = {
  sortType: SortType.ByDateCreated,
  sortDirection: SortDirection.Ascending,
}

function sort(data: File[], type: SortType, direction: SortDirection) {
  const lt = direction === SortDirection.Ascending ? -1 : 1
  const gt = direction === SortDirection.Ascending ? 1 : -1
  if (type === SortType.ByName) {
    data.sort((a, b) => {
      if (a.name.toLowerCase() < b.name.toLowerCase()) {
        return lt
      }
      if (a.name.toLowerCase() > b.name.toLowerCase()) {
        return gt
      }
      return 0
    })
  } else if (type === SortType.ByDateCreated) {
    data.sort((a, b) => {
      const dateA = new Date(a.createTime)
      const dateB = new Date(b.createTime)
      if (dateA < dateB) {
        return lt
      }
      if (dateA > dateB) {
        return gt
      }
      return 0
    })
  } else if (type === SortType.ByDateModified) {
    data.sort((a, b) => {
      const dateA = new Date(a.updateTime || a.createTime)
      const dateB = new Date(b.updateTime || b.createTime)
      if (dateA < dateB) {
        return lt
      }
      if (dateA > dateB) {
        return gt
      }
      return 0
    })
  }
}

const slice = createSlice({
  name: 'files',
  initialState,
  reducers: {
    filesAdded: (
      state,
      action: PayloadAction<{ id: string; files: File[] }>
    ) => {
      if (state.list && state.current === action.payload.id) {
        state.list.data.push(...action.payload.files)
      }
      if (state.list) {
        sort(state.list.data, state.sortType, state.sortDirection)
      }
    },
    filesRemoved: (
      state,
      action: PayloadAction<{ id: string; files: string[] }>
    ) => {
      if (state.list && state.current === action.payload.id) {
        state.list.data = state.list.data.filter(
          (e) => action.payload.files.findIndex((id) => e.id === id) === -1
        )
      }
      if (state.list) {
        sort(state.list.data, state.sortType, state.sortDirection)
      }
    },
    filesUpdated: (state, action: PayloadAction<File[]>) => {
      action.payload.forEach((e) => {
        if (!state.list) {
          return
        }
        const file = state.list.data.find((f) => f.id === e.id)
        if (file) {
          Object.assign(file, e)
        }
      })
      if (state.list) {
        sort(state.list.data, state.sortType, state.sortDirection)
      }
    },
    listUpdated: (state, action: PayloadAction<FileList>) => {
      state.list = action.payload
      if (state.list) {
        sort(state.list.data, state.sortType, state.sortDirection)
      }
    },
    listPatched: (state, action: PayloadAction<FileList>) => {
      if (!state.list) {
        return
      }
      state.list.data.push(...action.payload.data)
      state.list.page = action.payload.page
      state.list.size = action.payload.size
      state.list.totalElements = action.payload.totalElements
      state.list.totalPages = action.payload.totalPages
      if (state.list) {
        sort(state.list.data, state.sortType, state.sortDirection)
      }
    },
    folderUpdated: (state, action: PayloadAction<File>) => {
      state.folder = action.payload
    },
    currentUpdated: (state, action: PayloadAction<string>) => {
      state.current = action.payload
    },
    sortTypeUpdated: (state, action: PayloadAction<SortType>) => {
      state.sortType = action.payload
      if (state.list) {
        sort(state.list.data, state.sortType, state.sortDirection)
      }
    },
    sortDirectionUpdated: (state, action: PayloadAction<SortDirection>) => {
      state.sortDirection = action.payload
      if (state.list) {
        sort(state.list.data, state.sortType, state.sortDirection)
      }
    },
  },
})

export const {
  filesAdded,
  filesRemoved,
  filesUpdated,
  listUpdated,
  listPatched,
  folderUpdated,
  currentUpdated,
  sortTypeUpdated,
  sortDirectionUpdated,
} = slice.actions

export default slice.reducer
