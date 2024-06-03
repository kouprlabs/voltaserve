import { useRef } from 'react'
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
import TaskDrawerItem from './task-drawer-item'

const TaskDrawer = () => {
  const navigate = useNavigate()
  const buttonRef = useRef<HTMLButtonElement>(null)
  const { isOpen, onOpen, onClose } = useDisclosure()
  const { page, size, steps, setPage, setSize } = usePagePagination({
    navigate,
    location,
    storage: taskPaginationStorage(),
  })
  const { data: list } = TaskAPI.useList(
    { page, size, sortOrder: SortOrder.Desc },
    swrConfig(),
  )

  return (
    <>
      <NotificationBadge hasBadge={list && list.totalElements > 0}>
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
            {list && list.data.length > 0 ? (
              <div className={cx('flex', 'flex-col', 'gap-1.5')}>
                {list.data.map((task, index) => (
                  <div
                    key={index}
                    className={cx('flex', 'flex-col', 'gap-1.5')}
                  >
                    <TaskDrawerItem task={task} />
                  </div>
                ))}
              </div>
            ) : (
              <span>There are no tasks.</span>
            )}
            {list ? (
              <PagePagination
                style={{ alignSelf: 'end' }}
                totalElements={list.totalElements}
                totalPages={list.totalPages}
                page={page}
                size={size}
                steps={steps}
                setPage={setPage}
                setSize={setSize}
              />
            ) : null}
          </DrawerBody>
        </DrawerContent>
      </ChakraDrawer>
    </>
  )
}

export default TaskDrawer
