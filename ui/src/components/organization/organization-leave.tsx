// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

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
          <span>Are you sure you would like to leave this organization?</span>
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
              disabled={isLoading}
              isLoading={isLoading}
              onClick={() => handleConfirmation()}
            >
              Leave
            </Button>
          </div>
        </ModalFooter>
      </ModalContent>
    </Modal>
  )
}

export default OrganizationLeave
