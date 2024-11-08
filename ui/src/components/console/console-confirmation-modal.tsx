// Copyright 2024 Mateusz Kaźmierczak.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { ReactElement, useCallback, useState } from 'react'
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

export type ConsoleConfirmationModalProps = {
  header: ReactElement
  body: ReactElement
  isDestructive?: boolean
  isOpen: boolean
  onClose: () => void
  onRequest: ConsoleConfirmationModalRequest
}

export type ConsoleConfirmationModalRequest = () => Promise<void>

const ConsoleConfirmationModal = ({
  header,
  body,
  isDestructive,
  isOpen,
  onClose,
  onRequest,
}: ConsoleConfirmationModalProps) => {
  const [isLoading, setIsLoading] = useState(false)

  const handleRequest = useCallback(async () => {
    setIsLoading(true)
    try {
      await onRequest()
      onClose()
    } finally {
      setIsLoading(false)
    }
  }, [onRequest, onClose])

  return (
    <Modal isOpen={isOpen} onClose={onClose}>
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>{header}</ModalHeader>
        <ModalCloseButton />
        <ModalBody>{body}</ModalBody>
        <ModalFooter>
          <div className={cx('flex', 'flex-row', 'items-center', 'gap-1')}>
            <Button
              type="button"
              variant="outline"
              colorScheme="blue"
              disabled={isLoading}
              onClick={onClose}
            >
              Cancel
            </Button>
            <Button
              type="button"
              variant="solid"
              colorScheme={isDestructive ? 'red' : 'blue'}
              isLoading={isLoading}
              onClick={handleRequest}
            >
              Confirm
            </Button>
          </div>
        </ModalFooter>
      </ModalContent>
    </Modal>
  )
}

export default ConsoleConfirmationModal
