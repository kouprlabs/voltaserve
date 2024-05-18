import { KeyedMutator } from 'swr'
import { PayloadAction, createSlice } from '@reduxjs/toolkit'
import { List } from '@/client/api/organization'

type OrganizationsState = {
  isInviteModalOpen: boolean
  mutate?: KeyedMutator<List>
}

const initialState: OrganizationsState = {
  isInviteModalOpen: false,
}

const slice = createSlice({
  name: 'organizations',
  initialState,
  reducers: {
    inviteModalDidOpen: (state) => {
      state.isInviteModalOpen = true
    },
    inviteModalDidClose: (state) => {
      state.isInviteModalOpen = false
    },
    mutateUpdated: (state, action: PayloadAction<KeyedMutator<List>>) => {
      state.mutate = action.payload
    },
  },
})

export const { inviteModalDidOpen, inviteModalDidClose, mutateUpdated } =
  slice.actions

export default slice.reducer
