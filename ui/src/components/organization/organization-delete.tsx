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
  Text,
} from '@chakra-ui/react'
import { useSWRConfig } from 'swr'
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
  const { mutate } = useSWRConfig()
  const [isModalOpen, setIsModalOpen] = useState(false)
  const formSchema = Yup.object().shape({
    name: Yup.string()
      .required('Confirmation is required')
      .oneOf([organization.name], 'Invalid organization name'),
  })

  useEffect(() => {
    setIsModalOpen(open)
  }, [open])

  const handleSubmit = useCallback(
    async (_: FormValues, { setSubmitting }: FormikHelpers<FormValues>) => {
      setSubmitting(true)
      try {
        await OrganizationAPI.delete(organization.id)

        navigate('/organization')
        mutate('/organizations')
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
                  <Text>
                    Are you sure you would like to delete this organization?
                  </Text>
                  <Text>
                    Please type <b>{organization.name}</b> to confirm.
                  </Text>
                  <Field name="name">
                    {({ field }: FieldAttributes<FieldProps>) => (
                      <FormControl
                        isInvalid={errors.name && touched.name ? true : false}
                      >
                        <Input {...field} disabled={isSubmitting} />
                        <FormErrorMessage>{errors.name}</FormErrorMessage>
                      </FormControl>
                    )}
                  </Field>
                </div>
              </ModalBody>
              <ModalFooter>
                <Button
                  type="button"
                  variant="outline"
                  colorScheme="blue"
                  className={cx('mr-1')}
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
              </ModalFooter>
            </Form>
          )}
        </Formik>
      </ModalContent>
    </Modal>
  )
}

export default OrganizationDelete
