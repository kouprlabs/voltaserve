import { KeyedMutator } from 'swr'
import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import { List, SortBy, SortOrder } from '@/client/api/file'
import { List as SnapshotList } from '@/client/api/snapshot'
import {
  loadFileSortBy,
  loadFileSortOrder,
  loadFileViewType,
  loadIconScale,
  saveFileSortBy,
  saveFileSortOrder,
  saveFileViewType,
  saveIconScale,
} from '@/local-storage'
import { FileViewType } from '@/types/file'

export const SORT_ORDER_KEY = 'voltaserve_file_sort_order'

export type FilesState = {
  selection: string[]
  hidden: string[]
  snapshotSelection: string[]
  isMultiSelectActive: boolean
  isRangeSelectActive: boolean
  isMoveModalOpen: boolean
  isCopyModalOpen: boolean
  isCreateModalOpen: boolean
  isDeleteModalOpen: boolean
  isRenameModalOpen: boolean
  isShareModalOpen: boolean
  isSnapshotListModalOpen: boolean
  isSnapshotDeleteModalOpen: boolean
  isSelectionMode: boolean
  iconScale: number
  sortBy: SortBy
  sortOrder: SortOrder
  viewType: FileViewType
  mutate?: KeyedMutator<List | undefined>
  snapshotMutate?: KeyedMutator<SnapshotList | undefined>
}

const initialState: FilesState = {
  selection: [],
  hidden: [],
  snapshotSelection: [],
  isMultiSelectActive: false,
  isRangeSelectActive: false,
  isMoveModalOpen: false,
  isCopyModalOpen: false,
  isCreateModalOpen: false,
  isDeleteModalOpen: false,
  isRenameModalOpen: false,
  isShareModalOpen: false,
  isSnapshotListModalOpen: false,
  isSnapshotDeleteModalOpen: false,
  iconScale: loadIconScale() || 1,
  sortBy: loadFileSortBy() || SortBy.DateCreated,
  sortOrder: loadFileSortOrder() || SortOrder.Desc,
  viewType: loadFileViewType() || FileViewType.Grid,
  isSelectionMode: false,
}

const slice = createSlice({
  name: 'files',
  initialState,
  reducers: {
    selectionUpdated: (state, action: PayloadAction<string[]>) => {
      state.selection = action.payload
    },
    selectionAdded: (state, action: PayloadAction<string>) => {
      if (!state.selection.includes(action.payload)) {
        state.selection.push(action.payload)
      }
    },
    selectionRemoved: (state, action: PayloadAction<string>) => {
      state.selection = state.selection.filter((e) => e !== action.payload)
    },
    hiddenUpdated: (state, action: PayloadAction<string[]>) => {
      state.hidden = action.payload
    },
    snapshotSelectionUpdated: (state, action: PayloadAction<string[]>) => {
      state.snapshotSelection = action.payload
    },
    snapshotMutateUpdated: (
      state,
      action: PayloadAction<KeyedMutator<SnapshotList | undefined>>,
    ) => {
      state.snapshotMutate = action.payload
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
    snapshotListModalDidOpen: (state) => {
      state.isSnapshotListModalOpen = true
    },
    snapshotDeleteModalDidOpen: (state) => {
      state.isSnapshotDeleteModalOpen = true
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
    snapshotListModalDidClose: (state) => {
      state.isSnapshotListModalOpen = false
    },
    snapshotDeleteModalDidClose: (state) => {
      state.isSnapshotDeleteModalOpen = false
    },
    multiSelectKeyUpdated: (state, action: PayloadAction<boolean>) => {
      state.isMultiSelectActive = action.payload
    },
    rangeSelectKeyUpdated: (state, action: PayloadAction<boolean>) => {
      state.isRangeSelectActive = action.payload
    },
    iconScaleUpdated: (state, action: PayloadAction<number>) => {
      state.iconScale = action.payload
      saveIconScale(state.iconScale)
    },
    sortByUpdated: (state, action: PayloadAction<SortBy>) => {
      state.sortBy = action.payload
      saveFileSortBy(state.sortBy)
    },
    sortOrderToggled: (state) => {
      state.sortOrder =
        state.sortOrder === SortOrder.Asc ? SortOrder.Desc : SortOrder.Asc
      saveFileSortOrder(state.sortOrder)
    },
    viewTypeToggled: (state) => {
      state.viewType =
        state.viewType === FileViewType.Grid
          ? FileViewType.List
          : FileViewType.Grid
      saveFileViewType(state.viewType)
    },
    selectionModeToggled: (state) => {
      state.isSelectionMode = !state.isSelectionMode
    },
    mutateUpdated: (
      state,
      action: PayloadAction<KeyedMutator<List | undefined>>,
    ) => {
      state.mutate = action.payload
    },
  },
})

export const {
  selectionUpdated,
  selectionAdded,
  selectionRemoved,
  hiddenUpdated,
  snapshotSelectionUpdated,
  snapshotMutateUpdated,
  moveModalDidOpen,
  copyModalDidOpen,
  createModalDidOpen,
  deleteModalDidOpen,
  renameModalDidOpen,
  sharingModalDidOpen,
  snapshotListModalDidOpen,
  snapshotDeleteModalDidOpen,
  moveModalDidClose,
  copyModalDidClose,
  createModalDidClose,
  deleteModalDidClose,
  renameModalDidClose,
  sharingModalDidClose,
  snapshotListModalDidClose,
  snapshotDeleteModalDidClose,
  multiSelectKeyUpdated,
  rangeSelectKeyUpdated,
  iconScaleUpdated,
  sortByUpdated,
  sortOrderToggled,
  viewTypeToggled,
  selectionModeToggled,
  mutateUpdated,
} = slice.actions

export default slice.reducer
