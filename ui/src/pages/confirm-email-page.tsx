import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { Link as ChakraLink, Heading, Text } from '@chakra-ui/react'
import { Spinner } from '@koupr/ui'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import AccountAPI from '@/client/idp/account'
import Logo from '@/components/common/logo'
import LayoutFull from '@/components/layout/layout-full'

const ConfirmEmailPage = () => {
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
        await AccountAPI.confirmEmail({ token: token })
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
        <div className={cx('w-[64px]')}>
          <Logo isGlossy={true} />
        </div>
        {!isCompleted && !isFailed ? (
          <div className={cx('flex', 'flex-col', 'items-center', 'gap-1.5')}>
            <Heading className={cx('text-heading')}>
              Confirming your Email…
            </Heading>
            <Spinner />
          </div>
        ) : null}
        {isCompleted && !isFailed ? (
          <div className={cx('flex', 'flex-col', 'items-center', 'gap-1.5')}>
            <Heading className={cx('text-heading')}>Email confirmed</Heading>
            <div className={cx('flex', 'flex-col', 'items-center', 'gap-0.5')}>
              <Text>Click the link below to sign in.</Text>
              <ChakraLink as={Link} to="/sign-in">
                Sign In
              </ChakraLink>
            </div>
          </div>
        ) : null}
        {isFailed && (
          <Heading className={cx('text-heading')}>
            An error occurred while processing your request.
          </Heading>
        )}
      </div>
    </LayoutFull>
  )
}

export default ConfirmEmailPage
