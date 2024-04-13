import { useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { Heading } from '@chakra-ui/react'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import LayoutFull from '@/components/layout/layout-full'
import { clearToken } from '@/infra/token'

function SignOutPage() {
  const navigate = useNavigate()

  useEffect(() => {
    clearToken()
    navigate('/sign-in')
  }, [navigate])

  return (
    <LayoutFull>
      <>
        <Helmet>
          <title>Signing Out…</title>
        </Helmet>
        <Heading className={cx('text-heading')}>Signing out…</Heading>
      </>
    </LayoutFull>
  )
}

export default SignOutPage
