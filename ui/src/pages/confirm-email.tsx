import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import {
  CircularProgress,
  Link as ChakraLink,
  Text,
  useToast,
  VStack,
} from '@chakra-ui/react'
import { variables } from '@koupr/ui'
import { Helmet } from 'react-helmet-async'
import AccountAPI from '@/api/account'
import Logo from '@/components/common/logo'
import FullLayout from '@/components/layout/full'

const ConfirmEmailPage = () => {
  const params = useParams()
  const token = params.token as string
  const toast = useToast()
  const [isSuccessful, setIsSuccessful] = useState<boolean | null>(null)

  useEffect(() => {
    async function doRequest() {
      try {
        await AccountAPI.confirmEmail({ token: token })
        setIsSuccessful(true)
      } catch {
        setIsSuccessful(false)
      }
    }
    if (token) {
      doRequest()
    }
  }, [token, toast])

  return (
    <FullLayout>
      <>
        <Helmet>
          <title>Confirm Email</title>
        </Helmet>
        <VStack spacing="25px" w="100%">
          <VStack spacing={variables.spacing}>
            <Logo className="w-16" isGlossy={true} />
            {isSuccessful === null && (
              <>
                <h1 className="font-display text-2xl font-medium text-center">
                  Confirming your Email…
                </h1>
                <CircularProgress isIndeterminate />
              </>
            )}
            {isSuccessful && (
              <>
                <h1 className="font-display text-2xl font-medium text-center">
                  Email confirmed
                </h1>
                <Text>Click the link below to sign in.</Text>
                <ChakraLink as={Link} to="/sign-in">
                  Sign in
                </ChakraLink>
              </>
            )}
            {isSuccessful === false && (
              <h1 className="font-display text-2xl font-medium text-center">
                An error occurred while processing your request.
              </h1>
            )}
          </VStack>
        </VStack>
      </>
    </FullLayout>
  )
}

export default ConfirmEmailPage
