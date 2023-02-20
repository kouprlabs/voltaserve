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
import FileAPI from '@/api/file'
import { filesRemoved, filesUpdated } from '@/store/entities/files'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { moveModalDidClose, selectionUpdated } from '@/store/ui/files'
import variables from '@/theme/variables'
import FileBrowse from './browse'

const FileMove = () => {
  const params = useParams()
  const fileId = params.fileId as string
  const dispatch = useAppDispatch()
  const selection = useAppSelector((state) => state.ui.files.selection)
  const isModalOpen = useAppSelector((state) => state.ui.files.isMoveModalOpen)
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
        const { data: files } = await FileAPI.list(
          newFileId,
          FileAPI.DEFAULT_PAGE_SIZE,
          1
        )
        dispatch(filesUpdated(files))
      } else {
        dispatch(filesRemoved({ id: fileId, files: selection }))
      }
      dispatch(selectionUpdated([]))
      dispatch(moveModalDidClose())
    } finally {
      setLoading(false)
    }
  }, [newFileId, fileId, selection, dispatch])

  return (
    <Modal isOpen={isModalOpen} onClose={() => dispatch(moveModalDidClose())}>
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>Move {selection.length} Item(s) to…</ModalHeader>
        <ModalCloseButton />
        <ModalBody>
          <FileBrowse onChange={(id) => setNewFileId(id)} />
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
            Move here
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  )
}

export default FileMove
