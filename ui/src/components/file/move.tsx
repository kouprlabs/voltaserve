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
import { filesUpdated } from '@/store/entities/files'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { moveModalDidClose, selectionUpdated } from '@/store/ui/files'
import Browse from './browse'

const Move = () => {
  const { mutate } = useSWRConfig()
  const { fileId } = useParams()
  const dispatch = useAppDispatch()
  const selection = useAppSelector((state) => state.ui.files.selection)
  const isModalOpen = useAppSelector((state) => state.ui.files.isMoveModalOpen)
  const [loading, setLoading] = useState(false)
  const [targetId, setTargetId] = useState<string>()
  const fileListSearchParams = useFileListSearchParams()

  const handleMove = useCallback(async () => {
    if (!targetId) {
      return
    }
    try {
      setLoading(true)
      await FileAPI.move(targetId, { ids: selection })
      const list = await mutate<List>(
        `/files/${fileId}/list?${fileListSearchParams}`,
      )
      if (list) {
        dispatch(filesUpdated(list.data))
      }
      dispatch(selectionUpdated([]))
      dispatch(moveModalDidClose())
    } finally {
      setLoading(false)
    }
  }, [targetId, fileId, selection, fileListSearchParams, mutate, dispatch])

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
          <Browse onChange={(id) => setTargetId(id)} />
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
            isDisabled={targetId === fileId}
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
