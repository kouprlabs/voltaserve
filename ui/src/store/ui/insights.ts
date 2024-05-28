import { KeyedMutator } from 'swr'
import { PayloadAction, createSlice } from '@reduxjs/toolkit'
import { Metadata } from '@/client/api/insights'

type InsightsState = {
  isModalOpen: boolean
  mutateMetadata?: KeyedMutator<Metadata>
}

const initialState: InsightsState = {
  isModalOpen: false,
}

const slice = createSlice({
  name: 'insights',
  initialState,
  reducers: {
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

export const { modalDidOpen, modalDidClose, mutateMetadataUpdated } =
  slice.actions

export default slice.reducer
