import { createSlice } from '@reduxjs/toolkit'

type SearchFilterState = {
  isModalOpen: boolean
}

const initialState: SearchFilterState = {
  isModalOpen: false,
}

const slice = createSlice({
  name: 'search-filter',
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
