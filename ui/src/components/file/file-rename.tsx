// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { useCallback, useRef } from 'react'
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
import FileAPI from '@/client/api/file'
import useFocusAndSelectAll from '@/hooks/use-focus-and-select-all'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { renameModalDidClose } from '@/store/ui/files'

type FormValues = {
  name: string
}

const FileRename = () => {
  const dispatch = useAppDispatch()
  const { fileId } = useParams()
  const isModalOpen = useAppSelector(
    (state) => state.ui.files.isRenameModalOpen,
  )
  const id = useAppSelector((state) => state.ui.files.selection[0])
  const mutateList = useAppSelector((state) => state.ui.files.mutate)
  const { data: file, mutate: mutateFile } = FileAPI.useGet(id)
  const formSchema = Yup.object().shape({
    name: Yup.string().required('Name is required').max(255),
  })
  const inputRef = useRef<HTMLInputElement>(null)
  useFocusAndSelectAll(inputRef, isModalOpen)

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
        await mutateFile(await FileAPI.patchName(file.id, { name }))
        await mutateList?.()
        setSubmitting(false)
        dispatch(renameModalDidClose())
      } finally {
        setSubmitting(false)
      }
    },
    [file, fileId, dispatch, mutateFile, mutateList],
  )

  return (
    <Modal
      isOpen={isModalOpen}
      onClose={() => dispatch(renameModalDidClose())}
      closeOnOverlayClick={false}
    >
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>Rename Item</ModalHeader>
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
                      isInvalid={Boolean(errors.name && touched.name)}
                    >
                      <Input
                        ref={inputRef}
                        {...field}
                        placeholder="Name"
                        disabled={isSubmitting}
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
                </div>
              </ModalFooter>
            </Form>
          )}
        </Formik>
      </ModalContent>
    </Modal>
  )
}

export default FileRename
