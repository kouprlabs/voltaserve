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
import { listUpdated } from '@/store/entities/files'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { copyModalDidClose, selectedItemsUpdated } from '@/store/ui/files'
import Browse from './browse'

const Copy = () => {
  const { mutate } = useSWRConfig()
  const dispatch = useAppDispatch()
  const { fileId } = useParams()
  const isModalOpen = useAppSelector((state) => state.ui.files.isCopyModalOpen)
  const selectedItems = useAppSelector((state) => state.ui.files.selectedItems)
  const [loading, setLoading] = useState(false)
  const [targetId, setTargetId] = useState<string>()

  const handleMove = useCallback(async () => {
    if (!targetId) {
      return
    }
    try {
      setLoading(true)
      await FileAPI.copy(targetId, {
        ids: selectedItems,
      })
      if (fileId === targetId) {
        const list = await mutate<List>(`/files/${targetId}/list`)
        if (list) {
          dispatch(listUpdated(list))
        }
      }
      dispatch(selectedItemsUpdated([]))
      dispatch(copyModalDidClose())
    } finally {
      setLoading(false)
    }
  }, [targetId, fileId, selectedItems, mutate, dispatch])

  return (
    <Modal
      isOpen={isModalOpen}
      onClose={() => dispatch(copyModalDidClose())}
      closeOnOverlayClick={false}
    >
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>Copy {selectedItems.length} Item(s) to…</ModalHeader>
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
