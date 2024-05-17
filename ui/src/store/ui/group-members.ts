import { KeyedMutator } from 'swr'
import { PayloadAction, createSlice } from '@reduxjs/toolkit'
import { List } from '@/client/api/user'

type GroupMembersState = {
  mutate?: KeyedMutator<List>
}

const initialState: GroupMembersState = {}

const slice = createSlice({
  name: 'group-members',
  initialState,
  reducers: {
    mutateUpdated: (state, action: PayloadAction<KeyedMutator<List>>) => {
      state.mutate = action.payload
    },
  },
})

export const { mutateUpdated } = slice.actions

export default slice.reducer
