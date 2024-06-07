import { KeyedMutator } from 'swr'
import { PayloadAction, createSlice } from '@reduxjs/toolkit'
import { Info } from '@/client/api/insights'

type InsightsState = {
  isModalOpen: boolean
  mutateInfo?: KeyedMutator<Info>
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
    mutateInfoUpdated: (state, action: PayloadAction<KeyedMutator<Info>>) => {
      state.mutateInfo = action.payload
    },
  },
})

export const { modalDidOpen, modalDidClose, mutateInfoUpdated } = slice.actions

export default slice.reducer
