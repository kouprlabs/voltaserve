// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { useCallback, useEffect, useRef } from 'react'
import {
  Drawer as ChakraDrawer,
  DrawerBody,
  DrawerCloseButton,
  DrawerContent,
  DrawerHeader,
  DrawerOverlay,
  DrawerFooter,
  IconButton,
  useDisclosure,
  Button,
} from '@chakra-ui/react'
import cx from 'classnames'
import UploadList from '@/components/upload/upload-list'
import { IconClearAll, IconUpload } from '@/lib/components/icons'
import NotificationBadge from '@/lib/components/notification-badge'
import { completedUploadsCleared } from '@/store/entities/uploads'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { drawerDidClose } from '@/store/ui/uploads'

const UploadDrawer = () => {
  const dispatch = useAppDispatch()
  const hasPendingUploads = useAppSelector(
    (state) =>
      state.entities.uploads.items.filter((e) => !e.completed).length > 0,
  )
  const isDrawerOpen = useAppSelector((state) => state.ui.uploads.isDrawerOpen)
  const hasCompleted = useAppSelector(
    (state) =>
      state.entities.uploads.items.filter((e) => e.completed).length > 0,
  )
  const { isOpen, onOpen, onClose } = useDisclosure()
  const buttonRef = useRef<HTMLButtonElement>(null)

  useEffect(() => {
    if (isDrawerOpen) {
      onOpen()
    } else {
      onClose()
    }
  }, [isDrawerOpen, onOpen, onClose])

  const handleClearCompleted = useCallback(() => {
    dispatch(completedUploadsCleared())
  }, [dispatch])

  return (
    <>
      <NotificationBadge hasBadge={hasPendingUploads}>
        <IconButton
          ref={buttonRef}
          icon={<IconUpload />}
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
          <DrawerHeader>Uploads</DrawerHeader>
          <DrawerBody className={cx('p-2')}>
            <UploadList />
          </DrawerBody>
          <DrawerFooter>
            {hasCompleted ? (
              <Button
                className={cx('w-full')}
                size="sm"
                leftIcon={<IconClearAll />}
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

export default UploadDrawer
