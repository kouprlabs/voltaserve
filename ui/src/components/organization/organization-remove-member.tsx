import { useCallback, useState } from 'react'
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
import OrganizationAPI, { Organization } from '@/client/api/organization'
import { User } from '@/client/idp/user'
import userToString from '@/helpers/user-to-string'

export type OrganizationRemoveMemberProps = {
  organization: Organization
  user: User
  isOpen: boolean
  onClose?: () => void
  onCompleted?: () => void
}

const OrganizationRemoveMember = ({
  organization,
  user,
  isOpen,
  onCompleted,
  onClose,
}: OrganizationRemoveMemberProps) => {
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
    <Modal
      isOpen={isOpen}
      onClose={() => onClose?.()}
      closeOnOverlayClick={false}
    >
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>Remove Member</ModalHeader>
        <ModalCloseButton />
        <ModalBody>
          <Text>
            Are you sure you would like to remove member{' '}
            <Text as="span" className={cx('font-bold')}>
              {userToString(user)}
            </Text>{' '}
            from organization{' '}
            <Text as="span" className={cx('font-bold')}>
              {organization.name}
            </Text>
            ?
          </Text>
        </ModalBody>
        <ModalFooter>
          <div className={cx('flex', 'flex-row', 'items-center', 'gap-1')}>
            <Button
              type="button"
              variant="outline"
              colorScheme="blue"
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
          </div>
        </ModalFooter>
      </ModalContent>
    </Modal>
  )
}

export default OrganizationRemoveMember
