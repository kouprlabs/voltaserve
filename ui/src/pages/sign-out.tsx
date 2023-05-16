import { useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { Heading } from '@chakra-ui/react'
import { Helmet } from 'react-helmet-async'
import FullLayout from '@/components/layout/full'
import { clearToken } from '@/infra/token'

function SignOutPage() {
  const navigate = useNavigate()

  useEffect(() => {
    clearToken()
    navigate('/sign-in')
  }, [navigate])

  return (
    <FullLayout>
      <>
        <Helmet>
          <title>Signing Out…</title>
        </Helmet>
        <Heading fontSize="16px">Signing out…</Heading>
      </>
    </FullLayout>
  )
}

export default SignOutPage
