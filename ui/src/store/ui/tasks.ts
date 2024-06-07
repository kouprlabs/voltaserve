import { KeyedMutator } from 'swr'
import { PayloadAction, createSlice } from '@reduxjs/toolkit'
import { List } from '@/client/api/task'

type TaskaState = {
  isDrawerOpen: boolean
  mutateList?: KeyedMutator<List>
  mutateCount?: KeyedMutator<number>
}

const initialState: TaskaState = {
  isDrawerOpen: false,
}

const slice = createSlice({
  name: 'tasks',
  initialState,
  reducers: {
    drawerDidOpen: (state) => {
      state.isDrawerOpen = true
    },
    drawerDidClose: (state) => {
      state.isDrawerOpen = false
    },
    mutateListUpdated: (state, action: PayloadAction<KeyedMutator<List>>) => {
      state.mutateList = action.payload
    },
    mutateCountUpdated: (
      state,
      action: PayloadAction<KeyedMutator<number>>,
    ) => {
      state.mutateCount = action.payload
    },
  },
})

export const {
  drawerDidOpen,
  drawerDidClose,
  mutateListUpdated,
  mutateCountUpdated,
} = slice.actions

export default slice.reducer
