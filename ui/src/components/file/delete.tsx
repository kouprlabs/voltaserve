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
import { useSWRConfig } from 'swr'
import FileAPI, { List } from '@/client/api/file'
import useFileListSearchParams from '@/hooks/use-file-list-params'
import { useAppSelector } from '@/store/hook'
import { deleteModalDidClose, selectionUpdated } from '@/store/ui/files'

const Delete = () => {
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

export default Delete
