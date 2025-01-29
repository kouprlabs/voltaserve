// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { KeyedMutator } from 'swr'
import { PayloadAction, createSlice } from '@reduxjs/toolkit'
import { MosaicInfo } from '@/client/api/mosaic'

type MosaicState = {
  isModalOpen: boolean
  mutateInfo?: KeyedMutator<MosaicInfo>
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
    mutateInfoUpdated: (
      state,
      action: PayloadAction<KeyedMutator<MosaicInfo>>,
    ) => {
      state.mutateInfo = action.payload
    },
  },
})

export const {
  modalDidOpen,
  modalDidClose,
  allModalsDidClose,
  mutateInfoUpdated,
} = slice.actions

export default slice.reducer
