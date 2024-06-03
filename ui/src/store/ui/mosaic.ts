import { KeyedMutator } from 'swr'
import { PayloadAction, createSlice } from '@reduxjs/toolkit'
import { Metadata } from '@/client/api/mosaic'

type MosaicState = {
  isModalOpen: boolean
  mutateMetadata?: KeyedMutator<Metadata>
}

const initialState: MosaicState = {
  isModalOpen: false,
}

const slice = createSlice({
  name: 'mosaic',
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
