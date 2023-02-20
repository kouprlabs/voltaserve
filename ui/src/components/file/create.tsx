import { useCallback, useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import {
  Field,
  FieldAttributes,
  FieldProps,
  Form,
  Formik,
  FormikHelpers,
} from 'formik'
import * as Yup from 'yup'
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
import FileAPI from '@/api/file'
import { filesAdded } from '@/store/entities/files'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { createModalDidClose } from '@/store/ui/files'
import variables from '@/theme/variables'

type FormValues = {
  name: string
}

const FileCreate = () => {
  const params = useParams()
  const workspaceId = params.id as string
  const fileId = params.fileId as string
  const dispatch = useAppDispatch()
  const isModalOpen = useAppSelector(
    (state) => state.ui.files.isCreateModalOpen
  )
  const [inputRef, setInputRef] = useState<HTMLInputElement | null>()
  const formSchema = Yup.object().shape({
    name: Yup.string().required('Name is required').max(255),
  })

  useEffect(() => {
    if (inputRef) {
      inputRef.select()
    }
  }, [inputRef])

  const handleSubmit = useCallback(
    async (
      { name }: FormValues,
      { setSubmitting }: FormikHelpers<FormValues>
    ) => {
      setSubmitting(true)
      try {
        const result = await FileAPI.createFolder({
          name,
          workspaceId: workspaceId,
          parentId: fileId,
        })
        dispatch(filesAdded({ id: fileId, files: [result] }))
        setSubmitting(false)
        dispatch(createModalDidClose())
      } finally {
        setSubmitting(false)
      }
    },
    [fileId, workspaceId, dispatch]
  )

  return (
    <>
      <Modal
        isOpen={isModalOpen}
        onClose={() => dispatch(createModalDidClose())}
      >
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>New Folder</ModalHeader>
          <ModalCloseButton />
          <Formik
            enableReinitialize={true}
            initialValues={{ name: '' }}
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
                          ref={(r) => setInputRef(r)}
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
                    onClick={() => dispatch(createModalDidClose())}
                  >
                    Cancel
                  </Button>
                  <Button
                    type="submit"
                    variant="solid"
                    colorScheme="blue"
                    isLoading={isSubmitting}
                  >
                    Create
                  </Button>
                </ModalFooter>
              </Form>
            )}
          </Formik>
        </ModalContent>
      </Modal>
    </>
  )
}

export default FileCreate
