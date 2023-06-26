import { useCallback, useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import {
  Button,
  FormControl,
  FormErrorMessage,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  Select,
  Stack,
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
import GroupAPI, { Group } from '@/client/api/group'
import userToString from '@/helpers/user-to-string'

type AddMemberProps = {
  open: boolean
  group: Group
  onClose?: () => void
}

type FormValues = {
  userId: string
}

const AddMember = ({ group, open, onClose }: AddMemberProps) => {
  const navigate = useNavigate()
  const { mutate } = useSWRConfig()
  const [isModalOpen, setIsModalOpen] = useState(false)
  const { data: users } = GroupAPI.useGetAvailableUsers(group.id)
  const formSchema = Yup.object().shape({
    userId: Yup.string().required('User is required'),
  })

  useEffect(() => {
    setIsModalOpen(open)
  }, [open])

  const handleSubmit = useCallback(
    async (
      { userId }: FormValues,
      { setSubmitting }: FormikHelpers<FormValues>
    ) => {
      setSubmitting(true)
      try {
        await GroupAPI.addMember(group.id, {
          userId,
        })
        mutate(`/groups/${group.id}/get_members`)
        mutate(`/groups/${group.id}/get_available_users`)
        setSubmitting(false)
        onClose?.()
        navigate(`/group/${group.id}/member`)
      } finally {
        setSubmitting(false)
      }
    },
    [group.id, navigate, onClose, mutate]
  )

  return (
    <Modal isOpen={isModalOpen} onClose={() => onClose?.()}>
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>Add Member</ModalHeader>
        <ModalCloseButton />
        <Formik
          enableReinitialize={true}
          initialValues={{ userId: '' }}
          validationSchema={formSchema}
          validateOnBlur={false}
          onSubmit={handleSubmit}
        >
          {({ errors, touched, isSubmitting }) => (
            <Form>
              <ModalBody>
                <Stack spacing={variables.spacing}>
                  <Field name="userId">
                    {({ field }: FieldAttributes<FieldProps>) => (
                      <FormControl
                        isInvalid={
                          errors.userId && touched.userId ? true : false
                        }
                      >
                        <Select {...field} placeholder="Select a user">
                          {users?.map((u) => (
                            <option key={u.id} value={u.id}>
                              {userToString(u)}
                            </option>
                          ))}
                        </Select>
                        <FormErrorMessage>{errors.userId}</FormErrorMessage>
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
                  colorScheme="blue"
                  disabled={isSubmitting}
                  isLoading={isSubmitting}
                >
                  Add
                </Button>
              </ModalFooter>
            </Form>
          )}
        </Formik>
      </ModalContent>
    </Modal>
  )
}

export default AddMember
