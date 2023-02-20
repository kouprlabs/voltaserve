import { createSlice, PayloadAction } from '@reduxjs/toolkit'

type ErrorState = {
  value: string | null
}

const initialState: ErrorState = {
  value: null,
}

const slice = createSlice({
  name: 'error',
  initialState,
  reducers: {
    errorOccurred: (state, action: PayloadAction<string>) => {
      state.value = action.payload
    },
    errorCleared: (state) => {
      state.value = null
    },
  },
})

export const { errorOccurred, errorCleared } = slice.actions

export default slice.reducer
