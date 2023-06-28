import React, { useCallback, useEffect, useMemo, useState } from 'react'
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
import { SectionSpinner, variables } from '@koupr/ui'
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
import StorageAPI from '@/client/api/storage'
import WorkspaceAPI, { Workspace } from '@/client/api/workspace'
import { swrConfig } from '@/client/options'
import StorageInput from '@/components/common/storage-input'

type EditStorageCapacityProps = {
  open: boolean
  workspace: Workspace
  onClose?: () => void
}

type FormValues = {
  storageCapacity: number
}

const EditStorageCapacity = ({
  open,
  workspace,
  onClose,
}: EditStorageCapacityProps) => {
  const { mutate } = useSWRConfig()
  const [isModalOpen, setIsModalOpen] = useState(false)
  const { data: storageUsage, error } = StorageAPI.useGetWorkspaceUsage(
    workspace.id,
    swrConfig()
  )
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

  useEffect(() => {
    setIsModalOpen(open)
  }, [open])

  const handleSubmit = useCallback(
    async (
      { storageCapacity }: FormValues,
      { setSubmitting }: FormikHelpers<FormValues>
    ) => {
      setSubmitting(true)
      try {
        const result = await WorkspaceAPI.updateStorageCapacity(workspace.id, {
          storageCapacity,
        })
        mutate(`/workspaces/${workspace.id}`, result)
        setSubmitting(false)
        onClose?.()
      } finally {
        setSubmitting(false)
      }
    },
    [workspace.id, onClose, mutate]
  )

  return (
    <Modal
      isOpen={isModalOpen}
      onClose={() => onClose?.()}
      closeOnOverlayClick={false}
    >
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>Change Storage Capacity</ModalHeader>
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
                {!storageUsage && !error ? (
                  <SectionSpinner height="100px" />
                ) : null}
                {storageUsage && !error ? (
                  <Field name="storageCapacity">
                    {(props: FieldAttributes<FieldProps>) => (
                      <FormControl
                        maxW="500px"
                        isInvalid={
                          errors.storageCapacity && touched.storageCapacity
                            ? true
                            : false
                        }
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
                <Button
                  type="button"
                  variant="outline"
                  colorScheme="blue"
                  mr={variables.spacingSm}
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
              </ModalFooter>
            </Form>
          )}
        </Formik>
      </ModalContent>
    </Modal>
  )
}

export default EditStorageCapacity
