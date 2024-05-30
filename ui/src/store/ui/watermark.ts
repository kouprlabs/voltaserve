import { KeyedMutator } from 'swr'
import { PayloadAction, createSlice } from '@reduxjs/toolkit'
import { File } from '@/client/api/file'

type WatermarkState = {
  isCreating: boolean
  isUpdating: boolean
  isDeleting: boolean
  isModalOpen: boolean
  mutateFile?: KeyedMutator<File>
}

const initialState: WatermarkState = {
  isCreating: false,
  isUpdating: false,
  isDeleting: false,
  isModalOpen: false,
}

const slice = createSlice({
  name: 'watermark',
  initialState,
  reducers: {
    creatingDidStart: (state) => {
      state.isCreating = true
    },
    creatingDidStop: (state) => {
      state.isCreating = false
    },
    updatingDidStart: (state) => {
      state.isUpdating = true
    },
    updatingDidStop: (state) => {
      state.isUpdating = false
    },
    deletingDidStart: (state) => {
      state.isDeleting = true
    },
    deletingDidStop: (state) => {
      state.isDeleting = false
    },
    modalDidOpen: (state) => {
      state.isModalOpen = true
    },
    modalDidClose: (state) => {
      state.isModalOpen = false
    },
    mutateFileUpdated: (state, action: PayloadAction<KeyedMutator<File>>) => {
      state.mutateFile = action.payload
    },
  },
})

export const {
  modalDidOpen,
  modalDidClose,
  mutateFileUpdated,
  creatingDidStart,
  creatingDidStop,
  updatingDidStart,
  updatingDidStop,
  deletingDidStart,
  deletingDidStop,
} = slice.actions

export default slice.reducer
