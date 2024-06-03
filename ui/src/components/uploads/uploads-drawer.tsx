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
import UploadsList from '@/components/uploads/uploads-list'
import { IconClearAll, IconUpload } from '@/lib/components/icons'
import { completedUploadsCleared } from '@/store/entities/uploads'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { drawerDidClose } from '@/store/ui/uploads'

const UploadsDrawer = () => {
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
      <div className={cx('flex', 'items-center', 'justify-center', 'relative')}>
        <IconButton
          ref={buttonRef}
          icon={<IconUpload />}
          aria-label=""
          onClick={onOpen}
        />
        {hasPendingUploads ? (
          <Circle size="15px" bg="red" position="absolute" top={0} right={0} />
        ) : null}
      </div>
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
          <DrawerHeader>Uploads</DrawerHeader>
          <DrawerBody>
            <UploadsList />
          </DrawerBody>
          <DrawerFooter>
            {hasCompleted ? (
              <Button
                className={cx('w-full')}
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

export default UploadsDrawer
