// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

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
