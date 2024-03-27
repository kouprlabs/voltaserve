import { createSlice } from '@reduxjs/toolkit'

type OrganizationsState = {
  isInviteModalOpen: boolean
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
  },
})

export const { inviteModalDidOpen, inviteModalDidClose } = slice.actions

export default slice.reducer
