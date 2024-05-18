import { KeyedMutator } from 'swr'
import { PayloadAction, createSlice } from '@reduxjs/toolkit'
import { Organization } from '@/client/api/organization'

type AccountState = {
  mutate?: KeyedMutator<Organization>
}

const initialState: AccountState = {}

const slice = createSlice({
  name: 'organization',
  initialState,
  reducers: {
    mutateUpdated: (
      state,
      action: PayloadAction<KeyedMutator<Organization>>,
    ) => {
      state.mutate = action.payload
    },
  },
})

export const { mutateUpdated } = slice.actions

export default slice.reducer
