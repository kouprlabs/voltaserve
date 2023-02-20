import { configureStore } from '@reduxjs/toolkit'
import reducer from './reducer'

const store = configureStore({
  reducer,
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({
      serializableCheck: {
        ignoredActions: [
          'uploads/uploadAdded',
          'uploads/uploadUpdated',
          'uploads/uploadDeleted',
        ],
        ignoredPaths: ['entities.uploads.items'],
      },
    }),
})

export type RootState = ReturnType<typeof store.getState>
export type AppDispatch = typeof store.dispatch

export default store
