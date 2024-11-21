// Copyright (c) 2024 Mateusz Kaźmierczak.
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
import {
  IconWorkspaces,
  IconHome,
  IconGroup,
  IconFlag,
  IconPerson,
  Shell,
} from '@koupr/ui'
import { Logo } from '@koupr/ui'
import { cx } from '@emotion/css'
import AppBar from '@/components/app-bar'
import { getAdminStatus } from '@/infra/token'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { errorCleared } from '@/store/ui/error'

const LayoutConsole = () => {
  const toast = useToast()
  const error = useAppSelector((state) => state.ui.error.value)
  const navigate = useNavigate()
  const location = useLocation()
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

  useEffect(() => {
    if (!getAdminStatus()) {
      navigate('/')
    }
  }, [])

  return (
    <Shell
      storage={{ prefix: 'voltaserve', namespace: 'main' }}
      logo={
        <div className={cx('w-[16px]')}>
          <Logo type="voltaserve" />
        </div>
      }
      homeHref={
        location.pathname.startsWith('/console') ? '/console/dashboard' : '/'
      }
      topBar={<AppBar />}
      items={[
        {
          href: '/console/dashboard',
          icon: <IconHome />,
          primaryText: 'Overview',
          secondaryText: 'Basic information about your instance',
        },
        {
          href: '/console/workspaces',
          icon: <IconWorkspaces />,
          primaryText: 'Workspaces',
          secondaryText: 'Manage workspaces',
        },
        {
          href: '/console/groups',
          icon: <IconGroup />,
          primaryText: 'Groups',
          secondaryText: 'Manage groups',
        },
        {
          href: '/console/organizations',
          icon: <IconFlag />,
          primaryText: 'Organizations',
          secondaryText: 'Manage workspaces',
        },
        {
          href: '/console/users',
          icon: <IconPerson />,
          primaryText: 'Users',
          secondaryText: 'Manage users',
        },
      ]}
      navigateFn={navigate}
      pathnameFn={() => location.pathname}
    >
      <Outlet />
    </Shell>
  )
}

export default LayoutConsole
