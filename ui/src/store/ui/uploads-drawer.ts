import { createSlice } from '@reduxjs/toolkit'

export type UploadsDrawerState = {
  open: boolean
}

const initialState: UploadsDrawerState = {
  open: false,
}

const slice = createSlice({
  name: 'uploads-drawer',
  initialState,
  reducers: {
    uploadsDrawerOpened: (state) => {
      state.open = true
    },
    uploadsDrawerClosed: (state) => {
      state.open = false
    },
  },
})

export const { uploadsDrawerOpened, uploadsDrawerClosed } = slice.actions

export default slice.reducer
