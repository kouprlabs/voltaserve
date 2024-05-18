import { KeyedMutator } from 'swr'
import { PayloadAction, createSlice } from '@reduxjs/toolkit'
import { List } from '@/client/api/group'

type AccountState = {
  mutate?: KeyedMutator<List>
}

const initialState: AccountState = {}

const slice = createSlice({
  name: 'groups',
  initialState,
  reducers: {
    mutateUpdated: (state, action: PayloadAction<KeyedMutator<List>>) => {
      state.mutate = action.payload
    },
  },
})

export const { mutateUpdated } = slice.actions

export default slice.reducer
