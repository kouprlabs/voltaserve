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
import cx from 'classnames'
import TaskAPI from '@/client/api/task'
import { swrConfig } from '@/client/options'
import { IconStacks } from '@/lib/components/icons'
import NotificationBadge from '@/lib/components/notification-badge'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { drawerDidClose, mutateCountUpdated } from '@/store/ui/tasks'
import TasksList from './task-list'

const TaskDrawer = () => {
  const dispatch = useAppDispatch()
  const buttonRef = useRef<HTMLButtonElement>(null)
  const { isOpen, onOpen, onClose } = useDisclosure()
  const isDrawerOpen = useAppSelector((state) => state.ui.tasks.isDrawerOpen)
  const { data: count, mutate: mutateCount } = TaskAPI.useGetCount(swrConfig())

  useEffect(() => {
    if (isDrawerOpen) {
      onOpen()
    } else {
      onClose()
    }
  }, [isDrawerOpen, onOpen, onClose])

  useEffect(() => {
    if (mutateCount) {
      dispatch(mutateCountUpdated(mutateCount))
    }
  }, [mutateCount, dispatch])

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
        size="sm"
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
          <DrawerBody className={cx('p-2')}>
            <TasksList />
          </DrawerBody>
        </DrawerContent>
      </ChakraDrawer>
    </>
  )
}

export default TaskDrawer
