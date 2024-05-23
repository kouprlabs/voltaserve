import { KeyedMutator } from 'swr'
import { PayloadAction, createSlice } from '@reduxjs/toolkit'
import { Summary } from '@/client/api/insights'

type InsightsState = {
  isModalOpen: boolean
  mutateSummary?: KeyedMutator<Summary>
}

const initialState: InsightsState = {
  isModalOpen: false,
}

const slice = createSlice({
  name: 'analysis',
  initialState,
  reducers: {
    modalDidOpen: (state) => {
      state.isModalOpen = true
    },
    modalDidClose: (state) => {
      state.isModalOpen = false
    },
    mutateSummaryUpdated: (
      state,
      action: PayloadAction<KeyedMutator<Summary>>,
    ) => {
      state.mutateSummary = action.payload
    },
  },
})

export const { modalDidOpen, modalDidClose, mutateSummaryUpdated } =
  slice.actions

export default slice.reducer
