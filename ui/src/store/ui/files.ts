import { createSlice, PayloadAction } from '@reduxjs/toolkit'

export type FilesState = {
  selection: string[]
  isMultiSelectActive: boolean
  isRangeSelectActive: boolean
  isMoveModalOpen: boolean
  isCopyModalOpen: boolean
  isCreateModalOpen: boolean
  isDeleteModalOpen: boolean
  isRenameModalOpen: boolean
  isShareModalOpen: boolean
  iconScale: number
  sortType: SortType
  sortDirection: SortDirection
}

export enum SortType {
  None = 'none',
  ByName = 'by_name',
  ByDateCreated = 'by_date_created',
  ByDateModified = 'by_date_modified',
}

export enum SortDirection {
  Ascending = 'ascending',
  Descending = 'descending',
}

const initialState: FilesState = {
  selection: [],
  isMultiSelectActive: false,
  isRangeSelectActive: false,
  isMoveModalOpen: false,
  isCopyModalOpen: false,
  isCreateModalOpen: false,
  isDeleteModalOpen: false,
  isRenameModalOpen: false,
  isShareModalOpen: false,
  iconScale: 1,
  sortType: SortType.None,
  sortDirection: SortDirection.Descending,
}

const slice = createSlice({
  name: 'files',
  initialState,
  reducers: {
    selectionUpdated: (state, action: PayloadAction<string[]>) => {
      state.selection = action.payload
    },
    selectionAdded: (state, action: PayloadAction<string>) => {
      state.selection.push(action.payload)
    },
    selectionRemoved: (state, action: PayloadAction<string>) => {
      state.selection = state.selection.filter((e) => e !== action.payload)
    },
    moveModalDidOpen: (state) => {
      state.isMoveModalOpen = true
    },
    copyModalDidOpen: (state) => {
      state.isCopyModalOpen = true
    },
    createModalDidOpen: (state) => {
      state.isCreateModalOpen = true
    },
    deleteModalDidOpen: (state) => {
      state.isDeleteModalOpen = true
    },
    renameModalDidOpen: (state) => {
      state.isRenameModalOpen = true
    },
    sharingModalDidOpen: (state) => {
      state.isShareModalOpen = true
    },
    moveModalDidClose: (state) => {
      state.isMoveModalOpen = false
    },
    copyModalDidClose: (state) => {
      state.isCopyModalOpen = false
    },
    createModalDidClose: (state) => {
      state.isCreateModalOpen = false
    },
    deleteModalDidClose: (state) => {
      state.isDeleteModalOpen = false
    },
    renameModalDidClose: (state) => {
      state.isRenameModalOpen = false
    },
    sharingModalDidClose: (state) => {
      state.isShareModalOpen = false
    },
    multiSelectKeyUpdated: (state, action: PayloadAction<boolean>) => {
      state.isMultiSelectActive = action.payload
    },
    rangeSelectKeyUpdated: (state, action: PayloadAction<boolean>) => {
      state.isRangeSelectActive = action.payload
    },
    iconScaleUpdated: (state, action: PayloadAction<number>) => {
      state.iconScale = action.payload
    },
    sortTypeUpdated: (state, action: PayloadAction<SortType>) => {
      state.sortType = action.payload
    },
    sortDirectionUpdated: (state, action: PayloadAction<SortDirection>) => {
      state.sortDirection = action.payload
    },
  },
})

export const {
  selectionUpdated,
  selectionAdded,
  selectionRemoved,
  moveModalDidOpen,
  copyModalDidOpen,
  createModalDidOpen,
  deleteModalDidOpen,
  renameModalDidOpen,
  sharingModalDidOpen,
  moveModalDidClose,
  copyModalDidClose,
  createModalDidClose,
  deleteModalDidClose,
  renameModalDidClose,
  sharingModalDidClose,
  multiSelectKeyUpdated,
  rangeSelectKeyUpdated,
  iconScaleUpdated,
  sortTypeUpdated,
  sortDirectionUpdated,
} = slice.actions

export default slice.reducer
