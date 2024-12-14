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
import {
  Field,
  FieldAttributes,
  FieldProps,
  Form,
  Formik,
  FormikHelpers,
} from 'formik'
import * as Yup from 'yup'
import cx from 'classnames'
import OrganizationAPI, { Organization } from '@/client/api/organization'
import { useAppSelector } from '@/store/hook'

export type OrganizationDeleteProps = {
  open: boolean
  organization: Organization
  onClose?: () => void
}

type FormValues = {
  name: string
}

const OrganizationDelete = ({
  open,
  organization,
  onClose,
}: OrganizationDeleteProps) => {
  const navigate = useNavigate()
  const mutate = useAppSelector((state) => state.ui.organizations.mutate)
  const [isModalOpen, setIsModalOpen] = useState(false)
  const formSchema = Yup.object().shape({
    name: Yup.string()
      .required('Confirmation is required.')
      .oneOf([organization.name], 'Invalid organization name.'),
  })

  useEffect(() => {
    setIsModalOpen(open)
  }, [open])

  const handleSubmit = useCallback(
    async (_: FormValues, { setSubmitting }: FormikHelpers<FormValues>) => {
      setSubmitting(true)
      try {
        navigate('/organization')
        await OrganizationAPI.delete(organization.id)
        await mutate?.()
        onClose?.()
      } finally {
        setSubmitting(false)
      }
    },
    [organization.id, navigate, mutate, onClose],
  )

  return (
    <Modal
      isOpen={isModalOpen}
      onClose={() => onClose?.()}
      closeOnOverlayClick={false}
    >
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>Delete Organization</ModalHeader>
        <ModalCloseButton />
        <Formik
          initialValues={{ name: '' }}
          validationSchema={formSchema}
          validateOnBlur={false}
          onSubmit={handleSubmit}
        >
          {({ errors, touched, isSubmitting }) => (
            <Form>
              <ModalBody>
                <div className={cx('flex', 'flex-col', 'gap-1.5')}>
                  <span>
                    Are you sure you want to delete this organization?
                  </span>
                  <span>
                    Please type <b>{organization.name}</b> to confirm.
                  </span>
                  <Field name="name">
                    {({ field }: FieldAttributes<FieldProps>) => (
                      <FormControl
                        isInvalid={Boolean(errors.name && touched.name)}
                      >
                        <Input {...field} disabled={isSubmitting} />
                        <FormErrorMessage>{errors.name}</FormErrorMessage>
                      </FormControl>
                    )}
                  </Field>
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
                    Delete Permanently
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

export default OrganizationDelete
