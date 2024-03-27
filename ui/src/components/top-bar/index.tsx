import { useEffect } from 'react'
import { useLocation } from 'react-router-dom'
import cx from 'classnames'
import TopBarAccountMenu from '@/components/top-bar/account-menu'
import TopBarNotificationDrawer from '@/components/top-bar/notification-drawer'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { activeNavChanged, NavType } from '@/store/ui/nav'
import {
  CreateGroupButton,
  CreateOrganizationButton,
  CreateWorkspaceButton,
} from './top-bar-buttons'
import TopBarSearch from './top-bar-search'
import TopBarUploadDrawer from './top-bar-upload-drawer'

const TopBar = () => {
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
        'py-0',
        'px-3',
        'w-full',
        'h-[80px]',
      )}
    >
      <div className={cx('grow')}>
        <TopBarSearch />
      </div>
      <div className={cx('flex', 'flex-row', 'items-center', 'gap-1.5')}>
        {activeNav === NavType.Workspaces && <CreateWorkspaceButton />}
        {activeNav === NavType.Groups && <CreateGroupButton />}
        {activeNav === NavType.Organizations && <CreateOrganizationButton />}
        <TopBarUploadDrawer />
        <TopBarNotificationDrawer />
        <TopBarAccountMenu />
      </div>
    </div>
  )
}

export default TopBar
