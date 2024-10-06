// Copyright 2024 Mateusz Kaźmierczak.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { useEffect } from 'react'
import { Outlet, useNavigate } from 'react-router-dom'
import { useToast } from '@chakra-ui/react'
import { cx } from '@emotion/css'
import AppBar from '@/components/app-bar'
import Logo from '@/components/common/logo'
import { getAdminStatus } from '@/infra/token'
import {
  IconWorkspaces,
  IconHome,
  IconInvitations,
  IconGroup,
  IconFlag,
  IconPerson,
  IconDatabase,
} from '@/lib/components/icons'
import Shell from '@/lib/components/shell'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { errorCleared } from '@/store/ui/error'

const LayoutConsole = () => {
  const toast = useToast()
  const error = useAppSelector((state) => state.ui.error.value)
  const navigate = useNavigate()
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
          <Logo />
        </div>
      }
      topBar={<AppBar />}
      items={[
        {
          href: '/console/dashboard',
          icon: <IconHome />,
          primaryText: 'Cloud Console Overview',
          secondaryText: 'Basic information about instance',
        },
        {
          href: '/console/users',
          icon: <IconPerson />,
          primaryText: 'User Management',
          secondaryText: 'Manage users of your cloud instance',
        },
        {
          href: '/console/groups',
          icon: <IconGroup />,
          primaryText: 'Group Management',
          secondaryText: 'Manage groups of your cloud instance',
        },
        {
          href: '/console/workspaces',
          icon: <IconWorkspaces />,
          primaryText: 'Workspace Management',
          secondaryText: 'Manage workspaces of your cloud instance',
        },
        {
          href: '/console/organizations',
          icon: <IconFlag />,
          primaryText: 'Organization Management',
          secondaryText: 'Manage workspaces of your cloud instance',
        },
        {
          href: '/console/invitations',
          icon: <IconInvitations />,
          primaryText: 'Invitation Management',
          secondaryText: 'Manage invitations of your cloud instance',
        },
        {
          href: '/console/database',
          icon: <IconDatabase />,
          primaryText: 'Database Management',
          secondaryText: 'Manage database of your cloud instance',
        },
      ]}
    >
      <Outlet />
    </Shell>
  )
}

export default LayoutConsole
