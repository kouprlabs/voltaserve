import { useEffect } from 'react'
import {
  Modal,
  ModalCloseButton,
  ModalContent,
  ModalHeader,
  ModalOverlay,
} from '@chakra-ui/react'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { modalDidClose, wizardDidComplete } from '@/store/ui/analysis'
import AnalysisOverview from './analysis-overview'
import AnalysisWizard from './analysis-wizard'

const AnalysisModal = () => {
  const dispatch = useAppDispatch()
  const isModalOpen = useAppSelector((state) => state.ui.analysis.isModalOpen)
  const isWizardComplete = useAppSelector(
    (state) => state.ui.analysis.isWizardComplete,
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
        <ModalHeader>Analyze</ModalHeader>
        <ModalCloseButton />
        {!isWizardComplete ? <AnalysisWizard /> : null}
        {isWizardComplete ? <AnalysisOverview /> : null}
      </ModalContent>
    </Modal>
  )
}

export default AnalysisModal
