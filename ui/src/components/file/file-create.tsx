import { useCallback, useEffect, useState } from 'react'
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
import FileAPI, { List } from '@/client/api/file'
import useFileListSearchParams from '@/hooks/use-file-list-params'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { createModalDidClose } from '@/store/ui/files'

type FormValues = {
  name: string
}

const FileCreate = () => {
  const { mutate } = useSWRConfig()
  const { id, fileId } = useParams()
  const dispatch = useAppDispatch()
  const isModalOpen = useAppSelector(
    (state) => state.ui.files.isCreateModalOpen,
  )
  const [inputRef, setInputRef] = useState<HTMLInputElement | null>()
  const formSchema = Yup.object().shape({
    name: Yup.string().required('Name is required').max(255),
  })
  const fileListSearchParams = useFileListSearchParams()

  useEffect(() => {
    if (inputRef) {
      inputRef.select()
    }
  }, [inputRef])

  const handleSubmit = useCallback(
    async (
      { name }: FormValues,
      { setSubmitting }: FormikHelpers<FormValues>,
    ) => {
      setSubmitting(true)
      try {
        await FileAPI.createFolder({
          name,
          workspaceId: id!,
          parentId: fileId!,
        })
        await mutate<List>(`/files/${fileId}/list?${fileListSearchParams}`)
        setSubmitting(false)
        dispatch(createModalDidClose())
      } finally {
        setSubmitting(false)
      }
    },
    [id, fileId, fileListSearchParams, mutate, dispatch],
  )

  return (
    <>
      <Modal
        isOpen={isModalOpen}
        onClose={() => dispatch(createModalDidClose())}
        closeOnOverlayClick={false}
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
                  <div
                    className={cx('flex', 'flex-row', 'items-center', 'gap-1')}
                  >
                    <Button
                      type="button"
                      variant="outline"
                      colorScheme="blue"
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
                  </div>
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
