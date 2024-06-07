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
} from '@chakra-ui/react'
import { cx } from '@emotion/css'
import GroupAPI, { Group } from '@/client/api/group'
import { User } from '@/client/idp/user'
import userToString from '@/lib/helpers/user-to-string'

export type GroupRemoveMemberProps = {
  group: Group
  user: User
  isOpen: boolean
  onClose?: () => void
  onCompleted?: () => void
}

const GroupRemoveMember = ({
  group,
  user,
  isOpen,
  onCompleted,
  onClose,
}: GroupRemoveMemberProps) => {
  const [isLoading, setIsLoading] = useState(false)

  const handleRemoveMember = useCallback(async () => {
    try {
      setIsLoading(true)
      await GroupAPI.removeMember(group.id, {
        userId: user.id,
      })
      onCompleted?.()
      onClose?.()
    } finally {
      setIsLoading(false)
    }
  }, [group, user, onCompleted, onClose])

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
          <div>
            Are you sure you would like to remove member{' '}
            <span className={cx('font-bold')}>{userToString(user)}</span> from
            group <span className={cx('font-bold')}>{group.name}</span>?
          </div>
        </ModalBody>
        <ModalFooter>
          <div className={cx('flex', 'flex-row', 'items-center', 'gap-1')}>
            <Button
              type="button"
              variant="outline"
              colorScheme="blue"
              disabled={isLoading}
              onClick={() => onClose?.()}
            >
              Cancel
            </Button>
            <Button
              type="submit"
              variant="solid"
              colorScheme="red"
              isLoading={isLoading}
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

export default GroupRemoveMember
