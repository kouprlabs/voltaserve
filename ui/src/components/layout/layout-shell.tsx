// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { useEffect } from 'react'
import { Outlet, useLocation, useNavigate } from 'react-router-dom'
import { useToast } from '@chakra-ui/react'
import { IconFlag, IconGroup, IconWorkspaces, Logo, Shell } from '@koupr/ui'
import { cx } from '@emotion/css'
import AppBar from '@/components/app-bar'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { errorCleared } from '@/store/ui/error'

const LayoutShell = () => {
  const location = useLocation()
  const navigate = useNavigate()
  const toast = useToast()
  const error = useAppSelector((state) => state.ui.error.value)
  const dispatch = useAppDispatch()

  useEffect(() => {
    if (error) {
      toast({
        title: error,
        status: 'error',
        isClosable: true,
      })
      dispatch(errorCleared())
    }
  }, [error, toast, dispatch])

  return (
    <Shell
      storage={{ prefix: 'voltaserve', namespace: 'main' }}
      logo={
        <div className={cx('w-[16px]')}>
          <Logo type="voltaserve" size="sm" />
        </div>
      }
      topBar={<AppBar />}
      items={[
        {
          href: '/workspace',
          icon: <IconWorkspaces />,
          primaryText: 'Workspaces',
          secondaryText: 'Isolated containers for files and folders.',
        },
        {
          href: '/group',
          icon: <IconGroup />,
          primaryText: 'Groups',
          secondaryText: 'Allows assigning permissions to a group of users.',
        },
        {
          href: '/organization',
          icon: <IconFlag />,
          primaryText: 'Organizations',
          secondaryText: 'Umbrellas for workspaces and users.',
        },
      ]}
      navigateFn={navigate}
      pathnameFn={() => location.pathname}
    >
      <Outlet />
    </Shell>
  )
}

export default LayoutShell
