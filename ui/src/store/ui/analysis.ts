import { PayloadAction, createSlice } from '@reduxjs/toolkit'

type AnalysisState = {
  isModalOpen: boolean
  isWizardComplete: boolean
}

const initialState: AnalysisState = {
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
