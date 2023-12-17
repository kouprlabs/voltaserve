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
  Stack,
  Text,
} from '@chakra-ui/react'
import { variables } from '@koupr/ui'
import {
  Field,
  FieldAttributes,
  FieldProps,
  Form,
  Formik,
  FormikHelpers,
} from 'formik'
import * as Yup from 'yup'
import UserAPI from '@/client/idp/user'

type DeleteProps = {
  open: boolean
  onClose?: () => void
}

type FormValues = {
  password: string
}

const Delete = ({ open, onClose }: DeleteProps) => {
  const navigate = useNavigate()
  const [isModalOpen, setIsModalOpen] = useState(false)
  const formSchema = Yup.object().shape({
    password: Yup.string().required('Password is required'),
  })

  useEffect(() => {
    setIsModalOpen(open)
  }, [open])

  const handleSubmit = useCallback(
    async (
      { password }: FormValues,
      { setSubmitting }: FormikHelpers<FormValues>,
    ) => {
      setSubmitting(true)
      try {
        await UserAPI.delete({ password })
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
        <ModalHeader>Delete Account</ModalHeader>
        <ModalCloseButton />
        <Formik
          initialValues={{ password: '' }}
          validationSchema={formSchema}
          validateOnBlur={false}
          onSubmit={handleSubmit}
        >
          {({ errors, touched, isSubmitting }) => (
            <Form>
              <ModalBody>
                <Stack spacing={variables.spacingSm}>
                  <Text>
                    Are you sure you would like to delete your account
                    permanently?
                  </Text>
                  <Text fontWeight="semibold">
                    Type your password to confirm:
                  </Text>
                  <Field name="password">
                    {({ field }: FieldAttributes<FieldProps>) => (
                      <FormControl
                        isInvalid={
                          errors.password && touched.password ? true : false
                        }
                      >
                        <Input
                          {...field}
                          type="password"
                          placeholder="Password"
                          disabled={isSubmitting}
                        />
                        <FormErrorMessage>{errors.password}</FormErrorMessage>
                      </FormControl>
                    )}
                  </Field>
                </Stack>
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

export default Delete
