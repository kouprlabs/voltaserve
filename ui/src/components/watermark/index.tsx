import { useEffect } from 'react'
import {
  Modal,
  ModalCloseButton,
  ModalContent,
  ModalHeader,
  ModalOverlay,
} from '@chakra-ui/react'
import FileAPI from '@/client/api/file'
import { swrConfig } from '@/client/options'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { modalDidClose, mutateFileUpdated } from '@/store/ui/watermark'
import WatermarkCreate from './watermark-create'
import WatermarkOverview from './watermark-overview'

const Watermark = () => {
  const dispatch = useAppDispatch()
  const id = useAppSelector((state) =>
    state.ui.files.selection.length > 0
      ? state.ui.files.selection[0]
      : undefined,
  )
  const isModalOpen = useAppSelector((state) => state.ui.watermark.isModalOpen)
  const { data: file, mutate: mutateFile } = FileAPI.useGet(id, swrConfig())

  useEffect(() => {
    if (mutateFile) {
      dispatch(mutateFileUpdated(mutateFile))
    }
  }, [mutateFile])

  return (
    <Modal
      size="xl"
      isOpen={isModalOpen}
      onClose={() => dispatch(modalDidClose())}
      closeOnOverlayClick={false}
    >
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>Watermark</ModalHeader>
        <ModalCloseButton />
        {!file?.snapshot?.watermark ? <WatermarkCreate /> : null}
        {file?.snapshot?.watermark ? <WatermarkOverview /> : null}
      </ModalContent>
    </Modal>
  )
}

export default Watermark
