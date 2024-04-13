import { useEffect, useMemo, useState } from 'react'
import { Outlet, useLocation, useNavigate } from 'react-router-dom'
import {
  Avatar,
  Button,
  Heading,
  IconButton,
  Tab,
  TabList,
  Tabs,
  Tag,
} from '@chakra-ui/react'
import cx from 'classnames'
import NotificationAPI from '@/client/api/notification'
import UserAPI from '@/client/idp/user'
import { swrConfig } from '@/client/options'
import AccountEditPicture from '@/components/account/edit-picture'
import { IconEdit } from '@/lib'

const AccountLayout = () => {
  const location = useLocation()
  const navigate = useNavigate()
  const [isImageModalOpen, setIsImageModalOpen] = useState(false)
  const { data: user } = UserAPI.useGet(swrConfig())
  const { data: notfications } = NotificationAPI.useGetAll(swrConfig())
  const invitationCount = useMemo(
    () => notfications?.filter((e) => e.type === 'new_invitation').length,
    [notfications],
  )
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
              width="165px"
              height="165px"
              size="2xl"
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
            width="100%"
            type="submit"
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
                  <Tag className={cx('rounded-full')}>{invitationCount}</Tag>
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
