import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { Link as ChakraLink, Text, VStack } from '@chakra-ui/react'
import { variables, Spinner } from '@koupr/ui'
import { Helmet } from 'react-helmet-async'
import UserAPI from '@/client/idp/user'
import Logo from '@/components/common/logo'
import FullLayout from '@/components/layout/full'

const UpdateEmailPage = () => {
  const params = useParams()
  const [isCompleted, setIsCompleted] = useState(false)
  const [isFailed, setIsFailed] = useState(false)
  const [token, setToken] = useState<string>('')

  useEffect(() => {
    setToken(params.token as string)
  }, [params.token])

  useEffect(() => {
    async function doRequest() {
      try {
        await UserAPI.updateEmailConfirmation({ token: token })
        setIsCompleted(true)
      } catch {
        setIsFailed(true)
      } finally {
        setIsCompleted(true)
      }
    }
    if (token) {
      doRequest()
    }
  }, [token])

  return (
    <FullLayout>
      <>
        <Helmet>
          <title>Confirm Email</title>
        </Helmet>
        <VStack spacing={variables.spacingXl}>
          <Logo className="w-16" isGlossy={true} />
          {!isCompleted && !isFailed ? (
            <VStack spacing={variables.spacing}>
              <h1 className="font-display text-2xl font-medium text-center">
                Confirming your Email…
              </h1>
              <Spinner />
            </VStack>
          ) : null}
          {isCompleted && !isFailed ? (
            <VStack spacing={variables.spacing}>
              <h1 className="font-display text-2xl font-medium text-center">
                Email confirmed
              </h1>
              <VStack spacing={variables.spacingXs}>
                <Text>Click the link below to go back to your account.</Text>
                <ChakraLink as={Link} to="/account/settings">
                  Back to account
                </ChakraLink>
              </VStack>
            </VStack>
          ) : null}
          {isFailed && (
            <h1 className="font-display text-2xl font-medium text-center">
              An error occurred while processing your request.
            </h1>
          )}
        </VStack>
      </>
    </FullLayout>
  )
}

export default UpdateEmailPage
