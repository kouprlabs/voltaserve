// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { useCallback, useEffect, useState } from 'react'
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
import UserAPI, { User } from '@/client/idp/user'
import { getPictureUrl } from '@/lib/helpers/picture'
import { useAppSelector } from '@/store/hook'
import AccountUploadPicture from './account-upload-picture'

export type AccountEditPictureProps = {
  open: boolean
  user: User
  onClose?: () => void
}

type FormValues = {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  picture: any
}

const AccountEditPicture = ({
  open,
  user,
  onClose,
}: AccountEditPictureProps) => {
  const mutate = useAppSelector((state) => state.ui.account.mutate)
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [deletionInProgress, setDeletionInProgress] = useState(false)
  const formSchema = Yup.object().shape({
    picture: Yup.mixed()
      .required()
      .test(
        'fileSize',
        'Image is too big, should be less than 3 MB',
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        (value: any) => value === null || (value && value.size <= 3000000),
      )
      .test(
        'fileType',
        'Unsupported file format',
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        (value: any) =>
          value === null ||
          (value &&
            ['image/jpg', 'image/jpeg', 'image/gif', 'image/png'].includes(
              value.type,
            )),
      ),
  })

  useEffect(() => {
    setIsModalOpen(open)
  }, [open])

  const handleSubmit = useCallback(
    async (
      { picture }: FormValues,
      { setSubmitting }: FormikHelpers<FormValues>,
    ) => {
      setSubmitting(true)
      try {
        const result = await UserAPI.updatePicture(picture)
        await mutate?.(result)
        setSubmitting(false)
        onClose?.()
      } finally {
        setSubmitting(false)
      }
    },
    [onClose, mutate],
  )

  const handleDelete = useCallback(async () => {
    try {
      setDeletionInProgress(true)
      const result = await UserAPI.deletePicture()
      await mutate?.(result)
      onClose?.()
    } finally {
      setDeletionInProgress(false)
    }
  }, [onClose, mutate])

  return (
    <Modal
      isOpen={isModalOpen}
      onClose={() => onClose?.()}
      closeOnOverlayClick={false}
    >
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>Edit Picture</ModalHeader>
        <ModalCloseButton />
        <Formik
          enableReinitialize={true}
          initialValues={{
            picture: user.picture,
          }}
          validationSchema={formSchema}
          validateOnBlur={false}
          onSubmit={handleSubmit}
        >
          {({ errors, touched, isSubmitting, setFieldValue, values }) => (
            <Form>
              <ModalBody>
                <div
                  className={cx('flex', 'flex-col', 'items-center', 'gap-1')}
                >
                  <Field name="picture">
                    {({ field }: FieldAttributes<FieldProps>) => (
                      <FormControl
                        isInvalid={Boolean(errors.picture && touched.picture)}
                      >
                        <AccountUploadPicture
                          {...field}
                          initialValue={
                            user.picture
                              ? getPictureUrl(user.picture)
                              : undefined
                          }
                          disabled={isSubmitting}
                          onChange={async (event) => {
                            if (
                              event.target.files &&
                              event.target.files.length > 0
                            ) {
                              await setFieldValue(
                                'picture',
                                event.target.files[0],
                              )
                            }
                          }}
                        />
                        <FormErrorMessage>{errors.picture}</FormErrorMessage>
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
                    variant="outline"
                    colorScheme="red"
                    isLoading={deletionInProgress}
                    disabled={!user.picture}
                    onClick={handleDelete}
                  >
                    Delete
                  </Button>
                  <Button
                    type="submit"
                    variant="solid"
                    colorScheme="blue"
                    disabled={isSubmitting || values.picture === user.picture}
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

export default AccountEditPicture
