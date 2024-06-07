import { useEffect } from 'react'
import {
  Modal,
  ModalCloseButton,
  ModalContent,
  ModalHeader,
  ModalOverlay,
} from '@chakra-ui/react'
import FileAPI from '@/client/api/file'
import MosaicAPI from '@/client/api/mosaic'
import { swrConfig } from '@/client/options'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { mutateInfoUpdated } from '@/store/ui/mosaic'
import { modalDidClose } from '@/store/ui/mosaic'
import MosaicCreate from './mosaic-create'
import MosaicOverview from './mosaic-overview'

const Mosaic = () => {
  const dispatch = useAppDispatch()
  const id = useAppSelector((state) =>
    state.ui.files.selection.length > 0
      ? state.ui.files.selection[0]
      : undefined,
  )
  const isModalOpen = useAppSelector((state) => state.ui.mosaic.isModalOpen)
  const { data: info, mutate: mutateInfo } = MosaicAPI.useGetInfo(
    id,
    swrConfig(),
  )
  const { data: file } = FileAPI.useGet(id, swrConfig())

  useEffect(() => {
    if (file?.snapshot?.taskId) {
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
        <ModalHeader>Mosaic</ModalHeader>
        <ModalCloseButton />
        {info.isAvailable ? <MosaicOverview /> : <MosaicCreate />}
      </ModalContent>
    </Modal>
  )
}

export default Mosaic
