// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { KeyedMutator } from 'swr'
import { PayloadAction, createSlice } from '@reduxjs/toolkit'
import { Info } from '@/client/api/mosaic'

type MosaicState = {
  isModalOpen: boolean
  mutateInfo?: KeyedMutator<Info>
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
    mutateInfoUpdated: (state, action: PayloadAction<KeyedMutator<Info>>) => {
      state.mutateInfo = action.payload
    },
  },
})

export const { modalDidOpen, modalDidClose, mutateInfoUpdated } = slice.actions

export default slice.reducer
