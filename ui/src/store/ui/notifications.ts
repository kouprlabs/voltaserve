import { KeyedMutator } from 'swr'
import { PayloadAction, createSlice } from '@reduxjs/toolkit'
import { Notification } from '@/client/api/notification'

type AccountState = {
  mutate?: KeyedMutator<Notification[]>
}

const initialState: AccountState = {}

const slice = createSlice({
  name: 'notifications',
  initialState,
  reducers: {
    mutateUpdated: (
      state,
      action: PayloadAction<KeyedMutator<Notification[]>>,
    ) => {
      state.mutate = action.payload
    },
  },
})

export const { mutateUpdated } = slice.actions

export default slice.reducer
