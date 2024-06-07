import { useEffect, useState } from 'react'
import { Outlet, useLocation, useNavigate } from 'react-router-dom'
import {
  Avatar,
  Button,
  Heading,
  IconButton,
  Tab,
  TabList,
  Tabs,
} from '@chakra-ui/react'
import cx from 'classnames'
import InvitationAPI from '@/client/api/invitation'
import UserAPI from '@/client/idp/user'
import { swrConfig } from '@/client/options'
import AccountEditPicture from '@/components/account/edit-picture'
import { IconEdit } from '@/lib/components/icons'
import NumberTag from '@/lib/components/number-tag'
import { useAppDispatch } from '@/store/hook'
import { mutateUpdated } from '@/store/ui/account'

const AccountLayout = () => {
  const location = useLocation()
  const navigate = useNavigate()
  const dispatch = useAppDispatch()
  const [isImageModalOpen, setIsImageModalOpen] = useState(false)
  const { data: user, mutate } = UserAPI.useGet(swrConfig())
  const { data: invitationCount } =
    InvitationAPI.useGetIncomingCount(swrConfig())
  const [tabIndex, setTabIndex] = useState(0)

  useEffect(() => {
    const segments = location.pathname.split('/')
    const segment = segments[segments.length - 1]
    if (segment === 'settings') {
      setTabIndex(0)
    } else if (segment === 'invitation') {
      setTabIndex(1)
    }
  }, [location])

  useEffect(() => {
    if (mutate) {
      dispatch(mutateUpdated(mutate))
    }
  }, [mutate, dispatch])

  if (!user) {
    return null
  }

  return (
    <div className={cx('flex', 'flex-row', 'gap-2.5')}>
      <div
        className={cx('flex', 'flex-col', 'gap-2', 'items-center', 'w-[250px]')}
      >
        <div className={cx('flex', 'flex-col', 'gap-2', 'items-center')}>
          <div className={cx('relative', 'shrink-0')}>
            <Avatar
              name={user.fullName}
              src={user.picture}
              size="2xl"
              className={cx('w-[165px]', 'h-[165px]')}
            />
            <IconButton
              icon={<IconEdit />}
              variant="solid-gray"
              right="5px"
              bottom="10px"
              position="absolute"
              zIndex={1000}
              aria-label=""
              onClick={() => setIsImageModalOpen(true)}
            />
          </div>
          <Heading className={cx('text-center', 'text-heading')}>
            {user.fullName}
          </Heading>
        </div>
        <div className={cx('w-full', 'gap-1')}>
          <Button
            variant="outline"
            colorScheme="red"
            type="submit"
            className={cx('w-full')}
            onClick={() => navigate('/sign-out')}
          >
            Sign Out
          </Button>
        </div>
      </div>
      <div className={cx('w-full', 'pb-1.5')}>
        <Tabs
          variant="solid-rounded"
          colorScheme="gray"
          index={tabIndex}
          className={cx('pb-2.5')}
        >
          <TabList>
            <Tab onClick={() => navigate('/account/settings')}>Settings</Tab>
            <Tab onClick={() => navigate('/account/invitation')}>
              <div
                className={cx('flex', 'flex-row', 'items-center', 'gap-0.5')}
              >
                <span>Invitations</span>
                {invitationCount && invitationCount > 0 ? (
                  <NumberTag isActive={tabIndex === 1}>
                    {invitationCount}
                  </NumberTag>
                ) : null}
              </div>
            </Tab>
          </TabList>
        </Tabs>
        <Outlet />
      </div>
      <AccountEditPicture
        open={isImageModalOpen}
        user={user}
        onClose={() => setIsImageModalOpen(false)}
      />
    </div>
  )
}

export default AccountLayout
