import { useCallback, useState } from 'react'
import { useParams } from 'react-router-dom'
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
import { useSWRConfig } from 'swr'
import cx from 'classnames'
import SnapshotAPI from '@/client/api/snapshot'
import { useAppSelector } from '@/store/hook'
import { snapshotDeleteModalDidClose } from '@/store/ui/files'

const FileSnapshotDelete = () => {
  const { mutate } = useSWRConfig()
  const dispatch = useDispatch()
  const selection = useAppSelector((state) => state.ui.files.selection)
  const snapshotSelection = useAppSelector(
    (state) => state.ui.files.snapshotSelection,
  )
  const isModalOpen = useAppSelector(
    (state) => state.ui.files.isSnapshotDeleteModalOpen,
  )
  const [isLoading, setIsLoading] = useState(false)

  const handleDelete = useCallback(async () => {
    if (selection.length === 1 && snapshotSelection.length === 1) {
      try {
        setIsLoading(true)
        await SnapshotAPI.delete(selection[0], snapshotSelection[0])
        mutate(`/files/${selection[0]}/snapshots`)
        dispatch(snapshotDeleteModalDidClose())
      } finally {
        setIsLoading(false)
      }
    }
  }, [snapshotSelection, mutate, dispatch])

  return (
    <Modal
      isOpen={isModalOpen}
      onClose={() => dispatch(snapshotDeleteModalDidClose())}
      closeOnOverlayClick={false}
    >
      <ModalOverlay />
      <ModalContent>
        {snapshotSelection.length > 1 ? (
          <ModalHeader>
            Delete {snapshotSelection.length} Snapshot(s)
          </ModalHeader>
        ) : (
          <ModalHeader>Delete Snapshot</ModalHeader>
        )}
        <ModalCloseButton />
        <ModalBody>
          {snapshotSelection.length > 1 ? (
            <span>
              Are you sure you would like to delete ({snapshotSelection.length})
              snapshot(s)?
            </span>
          ) : (
            <span>Are you sure you would like to delete this snapshot?</span>
          )}
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
              onClick={handleDelete}
            >
              Delete
            </Button>
          </div>
        </ModalFooter>
      </ModalContent>
    </Modal>
  )
}

export default FileSnapshotDelete
