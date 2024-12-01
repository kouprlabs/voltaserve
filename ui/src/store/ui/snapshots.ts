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
import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import { List } from '@/client/api/snapshot'

export type SnapshotsState = {
  selection: string[]
  isListModalOpen: boolean
  isDetachModalOpen: boolean
  snapshotMutate?: KeyedMutator<List | undefined>
}

const initialState: SnapshotsState = {
  selection: [],
  isListModalOpen: false,
  isDetachModalOpen: false,
}

const slice = createSlice({
  name: 'snapshots',
  initialState,
  reducers: {
    selectionUpdated: (state, action: PayloadAction<string[]>) => {
      state.selection = action.payload
    },
    mutateUpdated: (state, action: PayloadAction<KeyedMutator<List | undefined>>) => {
      state.snapshotMutate = action.payload
    },
    listModalDidOpen: (state) => {
      state.isListModalOpen = true
    },
    detachModalDidOpen: (state) => {
      state.isDetachModalOpen = true
    },
    listModalDidClose: (state) => {
      state.isListModalOpen = false
    },
    detachModalDidClose: (state) => {
      state.isDetachModalOpen = false
    },
    allModalsDidClose: (state) => {
      state.isListModalOpen = false
      state.isDetachModalOpen = false
    },
  },
})

export const {
  selectionUpdated,
  mutateUpdated,
  listModalDidOpen,
  detachModalDidOpen,
  listModalDidClose,
  detachModalDidClose,
  allModalsDidClose,
} = slice.actions

export default slice.reducer
