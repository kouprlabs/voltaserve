import { KeyedMutator } from 'swr'
import { PayloadAction, createSlice } from '@reduxjs/toolkit'
import { Workspace } from '@/client/api/workspace'

type WorkspaceState = {
  mutate?: KeyedMutator<Workspace>
}

const initialState: WorkspaceState = {}

const slice = createSlice({
  name: 'workspace',
  initialState,
  reducers: {
    mutateUpdated: (state, action: PayloadAction<KeyedMutator<Workspace>>) => {
      state.mutate = action.payload
    },
  },
})

export const { mutateUpdated } = slice.actions

export default slice.reducer
