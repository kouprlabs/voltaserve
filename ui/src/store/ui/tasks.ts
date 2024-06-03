import { KeyedMutator } from 'swr'
import { PayloadAction, createSlice } from '@reduxjs/toolkit'
import { List } from '@/client/api/task'

type TaskaState = {
  isDrawerOpen: boolean
  mutate?: KeyedMutator<List>
}

const initialState: TaskaState = {
  isDrawerOpen: false,
}

const slice = createSlice({
  name: 'tasks',
  initialState,
  reducers: {
    mutateUpdated: (state, action: PayloadAction<KeyedMutator<List>>) => {
      state.mutate = action.payload
    },
    drawerDidOpen: (state) => {
      state.isDrawerOpen = true
    },
    drawerDidClose: (state) => {
      state.isDrawerOpen = false
    },
  },
})

export const { mutateUpdated, drawerDidOpen, drawerDidClose } = slice.actions

export default slice.reducer
