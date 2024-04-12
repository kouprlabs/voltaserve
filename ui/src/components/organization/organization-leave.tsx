import { useCallback, useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import {
  Button,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  Text,
} from '@chakra-ui/react'
import cx from 'classnames'
import OrganizationAPI from '@/client/api/organization'

export type OrganizationLeaveProps = {
  open: boolean
  id: string
  onClose?: () => void
}

const OrganizationLeave = ({ open, id, onClose }: OrganizationLeaveProps) => {
  const navigate = useNavigate()
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [isLoading, setIsLoading] = useState(false)

  useEffect(() => {
    setIsModalOpen(open)
  }, [open])

  const handleConfirmation = useCallback(async () => {
    setIsLoading(true)
    try {
      await OrganizationAPI.leave(id)
      navigate('/organization')
      onClose?.()
    } finally {
      setIsLoading(false)
    }
  }, [id, navigate, onClose])

  return (
    <Modal
      isOpen={isModalOpen}
      onClose={() => onClose?.()}
      closeOnOverlayClick={false}
    >
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>Leave Organization</ModalHeader>
        <ModalCloseButton />
        <ModalBody>
          <Text>Are you sure you would like to leave this organization?</Text>
        </ModalBody>
        <ModalFooter>
          <Button
            type="button"
            variant="outline"
            colorScheme="blue"
            className={cx('mr-1')}
            disabled={isLoading}
            onClick={() => onClose?.()}
          >
            Cancel
          </Button>
          <Button
            type="submit"
            variant="solid"
            colorScheme="red"
            disabled={isLoading}
            isLoading={isLoading}
            onClick={() => handleConfirmation()}
          >
            Leave
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  )
}

export default OrganizationLeave
