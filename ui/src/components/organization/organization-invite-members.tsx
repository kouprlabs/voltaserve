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
  FormLabel,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  Textarea,
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
import InvitationAPI from '@/client/api/invitation'
import EmailTokenizer from '@/components/common/email-tokenizer'
import parseEmailList from '@/lib/helpers/parse-email-list'
import { useAppSelector } from '@/store/hook'

export type OrganizationInviteMembersProps = {
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
  const mutateList = useAppSelector(
    (state) => state.ui.outgoingInvitations.mutate,
  )
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
      { setSubmitting }: FormikHelpers<FormValues>,
    ) => {
      setSubmitting(true)
      try {
        await InvitationAPI.create({
          organizationId: id,
          emails: parseEmailList(emails),
        })
        await mutateList?.()
        navigate(`/organization/${id}/invitation`)
        setSubmitting(false)
        onClose?.()
      } finally {
        setSubmitting(false)
      }
    },
    [id, navigate, onClose, mutateList],
  )

  return (
    <Modal
      isOpen={isModalOpen}
      onClose={() => onClose?.()}
      size="3xl"
      closeOnOverlayClick={false}
    >
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
          {({ values, errors, touched, isSubmitting }) => (
            <Form>
              <ModalBody>
                <div className={cx('flex', 'flex-col', 'gap-1.5')}>
                  <Field name="emails">
                    {({ field }: FieldAttributes<FieldProps>) => (
                      <FormControl
                        isInvalid={
                          errors.emails && touched.emails ? true : false
                        }
                      >
                        <FormLabel>Comma separated emails:</FormLabel>
                        <Textarea
                          {...field}
                          className={cx('h-[120px]')}
                          placeholder="alice@example.com, david@example.com"
                          disabled={isSubmitting}
                        />
                        <FormErrorMessage>{errors.emails}</FormErrorMessage>
                      </FormControl>
                    )}
                  </Field>
                  <EmailTokenizer value={values.emails} />
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
                    isLoading={isSubmitting}
                  >
                    Invite
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

export default OrganizationInviteMembers
