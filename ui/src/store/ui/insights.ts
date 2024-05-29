import { KeyedMutator } from 'swr'
import { PayloadAction, createSlice } from '@reduxjs/toolkit'
import { Metadata } from '@/client/api/insights'

type InsightsState = {
  isModalOpen: boolean
  isCreating: boolean
  isUpdating: boolean
  isDeleting: boolean
  mutateMetadata?: KeyedMutator<Metadata>
}

const initialState: InsightsState = {
  isCreating: false,
  isUpdating: false,
  isDeleting: false,
  isModalOpen: false,
}

const slice = createSlice({
  name: 'insights',
  initialState,
  reducers: {
    creatingDidStart: (state) => {
      state.isCreating = true
    },
    creatingDidStop: (state) => {
      state.isCreating = false
    },
    updatingDidStart: (state) => {
      state.isUpdating = true
    },
    updatingDidStop: (state) => {
      state.isUpdating = false
    },
    deletingDidStart: (state) => {
      state.isDeleting = true
    },
    deletingDidStop: (state) => {
      state.isDeleting = false
    },
    modalDidOpen: (state) => {
      state.isModalOpen = true
    },
    modalDidClose: (state) => {
      state.isModalOpen = false
    },
    mutateMetadataUpdated: (
      state,
      action: PayloadAction<KeyedMutator<Metadata>>,
    ) => {
      state.mutateMetadata = action.payload
    },
  },
})

export const {
  modalDidOpen,
  modalDidClose,
  mutateMetadataUpdated,
  creatingDidStart,
  creatingDidStop,
  updatingDidStart,
  updatingDidStop,
  deletingDidStart,
  deletingDidStop,
} = slice.actions

export default slice.reducer
