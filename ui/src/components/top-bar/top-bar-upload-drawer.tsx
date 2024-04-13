import { useCallback, useEffect, useRef } from 'react'
import {
  Circle,
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
import UploadList from '@/components/file/upload/upload-list'
import { IconDeleteListItem, IconUpload } from '@/lib'
import { completedUploadsCleared } from '@/store/entities/uploads'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { uploadsDrawerClosed } from '@/store/ui/uploads-drawer'

const TopBarUploadDrawer = () => {
  const dispatch = useAppDispatch()
  const hasPendingUploads = useAppSelector(
    (state) =>
      state.entities.uploads.items.filter((e) => !e.completed).length > 0,
  )
  const openDrawer = useAppSelector((state) => state.ui.uploadsDrawer.open)
  const hasCompleted = useAppSelector(
    (state) =>
      state.entities.uploads.items.filter((e) => e.completed).length > 0,
  )
  const { isOpen, onOpen, onClose } = useDisclosure()
  const buttonRef = useRef<HTMLButtonElement>(null)

  useEffect(() => {
    if (openDrawer) {
      onOpen()
    } else {
      onClose()
    }
  }, [openDrawer, onOpen, onClose])

  const handleClearCompleted = useCallback(() => {
    dispatch(completedUploadsCleared())
  }, [dispatch])

  return (
    <>
      <div className={cx('flex', 'items-center', 'justify-center', 'relative')}>
        <IconButton
          ref={buttonRef}
          icon={<IconUpload size="14px" />}
          aria-label=""
          onClick={onOpen}
        />
        {hasPendingUploads && (
          <Circle size="15px" bg="red" position="absolute" top={0} right={0} />
        )}
      </div>
      <ChakraDrawer
        isOpen={isOpen}
        placement="right"
        onClose={() => {
          onClose()
          dispatch(uploadsDrawerClosed())
        }}
        finalFocusRef={buttonRef}
      >
        <DrawerOverlay />
        <DrawerContent>
          <DrawerCloseButton />
          <DrawerHeader>Uploads</DrawerHeader>
          <DrawerBody>
            <UploadList />
          </DrawerBody>
          <DrawerFooter>
            {hasCompleted && (
              <Button
                width="100%"
                leftIcon={<IconDeleteListItem size="22px" />}
                onClick={handleClearCompleted}
              >
                Clear Completed Items
              </Button>
            )}
          </DrawerFooter>
        </DrawerContent>
      </ChakraDrawer>
    </>
  )
}

export default TopBarUploadDrawer
