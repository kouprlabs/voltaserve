import { useEffect } from 'react'
import { Outlet } from 'react-router-dom'
import { useToast } from '@chakra-ui/react'
import { cx } from '@emotion/css'
import Logo from '@/components/common/logo'
import TopBar from '@/components/top-bar'
import { IconGroup, IconFlag, IconWorkspaces } from '@/lib/components/icons'
import Shell from '@/lib/components/shell'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { errorCleared } from '@/store/ui/error'

const LayoutShell = () => {
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
          <Logo />
        </div>
      }
      topBar={<TopBar />}
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
    >
      <Outlet />
    </Shell>
  )
}

export default LayoutShell
