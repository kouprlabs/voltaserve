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
import { Group } from '@/client/api/group'

type GroupState = {
  mutate?: KeyedMutator<Group>
}

const initialState: GroupState = {}

const slice = createSlice({
  name: 'group',
  initialState,
  reducers: {
    mutateUpdated: (state, action: PayloadAction<KeyedMutator<Group>>) => {
      state.mutate = action.payload
    },
  },
})

export const { mutateUpdated } = slice.actions

export default slice.reducer
