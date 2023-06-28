import { useCallback, useState } from 'react'
import { useParams } from 'react-router-dom'
import {
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalCloseButton,
  ModalBody,
  ModalFooter,
  Button,
} from '@chakra-ui/react'
import { variables } from '@koupr/ui'
import FileAPI from '@/client/api/file'
import { filesRemoved, filesUpdated } from '@/store/entities/files'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { moveModalDidClose, selectionUpdated } from '@/store/ui/files'
import Browse from './browse'

const Move = () => {
  const params = useParams()
  const fileId = params.fileId as string
  const dispatch = useAppDispatch()
  const selection = useAppSelector((state) => state.ui.files.selection)
  const isModalOpen = useAppSelector((state) => state.ui.files.isMoveModalOpen)
  const sortBy = useAppSelector((state) => state.ui.files.sortBy)
  const sortOrder = useAppSelector((state) => state.ui.files.sortOrder)
  const [loading, setLoading] = useState(false)
  const [newFileId, setNewFileId] = useState<string>()

  const handleMove = useCallback(async () => {
    if (!newFileId) {
      return
    }
    try {
      setLoading(true)
      await FileAPI.move(newFileId, { ids: selection })
      if (fileId === newFileId) {
        const { data: files } = await FileAPI.list(newFileId, {
          page: 1,
          size: FileAPI.DEFAULT_PAGE_SIZE,
          sortBy,
          sortOrder,
        })
        dispatch(filesUpdated(files))
      } else {
        dispatch(filesRemoved({ id: fileId, files: selection }))
      }
      dispatch(selectionUpdated([]))
      dispatch(moveModalDidClose())
    } finally {
      setLoading(false)
    }
  }, [newFileId, fileId, selection, sortBy, sortOrder, dispatch])

  return (
    <Modal
      isOpen={isModalOpen}
      onClose={() => dispatch(moveModalDidClose())}
      closeOnOverlayClick={false}
    >
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>Move {selection.length} Item(s) to…</ModalHeader>
        <ModalCloseButton />
        <ModalBody>
          <Browse onChange={(id) => setNewFileId(id)} />
        </ModalBody>
        <ModalFooter>
          <Button
            type="button"
            variant="outline"
            colorScheme="blue"
            mr={variables.spacingSm}
            disabled={loading}
            onClick={() => dispatch(moveModalDidClose())}
          >
            Cancel
          </Button>
          <Button
            variant="solid"
            colorScheme="blue"
            isDisabled={newFileId === fileId}
            isLoading={loading}
            onClick={handleMove}
          >
            Move Here
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  )
}

export default Move
