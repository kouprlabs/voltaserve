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
import { useLocation } from 'react-router-dom'
import cx from 'classnames'
import AccountMenu from '@/components/account/menu'
import AdminButton from '@/components/admin/admin-button'
import TaskDrawer from '@/components/task/task-drawer'
import { getAdminStatus } from '@/infra/token'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { activeNavChanged, NavType } from '@/store/ui/nav'
import UploadDrawer from '../upload/upload-drawer'
import {
  CreateGroupButton,
  CreateOrganizationButton,
  CreateWorkspaceButton,
} from './app-bar-buttons'
import AppBarSearch from './app-bar-search'

const AppBar = () => {
  const dispatch = useAppDispatch()
  const location = useLocation()
  const activeNav = useAppSelector((state) => state.ui.nav.active)

  useEffect(() => {
    if (location.pathname.startsWith('/account')) {
      dispatch(activeNavChanged(NavType.Account))
    }
    if (location.pathname.startsWith('/organization')) {
      dispatch(activeNavChanged(NavType.Organizations))
    }
    if (location.pathname.startsWith('/group')) {
      dispatch(activeNavChanged(NavType.Groups))
    }
    if (location.pathname.startsWith('/workspace')) {
      dispatch(activeNavChanged(NavType.Workspaces))
    }
    if (location.pathname.startsWith('/admin')) {
      dispatch(activeNavChanged(NavType.Admin))
    }
  }, [location, dispatch])

  return (
    <div
      className={cx(
        'flex',
        'flex-row',
        'items-center',
        'gap-2',
        'shrink-0',
        'p-1.5',
        'w-full',
      )}
    >
      <div className={cx('grow')}>
        <AppBarSearch />
      </div>
      <div className={cx('flex', 'flex-row', 'items-center', 'gap-1.5')}>
        {activeNav === NavType.Workspaces ? <CreateWorkspaceButton /> : null}
        {activeNav === NavType.Groups ? <CreateGroupButton /> : null}
        {activeNav === NavType.Organizations ? (
          <CreateOrganizationButton />
        ) : null}
        {getAdminStatus() ? <AdminButton /> : null}
        <UploadDrawer />
        <TaskDrawer />
        <AccountMenu />
      </div>
    </div>
  )
}

export default AppBar
