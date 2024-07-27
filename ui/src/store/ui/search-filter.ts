// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { createSlice } from '@reduxjs/toolkit'

type SearchFilterState = {
  isModalOpen: boolean
}

const initialState: SearchFilterState = {
  isModalOpen: false,
}

const slice = createSlice({
  name: 'search-filter',
  initialState,
  reducers: {
    modalDidOpen: (state) => {
      state.isModalOpen = true
    },
    modalDidClose: (state) => {
      state.isModalOpen = false
    },
  },
})

export const { modalDidOpen, modalDidClose } = slice.actions

export default slice.reducer
