// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
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
import { Form, Formik, FormikHelpers } from 'formik'
import cx from 'classnames'
import { AuthUserAPI } from '@/client/idp/user'

export type AccountDeleteProps = {
  open: boolean
  onClose?: () => void
}

// eslint-disable-next-line @typescript-eslint/no-empty-object-type
type FormValues = {}

const AccountDelete = ({ open, onClose }: AccountDeleteProps) => {
  const navigate = useNavigate()
  const [isModalOpen, setIsModalOpen] = useState(false)

  useEffect(() => {
    setIsModalOpen(open)
  }, [open])

  const handleSubmit = useCallback(
    async (_: FormValues, { setSubmitting }: FormikHelpers<FormValues>) => {
      setSubmitting(true)
      try {
        await AuthUserAPI.delete()
        navigate('/sign-in')
        onClose?.()
      } finally {
        setSubmitting(false)
      }
    },
    [navigate, onClose],
  )

  return (
    <Modal
      isOpen={isModalOpen}
      onClose={() => onClose?.()}
      closeOnOverlayClick={false}
    >
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>Delete Account and Data</ModalHeader>
        <ModalCloseButton />
        <Formik
          initialValues={{}}
          validateOnBlur={false}
          onSubmit={handleSubmit}
        >
          {({ isSubmitting }) => (
            <Form>
              <ModalBody>
                <div className={cx('flex', 'flex-col', 'items-start', 'gap-1')}>
                  <span>
                    Are you sure you want to delete your account and data?
                  </span>
                </div>
              </ModalBody>
              <ModalFooter>
                <div
                  className={cx('flex', 'flex-row', 'items-center', 'gap-1')}
                >
                  <Button
                    type="button"
                    variant="outline"
                    colorScheme="blue"
                    disabled={isSubmitting}
                    onClick={() => onClose?.()}
                  >
                    Cancel
                  </Button>
                  <Button
                    type="submit"
                    variant="solid"
                    colorScheme="red"
                    isLoading={isSubmitting}
                  >
                    Delete
                  </Button>
                </div>
              </ModalFooter>
            </Form>
          )}
        </Formik>
      </ModalContent>
    </Modal>
  )
}

export default AccountDelete
