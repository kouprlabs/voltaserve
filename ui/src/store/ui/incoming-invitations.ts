import { KeyedMutator } from 'swr'
import { PayloadAction, createSlice } from '@reduxjs/toolkit'
import { List } from '@/client/api/invitation'

type OutgoingInvitationsState = {
  mutate?: KeyedMutator<List>
}

const initialState: OutgoingInvitationsState = {}

const slice = createSlice({
  name: 'incoming-invitations',
  initialState,
  reducers: {
    mutateUpdated: (state, action: PayloadAction<KeyedMutator<List>>) => {
      state.mutate = action.payload
    },
  },
})

export const { mutateUpdated } = slice.actions

export default slice.reducer
