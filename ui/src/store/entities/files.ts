import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import { List } from '@/client/api/file'

type FilesState = {
  list?: List
}

const initialState: FilesState = {}

const slice = createSlice({
  name: 'files',
  initialState,
  reducers: {
    listUpdated: (state, action: PayloadAction<List>) => {
      state.list = action.payload
    },
  },
})

export const { listUpdated } = slice.actions

export default slice.reducer
