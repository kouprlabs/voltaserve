// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { Link as ChakraLink, Heading } from '@chakra-ui/react'
import { Logo, Spinner } from '@koupr/ui'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import { AccountAPI } from '@/client/idp/account'
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
    if (token) {
      ;(async function (token: string) {
        try {
          await AccountAPI.confirmEmail({ token: token })
          setIsCompleted(true)
        } catch {
          setIsFailed(true)
        } finally {
          setIsCompleted(true)
        }
      })(token)
    }
  }, [token])

  return (
    <LayoutFull>
      <Helmet>
        <title>Confirm Email</title>
      </Helmet>
      <div className={cx('flex', 'flex-col', 'items-center', 'gap-3')}>
        <div className={cx('w-[64px]')}>
          <Logo type="voltaserve" size="md" isGlossy={true} />
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
              <span>Click the link below to sign in.</span>
              <ChakraLink as={Link} to="/sign-in">
                Sign In
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

export default ConfirmEmailPage
