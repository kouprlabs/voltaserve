import { ReactElement, useEffect } from 'react'
import { Stack, useToast } from '@chakra-ui/react'
import { Flex } from '@chakra-ui/react'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { errorCleared } from '@/store/ui/error'
import Drawer from '@/components/common/drawer'
import DrawerItem from '@/components/common/drawer-item'
import {
  IconGroup,
  IconOrganization,
  IconWorkspace,
} from '@/components/common/icon'
import TopBar from '@/components/top-bar'
import variables from '@/theme/variables'

type DrawerLayoutProps = {
  children?: ReactElement
}

const DrawerLayout = ({ children }: DrawerLayoutProps) => {
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
    <Stack direction="row" spacing={0} h="100%">
      <Drawer localStorageNamespace="main">
        <DrawerItem
          href="/workspace"
          icon={<IconWorkspace fontSize="18px" />}
          primaryText="Workspaces"
          secondaryText="Isolated containers for files and folders."
        />
        <DrawerItem
          href="/group"
          icon={<IconGroup fontSize="16px" />}
          primaryText="Groups"
          secondaryText="Allows assigning permissions to a group of users."
        />
        <DrawerItem
          href="/organization"
          icon={<IconOrganization fontSize="18px" />}
          primaryText="Organizations"
          secondaryText="Umbrellas for workspaces and users."
        />
      </Drawer>
      <Flex direction="column" alignItems="center" h="100%" w="100%">
        <TopBar />
        <Flex
          direction="column"
          width={{ base: 'full', '2xl': '1250px' }}
          px={variables.spacing2Xl}
          pt={variables.spacing2Xl}
          overflowY="auto"
          overflowX="hidden"
          flexGrow={1}
        >
          {children}
        </Flex>
      </Flex>
    </Stack>
  )
}

export default DrawerLayout
