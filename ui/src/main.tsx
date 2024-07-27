// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import React from 'react'
import ReactDOM from 'react-dom/client'
import { RouterProvider } from 'react-router-dom'
import { Provider } from 'react-redux'
import { ChakraProvider } from '@chakra-ui/react'
import { HelmetProvider } from 'react-helmet-async'
import theme from '@/lib/theme'
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
  </React.StrictMode>,
)
