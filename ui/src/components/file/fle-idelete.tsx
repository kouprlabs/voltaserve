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
import FileAPI, { List } from '@/client/api/file'
import useFileListSearchParams from '@/hooks/use-file-list-params'
import { useAppSelector } from '@/store/hook'
import { deleteModalDidClose, selectionUpdated } from '@/store/ui/files'

const FileDelete = () => {
  const { mutate } = useSWRConfig()
  const { fileId } = useParams()
  const dispatch = useDispatch()
  const selection = useAppSelector((state) => state.ui.files.selection)
  const isModalOpen = useAppSelector(
    (state) => state.ui.files.isDeleteModalOpen,
  )
  const [loading, setLoading] = useState(false)
  const fileListSearchParams = useFileListSearchParams()

  const handleDelete = useCallback(async () => {
    try {
      setLoading(true)
      await FileAPI.batchDelete({ ids: selection })
      await mutate<List>(`/files/${fileId}/list?${fileListSearchParams}`)
      dispatch(selectionUpdated([]))
      dispatch(deleteModalDidClose())
    } finally {
      setLoading(false)
    }
  }, [selection, fileId, fileListSearchParams, mutate, dispatch])

  return (
    <Modal
      isOpen={isModalOpen}
      onClose={() => dispatch(deleteModalDidClose())}
      closeOnOverlayClick={false}
    >
      <ModalOverlay />
      <ModalContent>
        {selection.length > 1 ? (
          <ModalHeader>Delete {selection.length} Item(s)</ModalHeader>
        ) : (
          <ModalHeader>Delete Item</ModalHeader>
        )}
        <ModalCloseButton />
        <ModalBody>
          {selection.length > 1 ? (
            <span>
              Are you sure you would like to delete ({selection.length})
              item(s)?
            </span>
          ) : (
            <span>Are you sure you would like to delete this item?</span>
          )}
        </ModalBody>
        <ModalFooter>
          <div className={cx('flex', 'flex-row', 'items-center', 'gap-1')}>
            <Button
              type="button"
              variant="outline"
              colorScheme="blue"
              disabled={loading}
              onClick={() => dispatch(deleteModalDidClose())}
            >
              Cancel
            </Button>
            <Button
              type="submit"
              variant="solid"
              colorScheme="red"
              isLoading={loading}
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

export default FileDelete
