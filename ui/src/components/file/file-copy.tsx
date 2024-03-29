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
import { useSWRConfig } from 'swr'
import FileAPI, { List } from '@/client/api/file'
import useFileListSearchParams from '@/hooks/use-file-list-params'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { copyModalDidClose, selectionUpdated } from '@/store/ui/files'
import FileBrowse from './file-browse'

const FileCopy = () => {
  const { mutate } = useSWRConfig()
  const dispatch = useAppDispatch()
  const { fileId } = useParams()
  const isModalOpen = useAppSelector((state) => state.ui.files.isCopyModalOpen)
  const selection = useAppSelector((state) => state.ui.files.selection)
  const [loading, setLoading] = useState(false)
  const [targetId, setTargetId] = useState<string>()
  const fileListSearchParams = useFileListSearchParams()

  const handleMove = useCallback(async () => {
    if (!targetId) {
      return
    }
    try {
      setLoading(true)
      await FileAPI.copy(targetId, {
        ids: selection,
      })
      if (fileId === targetId) {
        await mutate<List>(`/files/${targetId}/list?${fileListSearchParams}`)
      }
      dispatch(selectionUpdated([]))
      dispatch(copyModalDidClose())
    } finally {
      setLoading(false)
    }
  }, [targetId, fileId, selection, fileListSearchParams, mutate, dispatch])

  return (
    <Modal
      isOpen={isModalOpen}
      onClose={() => dispatch(copyModalDidClose())}
      closeOnOverlayClick={false}
    >
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>Copy {selection.length} Item(s) to…</ModalHeader>
        <ModalCloseButton />
        <ModalBody>
          <FileBrowse onChange={(id) => setTargetId(id)} />
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

export default FileCopy
