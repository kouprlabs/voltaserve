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
import { copyModalDidClose, selectionUpdated } from '@/store/ui/files'
import FileBrowse from './file-browse'

const FileCopy = () => {
  const dispatch = useAppDispatch()
  const { fileId } = useParams()
  const isModalOpen = useAppSelector((state) => state.ui.files.isCopyModalOpen)
  const selection = useAppSelector((state) => state.ui.files.selection)
  const mutateList = useAppSelector((state) => state.ui.files.mutate)
  const [isLoading, setIsLoading] = useState(false)
  const [targetId, setTargetId] = useState<string>()

  const handleMove = useCallback(async () => {
    if (!targetId) {
      return
    }
    try {
      setIsLoading(true)
      await FileAPI.copy(targetId, {
        ids: selection,
      })
      if (fileId === targetId) {
        mutateList?.()
      }
      dispatch(selectionUpdated([]))
      dispatch(copyModalDidClose())
    } finally {
      setIsLoading(false)
    }
  }, [targetId, fileId, selection, dispatch, mutateList])

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
          <div className={cx('flex', 'flex-row', 'items-center', 'gap-1')}>
            <Button
              type="button"
              variant="outline"
              colorScheme="blue"
              disabled={isLoading}
              onClick={() => dispatch(copyModalDidClose())}
            >
              Cancel
            </Button>
            <Button
              variant="solid"
              colorScheme="blue"
              isLoading={isLoading}
              onClick={handleMove}
            >
              Copy Here
            </Button>
          </div>
        </ModalFooter>
      </ModalContent>
    </Modal>
  )
}

export default FileCopy
