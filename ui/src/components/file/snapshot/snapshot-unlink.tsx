import { useCallback, useState } from 'react'
import { useDispatch } from 'react-redux'
import {
  Button,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
} from '@chakra-ui/react'
import cx from 'classnames'
import SnapshotAPI from '@/client/api/snapshot'
import { useAppSelector } from '@/store/hook'
import { deleteModalDidClose, selectionUpdated } from '@/store/ui/snapshots'

const FileSnapshotUnlink = () => {
  const dispatch = useDispatch()
  const id = useAppSelector((state) =>
    state.ui.snapshots.selection.length > 0
      ? state.ui.snapshots.selection[0]
      : undefined,
  )
  const fileId = useAppSelector((state) =>
    state.ui.files.selection.length > 0
      ? state.ui.files.selection[0]
      : undefined,
  )
  const mutate = useAppSelector((state) => state.ui.snapshots.snapshotMutate)
  const isModalOpen = useAppSelector(
    (state) => state.ui.snapshots.isDeleteModalOpen,
  )
  const [isLoading, setIsLoading] = useState(false)

  const handleUnlink = useCallback(async () => {
    async function unlink(id: string, fileId: string) {
      setIsLoading(true)
      try {
        await SnapshotAPI.unlink(id, { fileId })
        await mutate?.()
        dispatch(selectionUpdated([]))
        dispatch(deleteModalDidClose())
      } catch (error) {
        setIsLoading(false)
      } finally {
        setIsLoading(false)
      }
    }
    if (id && fileId) {
      unlink(id, fileId)
    }
  }, [id, fileId, dispatch, mutate])

  return (
    <Modal
      isOpen={isModalOpen}
      onClose={() => dispatch(deleteModalDidClose())}
      closeOnOverlayClick={false}
    >
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>Unlink Snapshot</ModalHeader>
        <ModalCloseButton />
        <ModalBody>
          <span>Are you sure you would like to unkink this snapshot?</span>
        </ModalBody>
        <ModalFooter>
          <div className={cx('flex', 'flex-row', 'items-center', 'gap-1')}>
            <Button
              type="button"
              variant="outline"
              colorScheme="blue"
              disabled={isLoading}
              onClick={() => dispatch(deleteModalDidClose())}
            >
              Cancel
            </Button>
            <Button
              type="submit"
              variant="solid"
              colorScheme="red"
              isLoading={isLoading}
              onClick={handleUnlink}
            >
              Unlink
            </Button>
          </div>
        </ModalFooter>
      </ModalContent>
    </Modal>
  )
}

export default FileSnapshotUnlink
