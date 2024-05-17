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
import cx from 'classnames'
import FileAPI from '@/client/api/file'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { moveModalDidClose, selectionUpdated } from '@/store/ui/files'
import FileBrowse from './file-browse'

const FileMove = () => {
  const { fileId } = useParams()
  const dispatch = useAppDispatch()
  const selection = useAppSelector((state) => state.ui.files.selection)
  const isModalOpen = useAppSelector((state) => state.ui.files.isMoveModalOpen)
  const mutateList = useAppSelector((state) => state.ui.files.mutate)
  const [isLoading, setIsLoading] = useState(false)
  const [targetId, setTargetId] = useState<string>()

  const handleMove = useCallback(async () => {
    if (!targetId) {
      return
    }
    try {
      setIsLoading(true)
      await FileAPI.move(targetId, { ids: selection })
      mutateList?.()
      dispatch(selectionUpdated([]))
      dispatch(moveModalDidClose())
    } finally {
      setIsLoading(false)
    }
  }, [targetId, fileId, selection, dispatch, mutateList])

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
          <FileBrowse onChange={(id) => setTargetId(id)} />
        </ModalBody>
        <ModalFooter>
          <div className={cx('flex', 'flex-row', 'items-center', 'gap-1')}>
            <Button
              type="button"
              variant="outline"
              colorScheme="blue"
              disabled={isLoading}
              onClick={() => dispatch(moveModalDidClose())}
            >
              Cancel
            </Button>
            <Button
              variant="solid"
              colorScheme="blue"
              isDisabled={targetId === fileId}
              isLoading={isLoading}
              onClick={handleMove}
            >
              Move Here
            </Button>
          </div>
        </ModalFooter>
      </ModalContent>
    </Modal>
  )
}

export default FileMove
