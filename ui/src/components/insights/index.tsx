import { useEffect } from 'react'
import {
  Modal,
  ModalCloseButton,
  ModalContent,
  ModalHeader,
  ModalOverlay,
} from '@chakra-ui/react'
import InsightsAPI from '@/client/api/insights'
import { swrConfig } from '@/client/options'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { modalDidClose, mutateMetadataUpdated } from '@/store/ui/insights'
import InsightsCreate from './insights-create'
import InsightsOverview from './insights-overview'

const Insights = () => {
  const dispatch = useAppDispatch()
  const id = useAppSelector((state) =>
    state.ui.files.selection.length > 0
      ? state.ui.files.selection[0]
      : undefined,
  )
  const isModalOpen = useAppSelector((state) => state.ui.insights.isModalOpen)
  const { data: metadata, mutate: mutateMetadata } = InsightsAPI.useGetMetadata(
    id,
    swrConfig(),
  )

  useEffect(() => {
    if (mutateMetadata) {
      dispatch(mutateMetadataUpdated(mutateMetadata))
    }
  }, [mutateMetadata])

  return (
    <Modal
      size="xl"
      isOpen={isModalOpen}
      onClose={() => dispatch(modalDidClose())}
      closeOnOverlayClick={false}
    >
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>Insights</ModalHeader>
        <ModalCloseButton />
        {!metadata ? <InsightsCreate /> : null}
        {metadata ? <InsightsOverview /> : null}
      </ModalContent>
    </Modal>
  )
}

export default Insights
