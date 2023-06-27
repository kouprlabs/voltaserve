import React, { useCallback, useState } from 'react'
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
import { variables } from '@koupr/ui'
import OrganizationAPI, { Organization } from '@/client/api/organization'
import { User } from '@/client/idp/user'

type RemoveMemberProps = {
  organization: Organization
  user: User
  isOpen: boolean
  onClose?: () => void
  onCompleted?: () => void
}

const RemoveMember = ({
  organization,
  user,
  isOpen,
  onCompleted,
  onClose,
}: RemoveMemberProps) => {
  const [loading, setLoading] = useState(false)

  const handleRemoveMember = useCallback(async () => {
    try {
      setLoading(true)
      await OrganizationAPI.removeMember(organization.id, {
        userId: user.id,
      })
      onCompleted?.()
      onClose?.()
    } finally {
      setLoading(false)
    }
  }, [organization, user, onClose, onCompleted])

  return (
    <Modal isOpen={isOpen} onClose={() => onClose?.()}>
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>Remove Member</ModalHeader>
        <ModalCloseButton />
        <ModalBody>
          <Text>
            Are you sure you would like to remove member{' '}
            <Text as="span" fontWeight="bold" whiteSpace="nowrap">
              {user.fullName}
            </Text>{' '}
            from organization{' '}
            <Text as="span" fontWeight="bold" whiteSpace="nowrap">
              {organization.name}
            </Text>
            ?
          </Text>
        </ModalBody>
        <ModalFooter>
          <Button
            type="button"
            variant="outline"
            colorScheme="blue"
            mr={variables.spacingSm}
            disabled={loading}
            onClick={() => onClose?.()}
          >
            Cancel
          </Button>
          <Button
            type="submit"
            variant="solid"
            colorScheme="red"
            isLoading={loading}
            onClick={handleRemoveMember}
          >
            Remove
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  )
}

export default RemoveMember
