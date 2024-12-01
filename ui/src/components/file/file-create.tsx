// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
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
import { Field, FieldAttributes, FieldProps, Form, Formik, FormikHelpers } from 'formik'
import * as Yup from 'yup'
import cx from 'classnames'
import FileAPI, { FileType } from '@/client/api/file'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { createModalDidClose } from '@/store/ui/files'

type FormValues = {
  name: string
}

const FileCreate = () => {
  const { id: workspaceId, fileId } = useParams()
  const dispatch = useAppDispatch()
  const isModalOpen = useAppSelector((state) => state.ui.files.isCreateModalOpen)
  const mutateList = useAppSelector((state) => state.ui.files.mutate)
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
    async ({ name }: FormValues, { setSubmitting }: FormikHelpers<FormValues>) => {
      setSubmitting(true)
      try {
        await FileAPI.create({
          type: FileType.Folder,
          name,
          workspaceId: workspaceId!,
          parentId: fileId!,
        })
        await mutateList?.()
        setSubmitting(false)
        dispatch(createModalDidClose())
      } finally {
        setSubmitting(false)
      }
    },
    [workspaceId, fileId, dispatch, mutateList],
  )

  return (
    <>
      <Modal isOpen={isModalOpen} onClose={() => dispatch(createModalDidClose())} closeOnOverlayClick={false}>
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
                      <FormControl isInvalid={Boolean(errors.name && touched.name)}>
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
                  <div className={cx('flex', 'flex-row', 'items-center', 'gap-1')}>
                    <Button
                      type="button"
                      variant="outline"
                      colorScheme="blue"
                      disabled={isSubmitting}
                      onClick={() => dispatch(createModalDidClose())}
                    >
                      Cancel
                    </Button>
                    <Button type="submit" variant="solid" colorScheme="blue" isLoading={isSubmitting}>
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
