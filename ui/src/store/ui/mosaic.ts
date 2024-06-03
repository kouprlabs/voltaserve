import { createSlice } from '@reduxjs/toolkit'

type MosaicState = {
  isModalOpen: boolean
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
  },
})

export const { modalDidOpen, modalDidClose } = slice.actions

export default slice.reducer
