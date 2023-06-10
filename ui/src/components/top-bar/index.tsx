import { useEffect } from 'react'
import { useLocation } from 'react-router-dom'
import { Box, HStack } from '@chakra-ui/react'
import { variables } from '@koupr/ui'
import AccountMenu from '@/components/common/account-menu'
import NotificationDrawer from '@/components/common/notification-drawer'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { activeNavChanged, NavType } from '@/store/ui/nav'
import {
  CreateGroupButton,
  CreateOrganizationButton,
  CreateWorkspaceButton,
} from './buttons'
import Search from './search'
import UploadDrawer from './upload-drawer'

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
    <HStack
      padding={`0 30px`}
      w="100%"
      h="80px"
      align="center"
      spacing={variables.spacingMd}
      flexShrink={0}
    >
      <Box flexGrow={1}>
        <Search />
      </Box>
      <HStack spacing={variables.spacing}>
        {activeNav === NavType.Workspaces && <CreateWorkspaceButton />}
        {activeNav === NavType.Groups && <CreateGroupButton />}
        {activeNav === NavType.Organizations && <CreateOrganizationButton />}
        <UploadDrawer />
        <NotificationDrawer />
        <AccountMenu />
      </HStack>
    </HStack>
  )
}

export default TopBar
