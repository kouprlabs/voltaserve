// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { useCallback } from 'react'
import { useNavigate, useParams, useSearchParams } from 'react-router-dom'
import {
  Button,
  FormControl,
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
import { decodeFileQuery, encodeFileQuery } from '@/lib/helpers/query'
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

const typeOptions = [
  { value: '', label: 'Any' },
  { value: 'file', label: 'File' },
  { value: 'folder', label: 'Folder' },
]

const FileSearchFilter = () => {
  const navigate = useNavigate()
  const dispatch = useAppDispatch()
  const { id: workspaceId, fileId } = useParams()
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
        const encodedQuery = encodeFileQuery({
          text: query?.text || '',
          type: type || undefined,
          createTimeBefore: createTimeBefore
            ? new Date(createTimeBefore).getTime()
            : undefined,
          createTimeAfter: createTimeAfter
            ? new Date(createTimeAfter).getTime()
            : undefined,
          updateTimeBefore: updateTimeBefore
            ? new Date(updateTimeBefore).getTime()
            : undefined,
          updateTimeAfter: updateTimeAfter
            ? new Date(updateTimeAfter).getTime()
            : undefined,
        })
        navigate(`/workspace/${workspaceId}/file/${fileId}?q=${encodedQuery}`)
        await mutateList?.()
        dispatch(modalDidClose())
        setSubmitting(false)
      } finally {
        setSubmitting(false)
      }
    },
    [query, workspaceId, fileId, mutateList, dispatch],
  )

  const handleClear = useCallback(async () => {
    navigate(`/workspace/${workspaceId}/file/${fileId}`)
    await mutateList?.()
    dispatch(modalDidClose())
  }, [workspaceId, fileId, mutateList, navigate, dispatch])

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
              type: query?.type || '',
              createTimeAfter: query?.createTimeAfter
                ? new Date(query?.createTimeAfter).toISOString().slice(0, 16)
                : '',
              createTimeBefore: query?.createTimeBefore
                ? new Date(query?.createTimeBefore).toISOString().slice(0, 16)
                : '',
              updateTimeAfter: query?.updateTimeAfter
                ? new Date(query?.updateTimeAfter).toISOString().slice(0, 16)
                : '',
              updateTimeBefore: query?.updateTimeBefore
                ? new Date(query?.updateTimeBefore).toISOString().slice(0, 16)
                : '',
            } as FormValues
          }
          validationSchema={formSchema}
          validateOnBlur={false}
          onSubmit={handleSubmit}
        >
          {({ isSubmitting, setFieldValue, values }) => (
            <Form>
              <ModalBody>
                <div className={cx('flex', 'flex-col', 'gap-1.5')}>
                  <FormControl>
                    <FormLabel>Type</FormLabel>
                    <Select
                      options={typeOptions}
                      defaultValue={
                        query?.type
                          ? typeOptions.find(
                              (option) => option.value === query.type,
                            )
                          : undefined
                      }
                      selectedOptionStyle="check"
                      chakraStyles={reactSelectStyles()}
                      isDisabled={isSubmitting}
                      onChange={async (event) => {
                        if (event) {
                          await setFieldValue('type', event.value)
                        }
                      }}
                    />
                  </FormControl>
                  <FormControl>
                    <FormLabel>Create Time (After - Before)</FormLabel>
                    <div className={cx('flex', 'items-center', 'gap-1.5')}>
                      <Field name="createTimeAfter">
                        {({ field }: FieldAttributes<FieldProps>) => (
                          <Input
                            {...field}
                            type="datetime-local"
                            disabled={isSubmitting}
                          />
                        )}
                      </Field>
                      <Field name="createTimeBefore">
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
                    <FormLabel>Update Time (After - Before)</FormLabel>
                    <div className={cx('flex', 'items-center', 'gap-1.5')}>
                      <Field name="updateTimeAfter">
                        {({ field }: FieldAttributes<FieldProps>) => (
                          <Input
                            {...field}
                            type="datetime-local"
                            disabled={isSubmitting}
                          />
                        )}
                      </Field>
                      <Field name="updateTimeBefore">
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
                  <Button
                    variant="outline"
                    colorScheme="red"
                    isDisabled={
                      !values.type &&
                      !values.createTimeAfter &&
                      !values.createTimeBefore &&
                      !values.updateTimeAfter &&
                      !values.updateTimeBefore
                    }
                    onClick={handleClear}
                  >
                    Clear Filter
                  </Button>
                  <Button
                    type="submit"
                    variant="solid"
                    colorScheme="blue"
                    isDisabled={
                      !values.type &&
                      !values.createTimeAfter &&
                      !values.createTimeBefore &&
                      !values.updateTimeAfter &&
                      !values.updateTimeBefore
                    }
                  >
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

export default FileSearchFilter
