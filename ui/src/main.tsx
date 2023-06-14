import React from 'react'
import ReactDOM from 'react-dom/client'
import { RouterProvider } from 'react-router-dom'
import { Provider } from 'react-redux'
import { ChakraProvider } from '@chakra-ui/react'
import { theme } from '@koupr/ui'
import '@koupr/ui/styles/index.css'
import { HelmetProvider } from 'react-helmet-async'
import store from '@/store/configure-store'
import '@/styles/index.css'
import router from './router'

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
    <Provider store={store}>
      <ChakraProvider theme={theme}>
        <HelmetProvider>
          <RouterProvider router={router} />
        </HelmetProvider>
      </ChakraProvider>
    </Provider>
  </React.StrictMode>
)
