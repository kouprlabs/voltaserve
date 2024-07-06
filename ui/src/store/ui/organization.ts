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
import { Organization } from '@/client/api/organization'

type AccountState = {
  mutate?: KeyedMutator<Organization>
}

const initialState: AccountState = {}

const slice = createSlice({
  name: 'organization',
  initialState,
  reducers: {
    mutateUpdated: (
      state,
      action: PayloadAction<KeyedMutator<Organization>>,
    ) => {
      state.mutate = action.payload
    },
  },
})

export const { mutateUpdated } = slice.actions

export default slice.reducer
