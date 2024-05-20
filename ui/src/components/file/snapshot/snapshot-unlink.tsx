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
import {
  snapshotDeleteModalDidClose,
  snapshotSelectionUpdated,
} from '@/store/ui/files'

const FileSnapshotUnlink = () => {
  const dispatch = useDispatch()
  const id = useAppSelector((state) =>
    state.ui.files.snapshotSelection.length > 0
      ? state.ui.files.snapshotSelection[0]
      : undefined,
  )
  const fileId = useAppSelector((state) =>
    state.ui.files.selection.length > 0
      ? state.ui.files.selection[0]
      : undefined,
  )
  const mutate = useAppSelector((state) => state.ui.files.snapshotMutate)
  const isModalOpen = useAppSelector(
    (state) => state.ui.files.isSnapshotDeleteModalOpen,
  )
  const [isLoading, setIsLoading] = useState(false)

  const handleUnlink = useCallback(async () => {
    async function unlink(id: string, fileId: string) {
      setIsLoading(true)
      try {
        await SnapshotAPI.unlink(id, { fileId })
        await mutate?.()
        dispatch(snapshotSelectionUpdated([]))
        dispatch(snapshotDeleteModalDidClose())
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
      onClose={() => dispatch(snapshotDeleteModalDidClose())}
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
              onClick={() => dispatch(snapshotDeleteModalDidClose())}
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
