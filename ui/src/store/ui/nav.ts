// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { createSlice, PayloadAction } from '@reduxjs/toolkit'

type NavState = {
  active: string | null
}

export enum NavType {
  Account = 'account',
  Workspaces = 'workspaces',
  Groups = 'groups',
  Organizations = 'organizations',
  Console = 'console',
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
