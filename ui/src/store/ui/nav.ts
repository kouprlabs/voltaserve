// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
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
