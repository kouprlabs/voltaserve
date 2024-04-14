import { useEffect } from 'react'
import { Outlet, useLocation, useNavigate } from 'react-router-dom'
import { useColorMode } from '@chakra-ui/react'
import { cx } from '@emotion/css'
import { Helmet } from 'react-helmet-async'

const RootPage = () => {
  const location = useLocation()
  const navigate = useNavigate()
  const { colorMode } = useColorMode()

  useEffect(() => {
    if (location.pathname === '/') {
      navigate('/workspace')
    }
  }, [location.pathname, navigate])

  useEffect(() => {
    const element = document.querySelector("link[rel='icon']")
    if (element) {
      window
        .matchMedia('(prefers-color-scheme: dark)')
        .addEventListener('change', (event: MediaQueryListEvent) => {
          if (event.matches) {
            element.setAttribute('href', '/favicon-dark.svg')
          } else {
            element.setAttribute('href', '/favicon.svg')
          }
        })
      if (
        window.matchMedia &&
        window.matchMedia('(prefers-color-scheme: dark)').matches
      ) {
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
    } else {
      body.classList.remove('dark')
    }
  }, [colorMode])

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
