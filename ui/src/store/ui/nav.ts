import { createSlice, PayloadAction } from '@reduxjs/toolkit'

type NavState = {
  active: string | null
}

export enum NavType {
  Account = 'account',
  Workspaces = 'workspaces',
  Groups = 'groups',
  Organizations = 'organizations',
}

const initialState: NavState = {
  active: null,
}

const slice = createSlice({
  name: 'nav',
  initialState,
  reducers: {
    activeNavChanged: (state, action: PayloadAction<NavType>) => {
      state.active = action.payload
    },
  },
})

export const { activeNavChanged } = slice.actions

export default slice.reducer
