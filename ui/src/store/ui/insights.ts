import { PayloadAction, createSlice } from '@reduxjs/toolkit'

type InsightsState = {
  isModalOpen: boolean
  isWizardComplete: boolean
}

const initialState: InsightsState = {
  isModalOpen: false,
  isWizardComplete: false,
}

const slice = createSlice({
  name: 'analysis',
  initialState,
  reducers: {
    wizardDidComplete: (state, action: PayloadAction<boolean>) => {
      state.isWizardComplete = action.payload
    },
    modalDidOpen: (state) => {
      state.isModalOpen = true
    },
    modalDidClose: (state) => {
      state.isModalOpen = false
    },
  },
})

export const { wizardDidComplete, modalDidOpen, modalDidClose } = slice.actions

export default slice.reducer
