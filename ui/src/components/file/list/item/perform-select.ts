import { File } from '@/client/api/file'
import store from '@/store/configure-store'
import {
  selectedItemAdded,
  selectedItemRemoved,
  selectedItemsUpdated,
} from '@/store/ui/files'

export function performRangeSelect(file: File) {
  const ui = store.getState().ui.files
  const data = store.getState().entities.files.list?.data
  if (data && ui.isRangeSelectActive && ui.selectedItems.length >= 1) {
    let startIndex = data.findIndex((e) => e.id === ui.selectedItems[0])
    let endIndex = data.findIndex((e) => e.id === file.id)
    if (startIndex !== -1 && endIndex !== -1 && startIndex !== endIndex) {
      if (startIndex > endIndex) {
        ;[startIndex, endIndex] = [endIndex, startIndex]
      }
      const ids = []
      for (let i = startIndex; i <= endIndex; i++) {
        ids.push(data[i].id)
      }
      store.dispatch(selectedItemsUpdated(ids))
    }
  }
}

export function performMultiSelect(file: File, isSelected: boolean) {
  if (isSelected) {
    store.dispatch(selectedItemRemoved(file.id))
  } else {
    store.dispatch(selectedItemAdded(file.id))
  }
}
