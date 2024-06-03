import { useEffect, useRef } from 'react'
import {
  Drawer as ChakraDrawer,
  DrawerBody,
  DrawerCloseButton,
  DrawerContent,
  DrawerHeader,
  DrawerOverlay,
  IconButton,
  useDisclosure,
} from '@chakra-ui/react'
import TaskAPI from '@/client/api/task'
import { swrConfig } from '@/client/options'
import { IconStacks } from '@/lib/components/icons'
import NotificationBadge from '@/lib/components/notification-badge'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { drawerDidClose } from '@/store/ui/tasks'
import TasksList from './tasks-list'

const TasksDrawer = () => {
  const dispatch = useAppDispatch()
  const buttonRef = useRef<HTMLButtonElement>(null)
  const { isOpen, onOpen, onClose } = useDisclosure()
  const isDrawerOpen = useAppSelector((state) => state.ui.tasks.isDrawerOpen)
  const { data: count } = TaskAPI.useGetCount(swrConfig())

  useEffect(() => {
    if (isDrawerOpen) {
      onOpen()
    } else {
      onClose()
    }
  }, [isDrawerOpen, onOpen, onClose])

  return (
    <>
      <NotificationBadge hasBadge={count !== undefined && count > 0}>
        <IconButton
          ref={buttonRef}
          icon={<IconStacks />}
          aria-label=""
          onClick={onOpen}
        />
      </NotificationBadge>
      <ChakraDrawer
        isOpen={isOpen}
        placement="right"
        onClose={() => {
          onClose()
          dispatch(drawerDidClose())
        }}
        finalFocusRef={buttonRef}
      >
        <DrawerOverlay />
        <DrawerContent>
          <DrawerCloseButton />
          <DrawerHeader>Tasks</DrawerHeader>
          <DrawerBody>
            <TasksList />
          </DrawerBody>
        </DrawerContent>
      </ChakraDrawer>
    </>
  )
}

export default TasksDrawer
