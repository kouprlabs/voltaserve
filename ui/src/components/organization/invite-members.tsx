import { useCallback, useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import {
  Button,
  FormControl,
  FormErrorMessage,
  FormHelperText,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  Text,
  Textarea,
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
import InvitationAPI from '@/api/invitation'

type OrganizationInviteMembersProps = {
  open: boolean
  id: string
  onClose?: () => void
}

type FormValues = {
  emails: string
}

const OrganizationInviteMembers = ({
  open,
  id,
  onClose,
}: OrganizationInviteMembersProps) => {
  const navigate = useNavigate()
  const { mutate } = useSWRConfig()
  const [isModalOpen, setIsModalOpen] = useState(false)
  const formSchema = Yup.object().shape({
    emails: Yup.string().required('Email(s) are required'),
  })

  useEffect(() => {
    setIsModalOpen(open)
  }, [open])

  const handleSubmit = useCallback(
    async (
      { emails }: FormValues,
      { setSubmitting }: FormikHelpers<FormValues>
    ) => {
      setSubmitting(true)
      try {
        await InvitationAPI.create({
          organizationId: id,
          emails: emails.split(',').map((e: string) => e.trim()),
        })
        mutate(
          `/invitations/get_outgoing?${new URLSearchParams({
            organization_id: id,
          })}`
        )
        navigate(`/organization/${id}/invitation`)
        setSubmitting(false)
        onClose?.()
      } finally {
        setSubmitting(false)
      }
    },
    [id, navigate, onClose, mutate]
  )

  return (
    <Modal isOpen={isModalOpen} onClose={() => onClose?.()} size="3xl">
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>Invite Members</ModalHeader>
        <ModalCloseButton />
        <Formik
          enableReinitialize={true}
          initialValues={{ emails: '' }}
          validationSchema={formSchema}
          validateOnBlur={false}
          onSubmit={handleSubmit}
        >
          {({ errors, touched, isSubmitting }) => (
            <Form>
              <ModalBody>
                <Field name="emails">
                  {({ field }: FieldAttributes<FieldProps>) => (
                    <FormControl
                      isInvalid={errors.emails && touched.emails ? true : false}
                    >
                      <Textarea
                        {...field}
                        placeholder="Comma separated emails"
                        disabled={isSubmitting}
                        h="120px"
                      />
                      <FormHelperText>
                        <Text fontSize={variables.bodyFontSize}>
                          Example: alice@example.com, david@example.com
                        </Text>
                      </FormHelperText>
                      <FormErrorMessage>{errors.emails}</FormErrorMessage>
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
                  Invite
                </Button>
              </ModalFooter>
            </Form>
          )}
        </Formik>
      </ModalContent>
    </Modal>
  )
}

export default OrganizationInviteMembers
