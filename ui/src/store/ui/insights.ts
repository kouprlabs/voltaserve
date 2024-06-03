import { createSlice } from '@reduxjs/toolkit'

type InsightsState = {
  isModalOpen: boolean
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
  },
})

export const { modalDidOpen, modalDidClose } = slice.actions

export default slice.reducer
