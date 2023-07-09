import React, { useCallback, useState } from 'react'
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
import OcrLanguageSelector from '@/components/common/ocr-language-selector'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { manageOcrModalDidClose } from '@/store/ui/files'

type FormValues = {
  ocrLanguageId: string
}

const ManageOcr = () => {
  const dispatch = useAppDispatch()
  const isModalOpen = useAppSelector(
    (state) => state.ui.files.isManageOcrModalOpen
  )
  const id = useAppSelector((state) => state.ui.files.selection[0])
  const { data: file } = FileAPI.useGetById(id)
  const formSchema = Yup.object().shape({
    ocrLanguageId: Yup.string().required('OCR language is required'),
  })
  const [isDeleting, setIsDeleting] = useState(false)

  const handleSubmit = useCallback(
    async (
      { ocrLanguageId }: FormValues,
      { setSubmitting }: FormikHelpers<FormValues>
    ) => {
      if (!file) {
        return
      }
      setSubmitting(true)
      try {
        await FileAPI.updateOcrLanguage(file.id, { ocrLanguageId })
        setSubmitting(false)
        dispatch(manageOcrModalDidClose())
      } finally {
        setSubmitting(false)
      }
    },
    [file, dispatch]
  )

  const handleDelete = useCallback(() => {
    setIsDeleting(true)
    setTimeout(() => setIsDeleting(false), 1500)
  }, [])

  return (
    <Modal
      isOpen={isModalOpen}
      onClose={() => dispatch(manageOcrModalDidClose())}
      closeOnOverlayClick={false}
    >
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>Manage OCR</ModalHeader>
        <ModalCloseButton />
        <Formik
          enableReinitialize={true}
          initialValues={{ ocrLanguageId: file?.ocr?.language || '' }}
          validationSchema={formSchema}
          validateOnBlur={false}
          onSubmit={handleSubmit}
        >
          {({ errors, touched, isSubmitting, values, setFieldValue }) => (
            <Form>
              <ModalBody>
                <Field name="ocrLanguageId">
                  {({ field }: FieldAttributes<FieldProps>) => (
                    <FormControl
                      maxW="400px"
                      isInvalid={
                        errors.ocrLanguageId && touched.ocrLanguageId
                          ? true
                          : false
                      }
                    >
                      <FormLabel>OCR Language</FormLabel>
                      <OcrLanguageSelector
                        defaultId={file?.ocr?.language}
                        isDisabled={isSubmitting || isDeleting}
                        onConfirm={(value) =>
                          setFieldValue(field.name, value.id)
                        }
                      />
                      <FormErrorMessage>
                        {errors.ocrLanguageId}
                      </FormErrorMessage>
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
                  isDisabled={isSubmitting}
                  onClick={() => dispatch(manageOcrModalDidClose())}
                >
                  Cancel
                </Button>
                <Button
                  type="button"
                  variant="solid"
                  colorScheme="red"
                  isDisabled={isSubmitting || isDeleting}
                  isLoading={isDeleting}
                  mr={variables.spacingSm}
                  onClick={handleDelete}
                >
                  Delete
                </Button>
                <Button
                  type="submit"
                  variant="solid"
                  colorScheme="blue"
                  isDisabled={values.ocrLanguageId === file?.ocr?.language}
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

export default ManageOcr
