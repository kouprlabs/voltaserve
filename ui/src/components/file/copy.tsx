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
import FileAPI from '@/api/file'
import { listUpdated } from '@/store/entities/files'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { copyModalDidClose, selectionUpdated } from '@/store/ui/files'
import Browse from './browse'

const Copy = () => {
  const dispatch = useAppDispatch()
  const { fileId: fileIdQuery } = useParams()
  const isModalOpen = useAppSelector((state) => state.ui.files.isCopyModalOpen)
  const selection = useAppSelector((state) => state.ui.files.selection)
  const sortBy = useAppSelector((state) => state.ui.files.sortBy)
  const sortOrder = useAppSelector((state) => state.ui.files.sortOrder)
  const [loading, setLoading] = useState(false)
  const [fileId, setFileId] = useState<string>()

  const handleMove = useCallback(async () => {
    if (!fileId) {
      return
    }
    try {
      setLoading(true)
      await FileAPI.copy(fileId, {
        ids: selection,
      })
      if (fileIdQuery === fileId) {
        const result = await FileAPI.list(
          fileId,
          FileAPI.DEFAULT_PAGE_SIZE,
          1,
          undefined,
          sortBy,
          sortOrder
        )
        dispatch(listUpdated(result))
      }
      dispatch(selectionUpdated([]))
      dispatch(copyModalDidClose())
    } finally {
      setLoading(false)
    }
  }, [fileId, fileIdQuery, selection, sortBy, sortOrder, dispatch])

  return (
    <Modal isOpen={isModalOpen} onClose={() => dispatch(copyModalDidClose())}>
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>Copy {selection.length} Item(s) to…</ModalHeader>
        <ModalCloseButton />
        <ModalBody>
          <Browse onChange={(id) => setFileId(id)} />
        </ModalBody>
        <ModalFooter>
          <Button
            type="button"
            variant="outline"
            colorScheme="blue"
            mr={variables.spacingSm}
            disabled={loading}
            onClick={() => dispatch(copyModalDidClose())}
          >
            Cancel
          </Button>
          <Button
            variant="solid"
            colorScheme="blue"
            isLoading={loading}
            onClick={handleMove}
          >
            Copy Here
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  )
}

export default Copy
