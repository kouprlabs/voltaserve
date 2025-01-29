// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import { FileList } from '@/client/api/file'

type FilesState = {
  list?: FileList
}

const initialState: FilesState = {}

const slice = createSlice({
  name: 'files',
  initialState,
  reducers: {
    listUpdated: (state, action: PayloadAction<FileList>) => {
      state.list = action.payload
    },
  },
})

export const { listUpdated } = slice.actions

export default slice.reducer
