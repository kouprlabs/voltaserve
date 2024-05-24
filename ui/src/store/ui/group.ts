import { KeyedMutator } from 'swr'
import { PayloadAction, createSlice } from '@reduxjs/toolkit'
import { Group } from '@/client/api/group'

type GroupState = {
  mutate?: KeyedMutator<Group>
}

const initialState: GroupState = {}

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
