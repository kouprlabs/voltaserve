import React from 'react'
import ReactDOM from 'react-dom/client'
import { RouterProvider } from 'react-router-dom'
import { Provider } from 'react-redux'
import { Theme } from '@radix-ui/themes'
import { HelmetProvider } from 'react-helmet-async'
import store from '@/store/configure-store'
import '@/styles/index.css'
import router from './router'

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
    <Provider store={store}>
      <Theme>
        <HelmetProvider>
          <RouterProvider router={router} />
        </HelmetProvider>
      </Theme>
    </Provider>
  </React.StrictMode>,
)
