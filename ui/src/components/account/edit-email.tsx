import React, { useCallback, useEffect, useState } from 'react'
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
import { variables } from '@koupr/ui'
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
import UserAPI, { User } from '@/client/idp/user'

type EditEmailProps = {
  open: boolean
  user: User
  onClose?: () => void
}

type FormValues = {
  email: string
}

const EditEmail = ({ open, user, onClose }: EditEmailProps) => {
  const { mutate } = useSWRConfig()
  const [isModalOpen, setIsModalOpen] = useState(false)
  const formSchema = Yup.object().shape({
    email: Yup.string()
      .required('Email is required')
      .email('Must be a valid email')
      .max(255),
  })

  useEffect(() => {
    setIsModalOpen(open)
  }, [open])

  const handleSubmit = useCallback(
    async (
      { email }: FormValues,
      { setSubmitting }: FormikHelpers<FormValues>
    ) => {
      setSubmitting(true)
      try {
        const result = await UserAPI.updateEmailRequest({
          email,
        })
        mutate(`/user`, result)
        setSubmitting(false)
        onClose?.()
      } finally {
        setSubmitting(false)
      }
    },
    [onClose, mutate]
  )

  return (
    <Modal isOpen={isModalOpen} onClose={() => onClose?.()}>
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>Edit Email</ModalHeader>
        <ModalCloseButton />
        <Formik
          enableReinitialize={true}
          initialValues={{ email: user?.pendingEmail || user?.email }}
          validationSchema={formSchema}
          validateOnBlur={false}
          onSubmit={handleSubmit}
        >
          {({ errors, touched, isSubmitting }) => (
            <Form>
              <ModalBody>
                <Field name="email">
                  {({ field }: FieldAttributes<FieldProps>) => (
                    <FormControl
                      isInvalid={errors.email && touched.email ? true : false}
                    >
                      <Input
                        {...field}
                        placeholder="Email"
                        disabled={isSubmitting}
                        autoFocus
                      />
                      <FormErrorMessage>{errors.email}</FormErrorMessage>
                    </FormControl>
                  )}
                </Field>
              </ModalBody>
              <ModalFooter>
                <Button
                  type="button"
                  variant="outline"
                  colorScheme="blue"
                  mr={variables.spacingSm}
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
              </ModalFooter>
            </Form>
          )}
        </Formik>
      </ModalContent>
    </Modal>
  )
}

export default EditEmail
