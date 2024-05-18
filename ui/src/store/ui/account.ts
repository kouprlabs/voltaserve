import { KeyedMutator } from 'swr'
import { PayloadAction, createSlice } from '@reduxjs/toolkit'
import { User } from '@/client/api/user'

type AccountState = {
  mutate?: KeyedMutator<User>
}

const initialState: AccountState = {}

const slice = createSlice({
  name: 'account',
  initialState,
  reducers: {
    mutateUpdated: (state, action: PayloadAction<KeyedMutator<User>>) => {
      state.mutate = action.payload
    },
  },
})

export const { mutateUpdated } = slice.actions

export default slice.reducer
