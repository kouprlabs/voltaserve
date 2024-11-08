// Copyright 2024 Mateusz Kaźmierczak.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { ReactElement, useCallback, useRef } from 'react'
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
import useFocusAndSelectAll from '@/hooks/use-focus-and-select-all'

interface ConsoleRenameModalProps {
  header: ReactElement
  currentName: string
  isOpen: boolean
  onClose: () => void
  onRequest: ConsoleRenameModalRequest
}

export type ConsoleRenameModalRequest = (name: string) => Promise<void>

type FormValues = {
  name: string
}

const ConsoleRenameModal = ({
  header,
  currentName,
  isOpen,
  onClose,
  onRequest,
}: ConsoleRenameModalProps) => {
  const inputRef = useRef<HTMLInputElement>(null)
  useFocusAndSelectAll(inputRef, isOpen)
  const formSchema = Yup.object().shape({
    name: Yup.string().required('Name is required').max(255),
  })

  const handleRequest = useCallback(
    async (
      { name }: FormValues,
      { setSubmitting }: FormikHelpers<FormValues>,
    ) => {
      setSubmitting(true)
      try {
        await onRequest(name)
        onClose()
      } finally {
        setSubmitting(false)
      }
    },
    [onRequest, onClose],
  )

  return (
    <Modal isOpen={isOpen} onClose={onClose}>
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>{header}</ModalHeader>
        <ModalCloseButton />
        <Formik
          enableReinitialize={true}
          initialValues={{ name: currentName }}
          validationSchema={formSchema}
          validateOnBlur={false}
          onSubmit={handleRequest}
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
                        disabled={isSubmitting}
                        autoFocus
                        required
                        placeholder="Name"
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
                    onClick={onClose}
                  >
                    Cancel
                  </Button>
                  <Button
                    type="submit"
                    variant="solid"
                    colorScheme="blue"
                    isLoading={isSubmitting}
                  >
                    Confirm
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

export default ConsoleRenameModal
