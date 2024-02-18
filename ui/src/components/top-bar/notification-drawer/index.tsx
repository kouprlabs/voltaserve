import { useRef } from 'react'
import {
  Divider,
  Drawer as ChakraDrawer,
  DrawerBody,
  DrawerCloseButton,
  DrawerContent,
  DrawerHeader,
  DrawerOverlay,
  Text,
  IconButton,
  useDisclosure,
  Circle,
} from '@chakra-ui/react'
import { IconNotification } from '@koupr/ui'
import classNames from 'classnames'
import NotificationAPI from '@/client/api/notification'
import { swrConfig } from '@/client/options'
import NotificationItem from './notification-item'

const NotificationDrawer = () => {
  const buttonRef = useRef<HTMLButtonElement>(null)
  const { isOpen, onOpen, onClose } = useDisclosure()
  const { data: notfications } = NotificationAPI.useGetAll(swrConfig())

  return (
    <>
      <div
        className={classNames(
          'flex',
          'items-center',
          'justify-center',
          'relative',
        )}
      >
        <IconButton
          ref={buttonRef}
          icon={<IconNotification />}
          aria-label=""
          onClick={onOpen}
        />
        {notfications && notfications.length > 0 && (
          <Circle size="15px" bg="red" position="absolute" top={0} right={0} />
        )}
      </div>
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
              <div className={classNames('flex', 'flex-col', 'gap-1.5')}>
                {notfications.map((n, index) => (
                  <div
                    key={index}
                    className={classNames('flex', 'flex-col', 'gap-1.5')}
                  >
                    <NotificationItem notification={n} />
                    {index !== notfications.length - 1 && <Divider />}
                  </div>
                ))}
              </div>
            ) : (
              <Text>There are no notifications.</Text>
            )}
          </DrawerBody>
        </DrawerContent>
      </ChakraDrawer>
    </>
  )
}

export default NotificationDrawer
