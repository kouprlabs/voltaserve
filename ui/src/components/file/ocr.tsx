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
import { filesUpdated } from '@/store/entities/files'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { manageOcrModalDidClose } from '@/store/ui/files'

type FormValues = {
  languageId: string
}

const OCR = () => {
  const dispatch = useAppDispatch()
  const isModalOpen = useAppSelector(
    (state) => state.ui.files.isManageOcrModalOpen
  )
  const id = useAppSelector((state) => state.ui.files.selection[0])
  const { data: file, mutate } = FileAPI.useGetById(id)
  const formSchema = Yup.object().shape({
    languageId: Yup.string().required('Language is required'),
  })
  const [isDeleting, setIsDeleting] = useState(false)

  const handleSubmit = useCallback(
    async (
      { languageId }: FormValues,
      { setSubmitting }: FormikHelpers<FormValues>
    ) => {
      if (!file) {
        return
      }
      setSubmitting(true)
      try {
        const result = await FileAPI.updateOcrLanguage(file.id, {
          id: languageId,
        })
        dispatch(filesUpdated([result]))
        mutate(result)
        setSubmitting(false)
        dispatch(manageOcrModalDidClose())
      } finally {
        setSubmitting(false)
      }
    },
    [file, mutate, dispatch]
  )

  const handleDelete = useCallback(async () => {
    if (!file) {
      return
    }
    setIsDeleting(true)
    try {
      const result = await FileAPI.deleteOcr(file.id)
      dispatch(filesUpdated([result]))
      mutate(result)
      setIsDeleting(false)
      dispatch(manageOcrModalDidClose())
    } finally {
      setIsDeleting(false)
    }
  }, [file, mutate, dispatch])

  return (
    <Modal
      isOpen={isModalOpen}
      onClose={() => dispatch(manageOcrModalDidClose())}
      closeOnOverlayClick={false}
    >
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>OCR</ModalHeader>
        <ModalCloseButton />
        <Formik
          enableReinitialize={true}
          initialValues={{ languageId: file?.ocr?.language || '' }}
          validationSchema={formSchema}
          validateOnBlur={false}
          onSubmit={handleSubmit}
        >
          {({
            errors,
            touched,
            isSubmitting,
            values,
            dirty,
            setFieldValue,
          }) => (
            <Form>
              <ModalBody>
                <Field name="languageId">
                  {({ field }: FieldAttributes<FieldProps>) => (
                    <FormControl
                      maxW="400px"
                      isInvalid={
                        errors.languageId && touched.languageId ? true : false
                      }
                    >
                      <FormLabel>Language</FormLabel>
                      <OcrLanguageSelector
                        buttonLabel="Select Language"
                        valueId={values.languageId}
                        isDisabled={isSubmitting || isDeleting}
                        onConfirm={(value) =>
                          setFieldValue(field.name, value.id)
                        }
                      />
                      <FormErrorMessage>{errors.languageId}</FormErrorMessage>
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
                  isDisabled={
                    isSubmitting || isDeleting || !file?.ocr?.language
                  }
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
                  isDisabled={!dirty}
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

export default OCR
