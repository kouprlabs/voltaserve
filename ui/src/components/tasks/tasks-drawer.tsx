import { useEffect, useRef } from 'react'
import { useNavigate } from 'react-router-dom'
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
import TaskAPI, { SortOrder } from '@/client/api/task'
import { swrConfig } from '@/client/options'
import { taskPaginationStorage } from '@/infra/pagination'
import { IconStacks } from '@/lib/components/icons'
import NotificationBadge from '@/lib/components/notification-badge'
import PagePagination from '@/lib/components/page-pagination'
import usePagePagination from '@/lib/hooks/page-pagination'
import { useAppSelector } from '@/store/hook'
import TasksDrawerItem from './tasks-item'
import TasksList from './tasks-list'

const TasksDrawer = () => {
  const navigate = useNavigate()
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
        onClose={onClose}
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
