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
import UserAPI, { User } from '@/client/idp/user'
import useFocusAndSelectAll from '@/hooks/use-focus-and-select-all'
import { useAppSelector } from '@/store/hook'

export type AccountEditEmailProps = {
  open: boolean
  user: User
  onClose?: () => void
}

type FormValues = {
  email: string
}

const AccountEditEmail = ({ open, user, onClose }: AccountEditEmailProps) => {
  const mutate = useAppSelector((state) => state.ui.account.mutate)
  const [isModalOpen, setIsModalOpen] = useState(false)
  const formSchema = Yup.object().shape({
    email: Yup.string()
      .required('Email is required')
      .email('Must be a valid email')
      .max(255),
  })
  const inputRef = useRef<HTMLInputElement>(null)
  useFocusAndSelectAll(inputRef, isModalOpen)

  useEffect(() => {
    setIsModalOpen(open)
  }, [open])

  const handleSubmit = useCallback(
    async (
      { email }: FormValues,
      { setSubmitting }: FormikHelpers<FormValues>,
    ) => {
      setSubmitting(true)
      try {
        const result = await UserAPI.updateEmailRequest({
          email,
        })
        mutate?.(result)
        setSubmitting(false)
        onClose?.()
      } finally {
        setSubmitting(false)
      }
    },
    [onClose, mutate],
  )

  return (
    <Modal
      isOpen={isModalOpen}
      onClose={() => onClose?.()}
      closeOnOverlayClick={false}
    >
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
                        ref={inputRef}
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

export default AccountEditEmail
