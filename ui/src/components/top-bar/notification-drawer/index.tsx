import { useEffect, useRef } from 'react'
import {
  Divider,
  Drawer as ChakraDrawer,
  DrawerBody,
  DrawerCloseButton,
  DrawerContent,
  DrawerHeader,
  DrawerOverlay,
  IconButton,
  useDisclosure,
} from '@chakra-ui/react'
import cx from 'classnames'
import NotificationAPI from '@/client/api/notification'
import { swrConfig } from '@/client/options'
import { IconNotifications } from '@/lib/components/icons'
import NotificationBadge from '@/lib/components/notification-badge'
import { useAppDispatch } from '@/store/hook'
import { mutateUpdated } from '@/store/ui/notifications'
import NotificationDrawerItem from './notification-drawer-item'

const TopBarNotificationDrawer = () => {
  const dispatch = useAppDispatch()
  const buttonRef = useRef<HTMLButtonElement>(null)
  const { isOpen, onOpen, onClose } = useDisclosure()
  const { data: notfications, mutate } = NotificationAPI.useGetAll(swrConfig())

  useEffect(() => {
    if (mutate) {
      dispatch(mutateUpdated(mutate))
    }
  }, [mutate])

  return (
    <>
      <NotificationBadge hasBadge={notfications && notfications.length > 0}>
        <IconButton
          ref={buttonRef}
          icon={<IconNotifications />}
          aria-label=""
          onClick={onOpen}
        />
      </NotificationBadge>
      <ChakraDrawer
        isOpen={isOpen}
        placement="right"
        onClose={onClose}
        finalFocusRef={buttonRef}
      >
        <DrawerOverlay />
        <DrawerContent>
          <DrawerCloseButton />
          <DrawerHeader>Notifications</DrawerHeader>
          <DrawerBody>
            {notfications && notfications.length > 0 ? (
              <div className={cx('flex', 'flex-col', 'gap-1.5')}>
                {notfications.map((nofitication, index) => (
                  <div
                    key={index}
                    className={cx('flex', 'flex-col', 'gap-1.5')}
                  >
                    <NotificationDrawerItem notification={nofitication} />
                    {index !== notfications.length - 1 && <Divider />}
                  </div>
                ))}
              </div>
            ) : (
              <span>There are no notifications.</span>
            )}
          </DrawerBody>
        </DrawerContent>
      </ChakraDrawer>
    </>
  )
}

export default TopBarNotificationDrawer
