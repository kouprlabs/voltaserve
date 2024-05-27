import { KeyedMutator } from 'swr'
import { PayloadAction, createSlice } from '@reduxjs/toolkit'
import { File } from '@/client/api/file'

type InsightsState = {
  isModalOpen: boolean
  mutateFile?: KeyedMutator<File>
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
    mutateFileUpdated: (state, action: PayloadAction<KeyedMutator<File>>) => {
      state.mutateFile = action.payload
    },
  },
})

export const { modalDidOpen, modalDidClose, mutateFileUpdated } = slice.actions

export default slice.reducer
