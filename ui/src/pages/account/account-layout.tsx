import { useEffect, useMemo, useState } from 'react'
import { Outlet, useLocation, useNavigate } from 'react-router-dom'
import {
  Avatar,
  SegmentedControl,
  Badge,
  Text,
  Button,
  Heading,
  IconButton,
} from '@radix-ui/themes'
import cx from 'classnames'
import initials from 'initials'
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
  const [segment, setSegment] = useState<string>()

  useEffect(() => {
    const segments = location.pathname.split('/')
    setSegment(segments[segments.length - 1])
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
              fallback={initials(user.fullName)}
              src={user.picture}
              size="8"
              radius="full"
              className={cx('w-[165px]', 'h-[165px]')}
            />
            <IconButton
              className={cx(
                'absolute',
                'right-[5px]',
                'bottom-[10px]',
                'z-[1000]',
              )}
              radius="full"
              aria-label=""
              onClick={() => setIsImageModalOpen(true)}
            >
              <IconEdit />
            </IconButton>
          </div>
          <Heading size="5">{user.fullName}</Heading>
        </div>
        <div className={cx('w-full', 'gap-1')}>
          <Button
            variant="outline"
            color="crimson"
            radius="full"
            type="submit"
            className={cx('w-full')}
            onClick={() => navigate('/sign-out')}
          >
            Sign Out
          </Button>
        </div>
      </div>
      <div className={cx('flex', 'flex-col', 'gap-2.5', 'w-full')}>
        <SegmentedControl.Root
          radius="full"
          defaultValue={segment}
          className={cx('self-start')}
        >
          <SegmentedControl.Item
            value="settings"
            onClick={() => navigate('/account/settings')}
          >
            Settings
          </SegmentedControl.Item>
          <SegmentedControl.Item
            value="invitation"
            onClick={() => navigate('/account/invitation')}
          >
            <div className={cx('flex', 'items-center', 'gap-0.5')}>
              <Text>Invitations</Text>
              {invitationCount && invitationCount > 0 ? (
                <Badge radius="full">{invitationCount}</Badge>
              ) : null}
            </div>
          </SegmentedControl.Item>
        </SegmentedControl.Root>
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
