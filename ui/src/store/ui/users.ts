// Copyright 2024 Mateusz Kaźmierczak.
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
import { UserManagementList } from '@/client/admin/admin'

type UsersState = {
  mutate?: KeyedMutator<UserManagementList>
}

const initialState: UsersState = {}

const slice = createSlice({
  name: 'indexes',
  initialState,
  reducers: {
    mutateUpdated: (
      state,
      action: PayloadAction<KeyedMutator<UserManagementList>>,
    ) => {
      state.mutate = action.payload
    },
  },
})

export const { mutateUpdated } = slice.actions

export default slice.reducer
