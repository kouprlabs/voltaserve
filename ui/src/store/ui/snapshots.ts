import { KeyedMutator } from 'swr'
import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import { List } from '@/client/api/snapshot'

export type SnapshotsState = {
  selection: string[]
  isListModalOpen: boolean
  isDeleteModalOpen: boolean
  snapshotMutate?: KeyedMutator<List | undefined>
}

const initialState: SnapshotsState = {
  selection: [],
  isListModalOpen: false,
  isDeleteModalOpen: false,
}

const slice = createSlice({
  name: 'snapshots',
  initialState,
  reducers: {
    selectionUpdated: (state, action: PayloadAction<string[]>) => {
      state.selection = action.payload
    },
    mutateUpdated: (
      state,
      action: PayloadAction<KeyedMutator<List | undefined>>,
    ) => {
      state.snapshotMutate = action.payload
    },
    listModalDidOpen: (state) => {
      state.isListModalOpen = true
    },
    deleteModalDidOpen: (state) => {
      state.isDeleteModalOpen = true
    },
    listModalDidClose: (state) => {
      state.isListModalOpen = false
    },
    deleteModalDidClose: (state) => {
      state.isDeleteModalOpen = false
    },
  },
})

export const {
  selectionUpdated,
  mutateUpdated,
  listModalDidOpen,
  deleteModalDidOpen,
  listModalDidClose,
  deleteModalDidClose,
} = slice.actions

export default slice.reducer
