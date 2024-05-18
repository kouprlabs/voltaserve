import { KeyedMutator } from 'swr'
import { PayloadAction, createSlice } from '@reduxjs/toolkit'
import { Group } from '@/client/api/group'

type AccountState = {
  mutate?: KeyedMutator<Group>
}

const initialState: AccountState = {}

const slice = createSlice({
  name: 'group',
  initialState,
  reducers: {
    mutateUpdated: (state, action: PayloadAction<KeyedMutator<Group>>) => {
      state.mutate = action.payload
    },
  },
})

export const { mutateUpdated } = slice.actions

export default slice.reducer
