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

type MosaicState = {
  isModalOpen: boolean
}

const initialState: MosaicState = {
  isModalOpen: false,
}

const slice = createSlice({
  name: 'mosaic',
  initialState,
  reducers: {
    modalDidOpen: (state) => {
      state.isModalOpen = true
    },
    modalDidClose: (state) => {
      state.isModalOpen = false
    },
    allModalsDidClose: (state) => {
      state.isModalOpen = false
    },
  },
})

export const { modalDidOpen, modalDidClose, allModalsDidClose } = slice.actions

export default slice.reducer
