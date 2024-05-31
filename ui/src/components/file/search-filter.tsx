import { useCallback } from 'react'
import { useSearchParams } from 'react-router-dom'
import {
  Button,
  FormControl,
  FormErrorMessage,
  FormLabel,
  Input,
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
import { Select } from 'chakra-react-select'
import cx from 'classnames'
import { FileType } from '@/client/api/file'
import { decodeFileQuery } from '@/helpers/query'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { modalDidClose } from '@/store/ui/search-filter'
import { reactSelectStyles } from '@/styles/react-select'

type FormValues = {
  type: FileType | ''
  createTimeAfter: string
  createTimeBefore: string
  updateTimeAfter: string
  updateTimeBefore: string
}

const SearchFilter = () => {
  const dispatch = useAppDispatch()
  const [searchParams] = useSearchParams()
  const query = decodeFileQuery(searchParams.get('q') as string)
  const isModalOpen = useAppSelector(
    (state) => state.ui.searchFilter.isModalOpen,
  )
  const mutateList = useAppSelector((state) => state.ui.files.mutate)
  const formSchema = Yup.object().shape({
    type: Yup.string(),
    createTimeAfter: Yup.string(),
    createTimeBefore: Yup.string(),
    updateTimeAfter: Yup.string(),
    updateTimeBefore: Yup.string(),
  })

  const handleSubmit = useCallback(
    async (
      {
        type,
        createTimeBefore,
        createTimeAfter,
        updateTimeBefore,
        updateTimeAfter,
      }: FormValues,
      { setSubmitting }: FormikHelpers<FormValues>,
    ) => {
      setSubmitting(true)
      try {
        // TODO: Implement search filter
        console.log('type', type)
        console.log('createTimeBefore', new Date(createTimeBefore).getTime())
        console.log('createTimeAfter', new Date(createTimeAfter).getTime())
        console.log('updateTimeBefore', new Date(updateTimeBefore).getTime())
        console.log('updateTimeAfter', new Date(updateTimeAfter).getTime())
        mutateList?.()
        setSubmitting(false)
      } finally {
        setSubmitting(false)
      }
    },
    [mutateList],
  )

  const handleClose = useCallback(() => {
    dispatch(modalDidClose())
  }, [dispatch])

  return (
    <Modal
      size="xl"
      isOpen={isModalOpen}
      onClose={handleClose}
      closeOnOverlayClick={false}
    >
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>Search Filter</ModalHeader>
        <ModalCloseButton />
        <Formik
          enableReinitialize={true}
          initialValues={
            {
              type: '',
              createTimeAfter: '',
              createTimeBefore: '',
              updateTimeAfter: '',
              updateTimeBefore: '',
            } as FormValues
          }
          validationSchema={formSchema}
          validateOnBlur={false}
          onSubmit={handleSubmit}
        >
          {({ isSubmitting, setFieldValue }) => (
            <Form>
              <ModalBody>
                <div className={cx('flex', 'flex-col', 'gap-1.5')}>
                  <FormControl>
                    <FormLabel>Type</FormLabel>
                    <Select
                      options={[
                        { value: 'file', label: 'File' },
                        { value: 'folder', label: 'Folder' },
                      ]}
                      selectedOptionStyle="check"
                      chakraStyles={reactSelectStyles()}
                      isDisabled={isSubmitting}
                      onChange={(event) => {
                        if (event) {
                          setFieldValue('type', event.value)
                        }
                      }}
                    />
                  </FormControl>
                  <FormControl>
                    <FormLabel>Create Time (Before - After)</FormLabel>
                    <div className={cx('flex', 'items-center', 'gap-1.5')}>
                      <Field name="createTimeBefore">
                        {({ field }: FieldAttributes<FieldProps>) => (
                          <Input
                            {...field}
                            type="datetime-local"
                            disabled={isSubmitting}
                          />
                        )}
                      </Field>
                      <Field name="createTimeAfter">
                        {({ field }: FieldAttributes<FieldProps>) => (
                          <Input
                            {...field}
                            type="datetime-local"
                            disabled={isSubmitting}
                          />
                        )}
                      </Field>
                    </div>
                  </FormControl>
                  <FormControl>
                    <FormLabel>Update Time (Before - After)</FormLabel>
                    <div className={cx('flex', 'items-center', 'gap-1.5')}>
                      <Field name="updateTimeBefore">
                        {({ field }: FieldAttributes<FieldProps>) => (
                          <Input
                            {...field}
                            type="datetime-local"
                            disabled={isSubmitting}
                          />
                        )}
                      </Field>
                      <Field name="updateTimeAfter">
                        {({ field }: FieldAttributes<FieldProps>) => (
                          <Input
                            {...field}
                            type="datetime-local"
                            disabled={isSubmitting}
                          />
                        )}
                      </Field>
                    </div>
                  </FormControl>
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
                    onClick={handleClose}
                  >
                    Close
                  </Button>
                  <Button variant="outline" colorScheme="red">
                    Clear Filter
                  </Button>
                  <Button type="submit" variant="solid" colorScheme="blue">
                    Apply Filter
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

export default SearchFilter
