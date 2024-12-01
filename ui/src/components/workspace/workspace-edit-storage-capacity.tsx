// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { useCallback, useEffect, useMemo, useState } from 'react'
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
import { SectionError, SectionSpinner } from '@koupr/ui'
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
import StorageAPI from '@/client/api/storage'
import WorkspaceAPI, { Workspace } from '@/client/api/workspace'
import { errorToString } from '@/client/error'
import { swrConfig } from '@/client/options'
import StorageInput from '@/components/common/storage-input'
import { useAppSelector } from '@/store/hook'

export type WorkspaceEditStorageCapacityProps = {
  open: boolean
  workspace: Workspace
  onClose?: () => void
}

type FormValues = {
  storageCapacity: number
}

const WorkspaceEditStorageCapacity = ({
  open,
  workspace,
  onClose,
}: WorkspaceEditStorageCapacityProps) => {
  const mutate = useAppSelector((state) => state.ui.workspace.mutate)
  const [isModalOpen, setIsModalOpen] = useState(false)
  const {
    data: storageUsage,
    error: storageUsageError,
    isLoading: storageUsageIsLoading,
  } = StorageAPI.useGetWorkspaceUsage(workspace.id, swrConfig())
  const formSchema = useMemo(() => {
    if (storageUsage) {
      return Yup.object().shape({
        storageCapacity: Yup.number()
          .required('Storage capacity is required')
          .positive()
          .min(storageUsage.bytes, 'Insufficient storage capacity'),
      })
    } else {
      return null
    }
  }, [storageUsage])
  const storageUsageIsReady = storageUsage && !storageUsageError

  useEffect(() => {
    setIsModalOpen(open)
  }, [open])

  const handleSubmit = useCallback(
    async (
      { storageCapacity }: FormValues,
      { setSubmitting }: FormikHelpers<FormValues>,
    ) => {
      setSubmitting(true)
      try {
        const result = await WorkspaceAPI.patchStorageCapacity(workspace.id, {
          storageCapacity,
        })
        await mutate?.(result)
        setSubmitting(false)
        onClose?.()
      } finally {
        setSubmitting(false)
      }
    },
    [workspace.id, onClose, mutate],
  )

  return (
    <Modal
      isOpen={isModalOpen}
      onClose={() => onClose?.()}
      closeOnOverlayClick={false}
    >
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>Edit Storage Capacity</ModalHeader>
        <ModalCloseButton />
        <Formik
          enableReinitialize={true}
          initialValues={{
            storageCapacity: workspace.storageCapacity,
          }}
          validationSchema={formSchema}
          validateOnBlur={false}
          onSubmit={handleSubmit}
        >
          {({ errors, touched, isSubmitting }) => (
            <Form>
              <ModalBody>
                {storageUsageIsLoading ? <SectionSpinner /> : null}
                {storageUsageError ? (
                  <SectionError text={errorToString(storageUsageError)} />
                ) : null}
                {storageUsageIsReady ? (
                  <Field name="storageCapacity">
                    {(props: FieldAttributes<FieldProps>) => (
                      <FormControl
                        maxW="500px"
                        isInvalid={Boolean(
                          errors.storageCapacity && touched.storageCapacity,
                        )}
                      >
                        <FormLabel>Storage capacity</FormLabel>
                        <StorageInput {...props} />
                        <FormErrorMessage>
                          {errors.storageCapacity}
                        </FormErrorMessage>
                      </FormControl>
                    )}
                  </Field>
                ) : null}
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
                    onClick={() => onClose?.()}
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

export default WorkspaceEditStorageCapacity
