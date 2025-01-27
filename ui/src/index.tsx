// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { RouterProvider } from 'react-router-dom'
import { Provider } from 'react-redux'
import { ChakraProvider } from '@chakra-ui/react'
import { theme } from '@koupr/ui'
import { HelmetProvider } from 'react-helmet-async'
import store from '@/store/configure-store'
import router from './router'
import './styles.css'

const Voltaserve = () => (
  <Provider store={store}>
    <ChakraProvider theme={theme}>
      <HelmetProvider>
        <RouterProvider router={router} />
      </HelmetProvider>
    </ChakraProvider>
  </Provider>
)

export default Voltaserve
