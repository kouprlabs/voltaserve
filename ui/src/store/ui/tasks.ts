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
import { TaskList } from '@/client/api/task'

type TasksState = {
  isDrawerOpen: boolean
  mutateList?: KeyedMutator<TaskList>
  mutateCount?: KeyedMutator<number>
}

const initialState: TasksState = {
  isDrawerOpen: false,
}

const slice = createSlice({
  name: 'tasks',
  initialState,
  reducers: {
    drawerDidOpen: (state) => {
      state.isDrawerOpen = true
    },
    drawerDidClose: (state) => {
      state.isDrawerOpen = false
    },
    mutateListUpdated: (
      state,
      action: PayloadAction<KeyedMutator<TaskList>>,
    ) => {
      state.mutateList = action.payload
    },
    mutateCountUpdated: (
      state,
      action: PayloadAction<KeyedMutator<number>>,
    ) => {
      state.mutateCount = action.payload
    },
  },
})

export const {
  drawerDidOpen,
  drawerDidClose,
  mutateListUpdated,
  mutateCountUpdated,
} = slice.actions

export default slice.reducer
