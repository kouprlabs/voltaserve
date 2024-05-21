import { useEffect } from 'react'
import {
  Modal,
  ModalCloseButton,
  ModalContent,
  ModalHeader,
  ModalOverlay,
} from '@chakra-ui/react'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { modalDidClose, wizardDidComplete } from '@/store/ui/ai'
import AiOverview from './ai-overview'
import AiWizard from './ai-wizard'

const AiModal = () => {
  const dispatch = useAppDispatch()
  const isModalOpen = useAppSelector((state) => state.ui.ai.isModalOpen)
  const isWizardComplete = useAppSelector(
    (state) => state.ui.ai.isWizardComplete,
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
        {!isWizardComplete ? <AiWizard /> : null}
        {isWizardComplete ? <AiOverview /> : null}
      </ModalContent>
    </Modal>
  )
}

export default AiModal
