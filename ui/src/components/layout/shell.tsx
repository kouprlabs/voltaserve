import { useEffect } from 'react'
import { Outlet } from 'react-router-dom'
import { useToast } from '@chakra-ui/react'
import { IconGroup, IconOrganization, IconWorkspace, Shell } from '@koupr/ui'
import Logo from '@/components/common/logo'
import TopBar from '@/components/top-bar'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { errorCleared } from '@/store/ui/error'

const ShellLayout = () => {
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
      localStorage={{ prefix: 'voltaserve', namespace: 'main' }}
      logo={<Logo className="w-4" />}
      topBar={<TopBar />}
      items={[
        {
          href: '/workspace',
          icon: <IconWorkspace fontSize="18px" />,
          primaryText: 'Workspaces',
          secondaryText: 'Isolated containers for files and folders.',
        },
        {
          href: '/group',
          icon: <IconGroup fontSize="16px" />,
          primaryText: 'Groups',
          secondaryText: 'Allows assigning permissions to a group of users.',
        },
        {
          href: '/organization',
          icon: <IconOrganization fontSize="18px" />,
          primaryText: 'Organizations',
          secondaryText: 'Umbrellas for workspaces and users.',
        },
      ]}
    >
      <Outlet />
    </Shell>
  )
}

export default ShellLayout
