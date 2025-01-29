// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { KeyedMutator } from 'swr'
import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import { FileList, FileSortBy, FileSortOrder } from '@/client/api/file'
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

export type FilesState = {
  selection: string[]
  hidden: string[]
  loading: string[]
  isMultiSelectActive: boolean
  isRangeSelectActive: boolean
  isMoveModalOpen: boolean
  isCopyModalOpen: boolean
  isCreateModalOpen: boolean
  isDeleteModalOpen: boolean
  isRenameModalOpen: boolean
  isSharingModalOpen: boolean
  isInfoModalOpen: boolean
  isFileUploadOpen: boolean
  isFolderUploadOpen: boolean
  isContextMenuOpen: boolean
  isSelectionMode: boolean
  iconScale: number
  sortBy: FileSortBy
  sortOrder: FileSortOrder
  viewType: FileViewType
  mutate?: KeyedMutator<FileList | undefined>
}

const initialState: FilesState = {
  selection: [],
  hidden: [],
  loading: [],
  isMultiSelectActive: false,
  isRangeSelectActive: false,
  isMoveModalOpen: false,
  isCopyModalOpen: false,
  isCreateModalOpen: false,
  isDeleteModalOpen: false,
  isRenameModalOpen: false,
  isSharingModalOpen: false,
  isFileUploadOpen: false,
  isFolderUploadOpen: false,
  isInfoModalOpen: false,
  isContextMenuOpen: false,
  iconScale: loadIconScale() || 1,
  sortBy: loadFileSortBy() || FileSortBy.DateCreated,
  sortOrder: loadFileSortOrder() || FileSortOrder.Desc,
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
    loadingAdded: (state, action: PayloadAction<string[]>) => {
      state.loading.push(...action.payload)
    },
    loadingRemoved: (state, action: PayloadAction<string[]>) => {
      state.loading = state.loading.filter((e) => !action.payload.includes(e))
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
      state.isSharingModalOpen = true
    },
    infoModalDidOpen: (state) => {
      state.isInfoModalOpen = true
    },
    fileUploadDidOpen: (state) => {
      state.isFileUploadOpen = true
    },
    folderUploadDidOpen: (state) => {
      state.isFolderUploadOpen = true
    },
    contextMenuDidOpen: (state) => {
      state.isContextMenuOpen = true
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
      state.isSharingModalOpen = false
    },
    infoModalDidClose: (state) => {
      state.isInfoModalOpen = false
    },
    fileUploadDidClose: (state) => {
      state.isFileUploadOpen = false
    },
    folderUploadDidClose: (state) => {
      state.isFolderUploadOpen = false
    },
    allModalsDidClose: (state) => {
      state.isMoveModalOpen = false
      state.isCopyModalOpen = false
      state.isCreateModalOpen = false
      state.isDeleteModalOpen = false
      state.isRenameModalOpen = false
      state.isSharingModalOpen = false
      state.isInfoModalOpen = false
    },
    contextMenuDidClose: (state) => {
      state.isContextMenuOpen = false
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
    sortByUpdated: (state, action: PayloadAction<FileSortBy>) => {
      state.sortBy = action.payload
      saveFileSortBy(state.sortBy)
    },
    sortOrderToggled: (state) => {
      state.sortOrder =
        state.sortOrder === FileSortOrder.Asc
          ? FileSortOrder.Desc
          : FileSortOrder.Asc
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
      action: PayloadAction<KeyedMutator<FileList | undefined>>,
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
  loadingAdded,
  loadingRemoved,
  moveModalDidOpen,
  copyModalDidOpen,
  createModalDidOpen,
  deleteModalDidOpen,
  renameModalDidOpen,
  sharingModalDidOpen,
  infoModalDidOpen,
  fileUploadDidOpen,
  folderUploadDidOpen,
  contextMenuDidOpen,
  moveModalDidClose,
  copyModalDidClose,
  createModalDidClose,
  deleteModalDidClose,
  renameModalDidClose,
  sharingModalDidClose,
  infoModalDidClose,
  fileUploadDidClose,
  folderUploadDidClose,
  allModalsDidClose,
  contextMenuDidClose,
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
