// Copyright 2024 Mateusz Kaźmierczak.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { useEffect } from 'react'
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

interface ConsoleConfirmationModalProps {
  action: string | undefined
  target: string
  verbose: string
  closeConfirmationWindow: () => void
  isOpen: boolean
  isSubmitting: boolean
  request: (
    id: string | null,
    target: string | null,
    action: boolean | null,
    confirm: boolean,
  ) => Promise<void>
}

const ConsoleConfirmationModal = (props: ConsoleConfirmationModalProps) => {
  useEffect(() => {
    if (
      props.isOpen &&
      (props.action === undefined || props.target == undefined)
    ) {
      setTimeout(() => {
        window.location.reload()
      }, 2000)
      throw new Error('No action or target provided')
    }
  }, [props.isOpen])
  return (
    <>
      <Modal
        isOpen={props.isOpen}
        onClose={() => {
          props.closeConfirmationWindow()
        }}
      >
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>{props.action}</ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            {props.verbose}
            <b>{props.target}</b>
          </ModalBody>
          <ModalFooter>
            <div className={cx('flex', 'flex-row', 'items-center', 'gap-1')}>
              <Button
                type="button"
                variant="outline"
                colorScheme="blue"
                disabled={props.isSubmitting}
                onClick={() => {
                  props.closeConfirmationWindow()
                }}
              >
                Cancel
              </Button>
              <Button
                type="button"
                variant="solid"
                colorScheme="blue"
                isLoading={props.isSubmitting}
                onClick={async () => {
                  await props.request(null, null, null, true)
                }}
              >
                Confirm
              </Button>
            </div>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </>
  )
}

export default ConsoleConfirmationModal
