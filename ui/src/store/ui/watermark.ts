import { createSlice } from '@reduxjs/toolkit'

type WatermarkState = {
  isModalOpen: boolean
}

const initialState: WatermarkState = {
  isModalOpen: false,
}

const slice = createSlice({
  name: 'watermark',
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
