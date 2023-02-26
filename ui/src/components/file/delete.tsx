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
  Text,
} from '@chakra-ui/react'
import { variables } from '@koupr/ui'
import FileAPI from '@/api/file'
import { filesRemoved } from '@/store/entities/files'
import { useAppSelector } from '@/store/hook'
import { deleteModalDidClose, selectionUpdated } from '@/store/ui/files'

const FileDelete = () => {
  const params = useParams()
  const dispatch = useDispatch()
  const selection = useAppSelector((state) => state.ui.files.selection)
  const isModalOpen = useAppSelector(
    (state) => state.ui.files.isDeleteModalOpen
  )
  const [loading, setLoading] = useState(false)

  const handleDelete = useCallback(async () => {
    try {
      setLoading(true)
      await FileAPI.batchDelete({ ids: selection })
      dispatch(
        filesRemoved({
          id: params.fileId as string,
          files: selection,
        })
      )
      dispatch(selectionUpdated([]))
      dispatch(deleteModalDidClose())
    } finally {
      setLoading(false)
    }
  }, [selection, params, dispatch])

  return (
    <Modal isOpen={isModalOpen} onClose={() => dispatch(deleteModalDidClose())}>
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>Delete File(s)</ModalHeader>
        <ModalCloseButton />
        <ModalBody>
          <Text>
            Are you sure you would like to delete ({selection.length}) item(s)?
          </Text>
        </ModalBody>
        <ModalFooter>
          <Button
            type="button"
            variant="outline"
            colorScheme="blue"
            mr={variables.spacingSm}
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
        </ModalFooter>
      </ModalContent>
    </Modal>
  )
}

export default FileDelete
