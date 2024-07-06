// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

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
import { User } from '@/client/api/user'
import { useAppSelector } from '@/store/hook'
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
  const mutateList = useAppSelector((state) => state.ui.groupMembers.mutate)
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
        mutateList?.()
        setSubmitting(false)
        onClose?.()
        navigate(`/group/${group.id}/member`)
      } finally {
        setSubmitting(false)
      }
    },
    [group.id, navigate, onClose, mutateList],
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
