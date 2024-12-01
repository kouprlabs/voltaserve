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
import {
  Button,
  FormControl,
  FormErrorMessage,
  Input,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
} from '@chakra-ui/react'
import { Field, FieldAttributes, FieldProps, Form, Formik, FormikHelpers } from 'formik'
import * as Yup from 'yup'
import cx from 'classnames'
import UserAPI, { User } from '@/client/idp/user'
import { useAppSelector } from '@/store/hook'

export type AccountChangePasswordProps = {
  open: boolean
  user: User
  onClose?: () => void
}

type FormValues = {
  currentPassword: string
  newPassword: string
}

const AccountChangePassword = ({ open, onClose }: AccountChangePasswordProps) => {
  const mutate = useAppSelector((state) => state.ui.account.mutate)
  const [isModalOpen, setIsModalOpen] = useState(false)
  const initialValues: FormValues = { currentPassword: '', newPassword: '' }
  const formSchema = Yup.object().shape({
    currentPassword: Yup.string().required('Current password is required'),
    newPassword: Yup.string().required('New password is required'),
  })

  useEffect(() => {
    setIsModalOpen(open)
  }, [open])

  const handleSubmit = useCallback(
    async ({ currentPassword, newPassword }: FormValues, { setSubmitting }: FormikHelpers<FormValues>) => {
      setSubmitting(true)
      try {
        const result = await UserAPI.updatePassword({
          currentPassword,
          newPassword,
        })
        await mutate?.(result)
        setSubmitting(false)
        onClose?.()
      } finally {
        setSubmitting(false)
      }
    },
    [onClose, mutate],
  )

  return (
    <Modal isOpen={isModalOpen} onClose={() => onClose?.()} closeOnOverlayClick={false}>
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>Change Password</ModalHeader>
        <ModalCloseButton />
        <Formik
          initialValues={initialValues}
          validationSchema={formSchema}
          validateOnBlur={false}
          onSubmit={handleSubmit}
        >
          {({ errors, touched, isSubmitting }) => (
            <Form>
              <ModalBody>
                <div className={cx('flex', 'flex-col', 'gap-1.5')}>
                  <Field name="currentPassword">
                    {({ field }: FieldAttributes<FieldProps>) => (
                      <FormControl isInvalid={Boolean(errors.currentPassword && touched.currentPassword)}>
                        <Input {...field} type="password" placeholder="Current password" disabled={isSubmitting} />
                        <FormErrorMessage>{errors.currentPassword}</FormErrorMessage>
                      </FormControl>
                    )}
                  </Field>
                  <Field name="newPassword">
                    {({ field }: FieldAttributes<FieldProps>) => (
                      <FormControl isInvalid={Boolean(errors.newPassword && touched.newPassword)}>
                        <Input {...field} type="password" placeholder="New password" disabled={isSubmitting} />
                        <FormErrorMessage>{errors.newPassword}</FormErrorMessage>
                      </FormControl>
                    )}
                  </Field>
                </div>
              </ModalBody>
              <ModalFooter>
                <div className={cx('flex', 'flex-row', 'items-center', 'gap-1')}>
                  <Button
                    type="button"
                    variant="outline"
                    colorScheme="blue"
                    disabled={isSubmitting}
                    onClick={() => onClose?.()}
                  >
                    Cancel
                  </Button>
                  <Button type="submit" variant="solid" colorScheme="blue" isLoading={isSubmitting}>
                    Save
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

export default AccountChangePassword
