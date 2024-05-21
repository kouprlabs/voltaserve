import { KeyedMutator } from 'swr'
import { PayloadAction, createSlice } from '@reduxjs/toolkit'
import { List } from '@/client/api/workspace'

type WorkspacesState = {
  mutate?: KeyedMutator<List>
}

const initialState: WorkspacesState = {}

const slice = createSlice({
  name: 'workspaces',
  initialState,
  reducers: {
    mutateUpdated: (state, action: PayloadAction<KeyedMutator<List>>) => {
      state.mutate = action.payload
    },
  },
})

export const { mutateUpdated } = slice.actions

export default slice.reducer
