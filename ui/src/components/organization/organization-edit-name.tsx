import { useCallback, useEffect, useRef, useState } from 'react'
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
import useFocusAndSelectAll from '@/hooks/use-focus-and-select-all'

export type OrganizationEditNameProps = {
  open: boolean
  organization: Organization
  onClose?: () => void
}

type FormValues = {
  name: string
}

const OrganizationEditName = ({
  open,
  organization,
  onClose,
}: OrganizationEditNameProps) => {
  const { mutate } = useSWRConfig()
  const [isModalOpen, setIsModalOpen] = useState(false)
  const formSchema = Yup.object().shape({
    name: Yup.string().required('Name is required').max(255),
  })
  const inputRef = useRef<HTMLInputElement>(null)
  useFocusAndSelectAll(inputRef, isModalOpen)

  useEffect(() => {
    setIsModalOpen(open)
  }, [open])

  const handleSubmit = useCallback(
    async (
      { name }: FormValues,
      { setSubmitting }: FormikHelpers<FormValues>,
    ) => {
      setSubmitting(true)
      try {
        const result = await OrganizationAPI.updateName(organization.id, {
          name,
        })
        mutate(`/organizations/${organization.id}`, result)
        setSubmitting(false)
        onClose?.()
      } finally {
        setSubmitting(false)
      }
    },
    [organization.id, onClose, mutate],
  )

  return (
    <Modal
      isOpen={isModalOpen}
      onClose={() => onClose?.()}
      closeOnOverlayClick={false}
    >
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>Edit Name</ModalHeader>
        <ModalCloseButton />
        <Formik
          enableReinitialize={true}
          initialValues={{ name: organization.name }}
          validationSchema={formSchema}
          validateOnBlur={false}
          onSubmit={handleSubmit}
        >
          {({ errors, touched, isSubmitting }) => (
            <Form>
              <ModalBody>
                <Field name="name">
                  {({ field }: FieldAttributes<FieldProps>) => (
                    <FormControl
                      isInvalid={errors.name && touched.name ? true : false}
                    >
                      <Input
                        ref={inputRef}
                        {...field}
                        placeholder="Name"
                        disabled={isSubmitting}
                        autoFocus
                      />
                      <FormErrorMessage>{errors.name}</FormErrorMessage>
                    </FormControl>
                  )}
                </Field>
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
                    colorScheme="blue"
                    isLoading={isSubmitting}
                  >
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

export default OrganizationEditName
