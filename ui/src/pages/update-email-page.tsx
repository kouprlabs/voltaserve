// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { Link as ChakraLink, Heading } from '@chakra-ui/react'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import UserAPI from '@/client/idp/user'
import Logo from '@/components/common/logo'
import LayoutFull from '@/components/layout/layout-full'
import Spinner from '@/lib/components/spinner'

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
            <div className={cx('flex', 'flex-row', 'items-center', 'gap-0.5')}>
              <span>Click the link below to go back to your account.</span>
              <ChakraLink as={Link} to="/account/settings">
                Back to account
              </ChakraLink>
            </div>
          </div>
        ) : null}
        {isFailed ? (
          <Heading className={cx('text-heading')}>
            An error occurred while processing your request.
          </Heading>
        ) : null}
      </div>
    </LayoutFull>
  )
}

export default UpdateEmailPage
