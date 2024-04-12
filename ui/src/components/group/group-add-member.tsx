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
import GroupAPI, { Group } from '@/client/api/group'
import UserAPI, { User } from '@/client/api/user'
import UserSelector from '../common/user-selector'

export type GroupAddMemberProps = {
  open: boolean
  group: Group
  onClose?: () => void
}

type FormValues = {
  userId: string
}

const GroupAddMember = ({ group, open, onClose }: GroupAddMemberProps) => {
  const navigate = useNavigate()
  const { mutate } = useSWRConfig()
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [activeUser, setActiveUser] = useState<User>()
  const formSchema = Yup.object().shape({
    userId: Yup.string().required('User is required'),
  })

  useEffect(() => {
    setIsModalOpen(open)
  }, [open])

  const handleSubmit = useCallback(
    async (
      { userId }: FormValues,
      { setSubmitting }: FormikHelpers<FormValues>,
    ) => {
      setSubmitting(true)
      try {
        await GroupAPI.addMember(group.id, {
          userId,
        })
        mutate(
          `/users?${UserAPI.paramsFromListOptions({
            groupId: group.id,
            nonGroupMembersOnly: true,
          })}`,
        )
        mutate(`/groups/${group.id}/get_available_users`)
        setSubmitting(false)
        onClose?.()
        navigate(`/group/${group.id}/member`)
      } finally {
        setSubmitting(false)
      }
    },
    [group.id, navigate, onClose, mutate],
  )

  return (
    <Modal
      isOpen={isModalOpen}
      onClose={() => onClose?.()}
      closeOnOverlayClick={false}
    >
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
          {({ errors, touched, isSubmitting, setFieldValue }) => (
            <Form>
              <ModalBody>
                <div className={cx('flex', 'flex-col', 'gap-1.5')}>
                  <Field name="userId">
                    {({ field }: FieldAttributes<FieldProps>) => (
                      <FormControl
                        isInvalid={
                          errors.userId && touched.userId ? true : false
                        }
                      >
                        <UserSelector
                          value={activeUser}
                          organizationId={group.organization.id}
                          groupId={group.id}
                          nonGroupMembersOnly={true}
                          onConfirm={(value) => {
                            setActiveUser(value)
                            setFieldValue(field.name, value.id)
                          }}
                        />
                        <FormErrorMessage>{errors.userId}</FormErrorMessage>
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
                    colorScheme="blue"
                    disabled={isSubmitting}
                    isLoading={isSubmitting}
                  >
                    Add
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

export default GroupAddMember
