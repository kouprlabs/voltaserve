import { useEffect } from 'react'
import { useLocation } from 'react-router-dom'
import cx from 'classnames'
import AccountMenu from '@/components/account/menu'
import TaskDrawer from '@/components/tasks/tasks-drawer'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { activeNavChanged, NavType } from '@/store/ui/nav'
import UploadsDrawer from '../uploads/uploads-drawer'
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
        {activeNav === NavType.Workspaces && <CreateWorkspaceButton />}
        {activeNav === NavType.Groups && <CreateGroupButton />}
        {activeNav === NavType.Organizations && <CreateOrganizationButton />}
        <UploadsDrawer />
        <TaskDrawer />
        <AccountMenu />
      </div>
    </div>
  )
}

export default AppBar