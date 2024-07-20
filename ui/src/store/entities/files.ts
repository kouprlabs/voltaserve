// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import { List } from '@/client/api/file'

type FilesState = {
  list?: List
}

const initialState: FilesState = {}

const slice = createSlice({
  name: 'files',
  initialState,
  reducers: {
    listUpdated: (state, action: PayloadAction<List>) => {
      state.list = action.payload
    },
    fileRemoved: (state, action: PayloadAction<string>) => {
      if (state.list) {
        state.list = {
          ...state.list,
          data: state.list.data.filter((item) => item.id !== action.payload),
        }
      }
    },
  },
})

export const { listUpdated, fileRemoved } = slice.actions

export default slice.reducer
