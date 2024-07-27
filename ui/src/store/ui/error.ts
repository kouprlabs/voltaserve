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

type ErrorState = {
  value: string | null
}

const initialState: ErrorState = {
  value: null,
}

const slice = createSlice({
  name: 'error',
  initialState,
  reducers: {
    errorOccurred: (state, action: PayloadAction<string>) => {
      state.value = action.payload
    },
    errorCleared: (state) => {
      state.value = null
    },
  },
})

export const { errorOccurred, errorCleared } = slice.actions

export default slice.reducer
