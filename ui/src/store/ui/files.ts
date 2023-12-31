import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import { SortBy, SortOrder } from '@/client/api/file'

export const SORT_BY_KEY = 'voltaserve_file_sort_by'
export const SORT_ORDER_KEY = 'voltaserve_file_sort_order'

export type FilesState = {
  selectedItems: string[]
  hiddenItems: string[]
  isMultiSelectActive: boolean
  isRangeSelectActive: boolean
  isMoveModalOpen: boolean
  isCopyModalOpen: boolean
  isCreateModalOpen: boolean
  isDeleteModalOpen: boolean
  isRenameModalOpen: boolean
  isShareModalOpen: boolean
  iconScale: number
  sortBy: SortBy
  sortOrder: SortOrder
}

const initialState: FilesState = {
  selectedItems: [],
  hiddenItems: [],
  isMultiSelectActive: false,
  isRangeSelectActive: false,
  isMoveModalOpen: false,
  isCopyModalOpen: false,
  isCreateModalOpen: false,
  isDeleteModalOpen: false,
  isRenameModalOpen: false,
  isShareModalOpen: false,
  iconScale: 1,
  sortBy: (localStorage.getItem(SORT_BY_KEY) as SortBy) || SortBy.DateCreated,
  sortOrder:
    (localStorage.getItem(SORT_ORDER_KEY) as SortOrder) || SortOrder.Desc,
}

const slice = createSlice({
  name: 'files',
  initialState,
  reducers: {
    selectedItemsUpdated: (state, action: PayloadAction<string[]>) => {
      state.selectedItems = action.payload
    },
    selectedItemAdded: (state, action: PayloadAction<string>) => {
      state.selectedItems.push(action.payload)
    },
    selectedItemRemoved: (state, action: PayloadAction<string>) => {
      state.selectedItems = state.selectedItems.filter(
        (e) => e !== action.payload,
      )
    },
    hiddenItemsUpdated: (state, action: PayloadAction<string[]>) => {
      state.hiddenItems = action.payload
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
    sortByUpdated: (state, action: PayloadAction<SortBy>) => {
      state.sortBy = action.payload
    },
    sortOrderUpdated: (state, action: PayloadAction<SortOrder>) => {
      state.sortOrder = action.payload
    },
  },
})

export const {
  selectedItemsUpdated,
  selectedItemAdded,
  selectedItemRemoved,
  hiddenItemsUpdated,
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
  sortByUpdated,
  sortOrderUpdated,
} = slice.actions

export default slice.reducer
