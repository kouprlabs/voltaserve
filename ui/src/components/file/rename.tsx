import React, { useCallback } from 'react'
import { useParams } from 'react-router-dom'
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
import FileAPI, { List } from '@/client/api/file'
import useFileListSearchParams from '@/hooks/use-file-list-params'
import { filesUpdated, listUpdated } from '@/store/entities/files'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { renameModalDidClose } from '@/store/ui/files'

type FormValues = {
  name: string
}

const Rename = () => {
  const { mutate } = useSWRConfig()
  const dispatch = useAppDispatch()
  const { fileId } = useParams()
  const isModalOpen = useAppSelector(
    (state) => state.ui.files.isRenameModalOpen,
  )
  const id = useAppSelector((state) => state.ui.files.selection[0])
  const { data: file } = FileAPI.useGetById(id)
  const formSchema = Yup.object().shape({
    name: Yup.string().required('Name is required').max(255),
  })
  const fileListSearchParams = useFileListSearchParams()

  const handleSubmit = useCallback(
    async (
      { name }: FormValues,
      { setSubmitting }: FormikHelpers<FormValues>,
    ) => {
      if (!file) {
        return
      }
      setSubmitting(true)
      try {
        const result = await FileAPI.rename(file.id, {
          name,
        })
        const list = await mutate<List>(
          `/files/${fileId}/list?${fileListSearchParams}`,
        )
        if (list) {
          dispatch(listUpdated(list))
        }
        setSubmitting(false)
        dispatch(filesUpdated([result]))
        dispatch(renameModalDidClose())
      } finally {
        setSubmitting(false)
      }
    },
    [file, fileId, fileListSearchParams, dispatch, mutate],
  )

  return (
    <Modal
      isOpen={isModalOpen}
      onClose={() => dispatch(renameModalDidClose())}
      closeOnOverlayClick={false}
    >
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
