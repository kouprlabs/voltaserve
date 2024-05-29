import { useEffect } from 'react'
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
import { modalDidClose, mutateMetadataUpdated } from '@/store/ui/mosaic'
import PerformanceCreate from './performance-create'
import PerformanceOverview from './performance-overview'

const Mosaic = () => {
  const dispatch = useAppDispatch()
  const id = useAppSelector((state) =>
    state.ui.files.selection.length > 0
      ? state.ui.files.selection[0]
      : undefined,
  )
  const isLoading = useAppSelector(
    (state) =>
      state.ui.mosaic.isCreating ||
      state.ui.mosaic.isUpdating ||
      state.ui.mosaic.isDeleting,
  )
  const isModalOpen = useAppSelector((state) => state.ui.mosaic.isModalOpen)
  const { data: metadata, mutate: mutateMetadata } = MosaicAPI.useGetMetadata(
    id,
    swrConfig(),
  )

  useEffect(() => {
    if (id) {
      dispatch(mutateMetadataUpdated(mutateMetadata))
    }
  }, [mutateMetadata])

  return (
    <Modal
      size="xl"
      isOpen={isModalOpen}
      onClose={() => dispatch(modalDidClose())}
      closeOnOverlayClick={false}
      closeOnEsc={!isLoading}
    >
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>Performance</ModalHeader>
        <ModalCloseButton isDisabled={isLoading} />
        {!metadata ? <PerformanceCreate /> : null}
        {metadata ? <PerformanceOverview /> : null}
      </ModalContent>
    </Modal>
  )
}

export default Mosaic
