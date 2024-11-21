// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { createSlice } from '@reduxjs/toolkit'

export type UploadsState = {
  isDrawerOpen: boolean
}

const initialState: UploadsState = {
  isDrawerOpen: false,
}

const slice = createSlice({
  name: 'uploads',
  initialState,
  reducers: {
    drawerDidOpen: (state) => {
      state.isDrawerOpen = true
    },
    drawerDidClose: (state) => {
      state.isDrawerOpen = false
    },
  },
})

export const { drawerDidOpen, drawerDidClose } = slice.actions

export default slice.reducer
