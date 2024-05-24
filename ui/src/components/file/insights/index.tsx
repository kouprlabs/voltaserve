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
import { modalDidClose, mutateSummaryUpdated } from '@/store/ui/insights'
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
  const { data: summary, mutate: mutateSummary } = InsightsAPI.useGetSummary(
    id,
    swrConfig(),
  )

  useEffect(() => {
    if (id) {
      dispatch(mutateSummaryUpdated(mutateSummary))
    }
  }, [mutateSummary])

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
        {!summary?.hasEntities ? <InsightsCreate /> : null}
        {summary?.hasEntities ? <InsightsOverview /> : null}
      </ModalContent>
    </Modal>
  )
}

export default Insights
