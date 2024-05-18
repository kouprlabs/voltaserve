import { KeyedMutator } from 'swr'
import { PayloadAction, createSlice } from '@reduxjs/toolkit'
import { Workspace } from '@/client/api/workspace'

type AccountState = {
  mutate?: KeyedMutator<Workspace>
}

const initialState: AccountState = {}

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
