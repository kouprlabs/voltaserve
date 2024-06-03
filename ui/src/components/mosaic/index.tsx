import {
  Modal,
  ModalCloseButton,
  ModalContent,
  ModalHeader,
  ModalOverlay,
} from '@chakra-ui/react'
import MosaicAPI from '@/client/api/mosaic'
import { swrConfig } from '@/client/options'
import { useAppDispatch, useAppSelector } from '@/store/hook'
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
  const { data: metadata } = MosaicAPI.useGetMetadata(id, swrConfig())

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
        {!metadata ? <MosaicCreate /> : null}
        {metadata ? <MosaicOverview /> : null}
      </ModalContent>
    </Modal>
  )
}

export default Mosaic
