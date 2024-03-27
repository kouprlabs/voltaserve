import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { Link as ChakraLink, Heading, Text } from '@chakra-ui/react'
import { variables, Spinner } from '@koupr/ui'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import UserAPI from '@/client/idp/user'
import Logo from '@/components/common/logo'
import LayoutFull from '@/components/layout/layout-full'

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
    <LayoutFull>
      <Helmet>
        <title>Confirm Email</title>
      </Helmet>
      <div className={cx('flex', 'flex-col', 'items-center', 'gap-3')}>
        <div className="w-16">
          <Logo isGlossy={true} />
        </div>
        {!isCompleted && !isFailed ? (
          <div className={cx('flex', 'flex-col', 'items-center', 'gap-1.5')}>
            <Heading fontSize={variables.headingFontSize}>
              Confirming your Email…
            </Heading>
            <Spinner />
          </div>
        ) : null}
        {isCompleted && !isFailed ? (
          <div className={cx('flex', 'flex-col', 'items-center', 'gap-1.5')}>
            <Heading fontSize={variables.headingFontSize}>
              Email confirmed
            </Heading>
            <div className={cx('flex', 'flex-row', 'items-center', 'gap-0.5')}>
              <Text>Click the link below to go back to your account.</Text>
              <ChakraLink as={Link} to="/account/settings">
                Back to account
              </ChakraLink>
            </div>
          </div>
        ) : null}
        {isFailed && (
          <Heading fontSize={variables.headingFontSize}>
            An error occurred while processing your request.
          </Heading>
        )}
      </div>
    </LayoutFull>
  )
}

export default UpdateEmailPage
