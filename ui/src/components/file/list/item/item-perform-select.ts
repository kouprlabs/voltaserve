// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { File } from '@/client/api/file'
import store from '@/store/configure-store'
import {
  selectionAdded,
  selectionRemoved,
  selectionUpdated,
} from '@/store/ui/files'

export function performRangeSelect(file: File) {
  const ui = store.getState().ui.files
  const data = store.getState().entities.files.list?.data || []
  if (data && ui.isRangeSelectActive && ui.selection.length >= 1) {
    let startIndex = data.findIndex((e) => e.id === ui.selection[0])
    let endIndex = data.findIndex((e) => e.id === file.id)
    if (startIndex !== -1 && endIndex !== -1 && startIndex !== endIndex) {
      if (startIndex > endIndex) {
        ;[startIndex, endIndex] = [endIndex, startIndex]
      }
      const ids = []
      for (let i = startIndex; i <= endIndex; i++) {
        ids.push(data[i].id)
      }
      store.dispatch(selectionUpdated(ids))
    }
  }
}

export function performMultiSelect(file: File, isSelected: boolean) {
  if (isSelected) {
    store.dispatch(selectionRemoved(file.id))
  } else {
    store.dispatch(selectionAdded(file.id))
  }
}
