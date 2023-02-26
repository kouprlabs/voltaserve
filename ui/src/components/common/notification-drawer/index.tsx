import { useRef } from 'react'
import {
  Box,
  Center,
  Divider,
  Drawer as ChakraDrawer,
  DrawerBody,
  DrawerCloseButton,
  DrawerContent,
  DrawerHeader,
  DrawerOverlay,
  Stack,
  Text,
  IconButton,
  useDisclosure,
  Circle,
} from '@chakra-ui/react'
import { variables } from '@koupr/ui'
import { IconNotification } from '@koupr/ui'
import NotificationAPI from '@/api/notification'
import { swrConfig } from '@/api/options'
import NotificationItem from './notification-item'

const NotificationDrawer = () => {
  const buttonRef = useRef<HTMLButtonElement>(null)
  const { isOpen, onOpen, onClose } = useDisclosure()
  const { data: notfications } = NotificationAPI.useGetAll(swrConfig())

  return (
    <>
      <Box>
        <Center position="relative">
          <IconButton
            ref={buttonRef}
            icon={<IconNotification />}
            aria-label=""
            onClick={onOpen}
          />
          {notfications && notfications.length > 0 && (
            <Circle
              size="15px"
              bg="red"
              position="absolute"
              top={0}
              right={0}
            />
          )}
        </Center>
      </Box>
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
              <Stack spacing={variables.spacing}>
                {notfications.map((n, i) => (
                  <Stack key={i} spacing={variables.spacing}>
                    <NotificationItem notification={n} />
                    {i !== notfications.length - 1 && <Divider />}
                  </Stack>
                ))}
              </Stack>
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
