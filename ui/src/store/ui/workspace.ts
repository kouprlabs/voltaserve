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
import { Workspace } from '@/client/api/workspace'

type WorkspaceState = {
  mutate?: KeyedMutator<Workspace>
}

const initialState: WorkspaceState = {}

const slice = createSlice({
  name: 'workspace',
  initialState,
  reducers: {
    mutateUpdated: (state, action: PayloadAction<KeyedMutator<Workspace>>) => {
      state.mutate = action.payload
    },
  },
})

export const { mutateUpdated } = slice.actions

export default slice.reducer
