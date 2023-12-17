import React, { useCallback, useEffect, useState } from 'react'
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
  VStack,
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
import ImageUpload from './image-upload'

type EditPictureProps = {
  open: boolean
  user: User
  onClose?: () => void
}

type FormValues = {
  picture: any
}

const EditPicture = ({ open, user, onClose }: EditPictureProps) => {
  const { mutate } = useSWRConfig()
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [deletionInProgress, setDeletionInProgress] = useState(false)
  const formSchema = Yup.object().shape({
    picture: Yup.mixed()
      .required()
      .test(
        'fileSize',
        'Image is too big, should be less than 3 MB',
        (value: any) => value === null || (value && value.size <= 3000000),
      )
      .test(
        'fileType',
        'Unsupported file format',
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
        mutate(`/user`, result)
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
      mutate(`/user`, result)
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
                <VStack spacing={variables.spacingSm}>
                  <Field name="picture">
                    {({ field }: FieldAttributes<FieldProps>) => (
                      <FormControl
                        isInvalid={
                          errors.picture && touched.picture ? true : false
                        }
                      >
                        <ImageUpload
                          {...field}
                          initialValue={user.picture}
                          disabled={isSubmitting}
                          onChange={(e) =>
                            setFieldValue('picture', e.target.files[0])
                          }
                        />
                        <FormErrorMessage>{errors.picture}</FormErrorMessage>
                      </FormControl>
                    )}
                  </Field>
                </VStack>
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
                  variant="outline"
                  colorScheme="red"
                  mr={variables.spacingSm}
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
              </ModalFooter>
            </Form>
          )}
        </Formik>
      </ModalContent>
    </Modal>
  )
}

export default EditPicture
