// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { useCallback, useEffect } from 'react'
import { Outlet, useLocation, useNavigate } from 'react-router-dom'
import { useColorMode } from '@chakra-ui/react'
import { Helmet } from 'react-helmet-async'
import { useAppDispatch } from '@/store/hook'
import { allModalsDidClose as allModalsDidCloseForFiles } from '@/store/ui/files'
import { allModalsDidClose as allModalsDidCloseForInsights } from '@/store/ui/insights'
import { allModalsDidClose as allModalsDidCloseForMosaic } from '@/store/ui/mosaic'
import { allModalsDidClose as allModalsDidCloseForOrganizations } from '@/store/ui/organizations'
import { allModalsDidClose as allModalsDidCloseForSearchFilter } from '@/store/ui/search-filter'
import { allModalsDidClose as allModalsDidCloseForSnapshots } from '@/store/ui/snapshots'

const RootPage = () => {
  const location = useLocation()
  const navigate = useNavigate()
  const dispatch = useAppDispatch()
  const { colorMode } = useColorMode()

  useEffect(() => {
    closeAllModals()
    if (location.pathname === '/') {
      navigate('/workspace')
    }
  }, [location.pathname, navigate])

  useEffect(() => {
    const element = document.querySelector("link[rel='icon']")
    if (element) {
      window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', (event: MediaQueryListEvent) => {
        if (event.matches) {
          element.setAttribute('href', '/favicon-dark.svg')
        } else {
          element.setAttribute('href', '/favicon.svg')
        }
      })
      if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
        element.setAttribute('href', '/favicon-dark.svg')
      } else {
        element.setAttribute('href', '/favicon.svg')
      }
    }
  }, [])

  useEffect(() => {
    const body = document.getElementsByTagName('body')[0]
    if (colorMode === 'dark') {
      body.classList.add('dark')
      body.classList.add('koupr-dark')
    } else {
      body.classList.remove('dark')
      body.classList.remove('koupr-dark')
    }
  }, [colorMode])

  const closeAllModals = useCallback(() => {
    dispatch(allModalsDidCloseForFiles())
    dispatch(allModalsDidCloseForInsights())
    dispatch(allModalsDidCloseForMosaic())
    dispatch(allModalsDidCloseForSnapshots())
    dispatch(allModalsDidCloseForOrganizations())
    dispatch(allModalsDidCloseForSearchFilter())
  }, [dispatch])

  return (
    <>
      <Helmet>
        <title>Voltaserve</title>
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link href="/favicon.svg" rel="icon" type="image/svg+xml" />
      </Helmet>
      <Outlet />
    </>
  )
}

export default RootPage
