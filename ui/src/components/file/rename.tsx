import React, { useCallback } from 'react'
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
import {
  Field,
  FieldAttributes,
  FieldProps,
  Form,
  Formik,
  FormikHelpers,
} from 'formik'
import * as Yup from 'yup'
import FileAPI from '@/client/api/file'
import { filesUpdated } from '@/store/entities/files'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { renameModalDidClose } from '@/store/ui/files'

type FormValues = {
  name: string
}

const Rename = () => {
  const dispatch = useAppDispatch()
  const isModalOpen = useAppSelector(
    (state) => state.ui.files.isRenameModalOpen
  )
  const id = useAppSelector((state) => state.ui.files.selection[0])
  const { data: file, mutate } = FileAPI.useGetById(id)
  const formSchema = Yup.object().shape({
    name: Yup.string().required('Name is required').max(255),
  })

  const handleSubmit = useCallback(
    async (
      { name }: FormValues,
      { setSubmitting }: FormikHelpers<FormValues>
    ) => {
      if (!file) {
        return
      }
      setSubmitting(true)
      try {
        const result = await FileAPI.rename(file.id, {
          name,
        })
        mutate()
        setSubmitting(false)
        dispatch(filesUpdated([result]))
        dispatch(renameModalDidClose())
      } finally {
        setSubmitting(false)
      }
    },
    [file, dispatch, mutate]
  )

  return (
    <Modal isOpen={isModalOpen} onClose={() => dispatch(renameModalDidClose())}>
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>Rename File</ModalHeader>
        <ModalCloseButton />
        <Formik
          enableReinitialize={true}
          initialValues={{ name: file?.name || '' }}
          validationSchema={formSchema}
          validateOnBlur={false}
          onSubmit={handleSubmit}
        >
          {({ errors, touched, isSubmitting }) => (
            <Form>
              <ModalBody>
                <Field name="name">
                  {({ field }: FieldAttributes<FieldProps>) => (
                    <FormControl
                      isInvalid={errors.name && touched.name ? true : false}
                    >
                      <Input
                        {...field}
                        placeholder="Name"
                        disabled={isSubmitting}
                        autoFocus
                      />
                      <FormErrorMessage>{errors.name}</FormErrorMessage>
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
                  onClick={() => dispatch(renameModalDidClose())}
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

export default Rename
