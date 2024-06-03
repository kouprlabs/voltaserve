import { createSlice } from '@reduxjs/toolkit'

export type UploadsState = {
  isDrawerOpen: boolean
}

const initialState: UploadsState = {
  isDrawerOpen: false,
}

const slice = createSlice({
  name: 'uploads',
  initialState,
  reducers: {
    drawerDidOpen: (state) => {
      state.isDrawerOpen = true
    },
    drawerDidClose: (state) => {
      state.isDrawerOpen = false
    },
  },
})

export const { drawerDidOpen, drawerDidClose } = slice.actions

export default slice.reducer
