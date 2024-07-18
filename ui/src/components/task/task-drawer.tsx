// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { useCallback, useEffect, useRef, useState } from 'react'
import {
  Button,
  Drawer as ChakraDrawer,
  DrawerBody,
  DrawerCloseButton,
  DrawerContent,
  DrawerFooter,
  DrawerHeader,
  DrawerOverlay,
  IconButton,
  useDisclosure,
} from '@chakra-ui/react'
import cx from 'classnames'
import TaskAPI from '@/client/api/task'
import { swrConfig } from '@/client/options'
import { IconClearAll, IconStacks } from '@/lib/components/icons'
import NotificationBadge from '@/lib/components/notification-badge'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { drawerDidClose, mutateCountUpdated } from '@/store/ui/tasks'
import TasksList from './task-list'

const TaskDrawer = () => {
  const dispatch = useAppDispatch()
  const buttonRef = useRef<HTMLButtonElement>(null)
  const { isOpen, onOpen, onClose } = useDisclosure()
  const [isDismissing, setIsDismissing] = useState(false)
  const isDrawerOpen = useAppSelector((state) => state.ui.tasks.isDrawerOpen)
  const mutateList = useAppSelector((state) => state.ui.tasks.mutateList)
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

  const handleClearCompleted = useCallback(async () => {
    try {
      setIsDismissing(true)
      await TaskAPI.dismissAll()
      mutateList?.(await TaskAPI.list())
    } finally {
      setIsDismissing(false)
    }
  }, [dispatch, mutateList])

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
          <DrawerFooter>
            {count && count > 0 ? (
              <Button
                className={cx('w-full')}
                size="sm"
                leftIcon={<IconClearAll />}
                isLoading={isDismissing}
                onClick={handleClearCompleted}
              >
                Clear Completed Items
              </Button>
            ) : null}
          </DrawerFooter>
        </DrawerContent>
      </ChakraDrawer>
    </>
  )
}

export default TaskDrawer
