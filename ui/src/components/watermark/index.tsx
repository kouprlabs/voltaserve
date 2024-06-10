import { useEffect } from 'react'
import {
  Modal,
  ModalCloseButton,
  ModalContent,
  ModalHeader,
  ModalOverlay,
} from '@chakra-ui/react'
import FileAPI from '@/client/api/file'
import WatermarkAPI from '@/client/api/watermark'
import { swrConfig } from '@/client/options'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { modalDidClose, mutateInfoUpdated } from '@/store/ui/watermark'
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
  const { data: info, mutate: mutateInfo } = WatermarkAPI.useGetInfo(
    id,
    swrConfig(),
  )
  const { data: file } = FileAPI.useGet(id, swrConfig())

  useEffect(() => {
    if (file?.snapshot?.task?.isPending) {
      dispatch(modalDidClose())
    }
  }, [file])

  useEffect(() => {
    if (mutateInfo) {
      dispatch(mutateInfoUpdated(mutateInfo))
    }
  }, [mutateInfo])

  if (!info) {
    return null
  }

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
        {info.isAvailable ? <WatermarkOverview /> : <WatermarkCreate />}
      </ModalContent>
    </Modal>
  )
}

export default Watermark
