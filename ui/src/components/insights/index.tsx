import { useEffect } from 'react'
import {
  Modal,
  ModalCloseButton,
  ModalContent,
  ModalHeader,
  ModalOverlay,
} from '@chakra-ui/react'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { modalDidClose, wizardDidComplete } from '@/store/ui/insights'
import InsightsOverview from './insights-overview'
import InsightsWizard from './insights-wizard'

const Insights = () => {
  const dispatch = useAppDispatch()
  const isModalOpen = useAppSelector((state) => state.ui.insights.isModalOpen)
  const isWizardComplete = useAppSelector(
    (state) => state.ui.insights.isWizardComplete,
  )

  useEffect(() => {
    dispatch(wizardDidComplete(false))
  }, [dispatch])

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
        {!isWizardComplete ? <InsightsWizard /> : null}
        {isWizardComplete ? <InsightsOverview /> : null}
      </ModalContent>
    </Modal>
  )
}

export default Insights
