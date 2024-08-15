// Copyright 2024 Mateusz Kaźmierczak.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { useEffect, useRef } from 'react'
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
import { Field, FieldAttributes, FieldProps, Form, Formik } from 'formik'
import * as Yup from 'yup'

interface AdminRenameModalProps {
  closeConfirmationWindow: () => void
  isOpen: boolean
  isSubmitting: boolean
  previousName: string | undefined
  object: string
  formSchema: Yup.ObjectSchema<
    { name: string },
    Yup.AnyObject,
    { name: undefined }
  >
  request: (
    id: string | null,
    currentName: string | null,
    newName: string | null,
    confirm: boolean,
  ) => Promise<void>
}
const AdminRenameModal = (props: AdminRenameModalProps) => {
  const inputRef = useRef<HTMLInputElement>(null)
  useEffect(() => {
    if (
      props.isOpen &&
      (props.previousName === undefined || props.formSchema == undefined)
    ) {
      setTimeout(() => {
        window.location.reload()
      }, 2000)
      throw new Error('No action or target provided')
    }
  }, [props.isOpen])

  return (
    <Modal
      isOpen={props.isOpen}
      onClose={() => {
        props.closeConfirmationWindow()
      }}
    >
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>Rename {props.object}</ModalHeader>
        <ModalCloseButton />
        <Formik
          enableReinitialize={true}
          initialValues={{ name: '' }}
          validationSchema={props.formSchema}
          validateOnBlur={false}
          onSubmit={async (event) => {
            await props.request(null, null, event.name, true)
          }}
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
                        ref={inputRef}
                        {...field}
                        disabled={isSubmitting}
                        autoFocus
                        required
                        placeholder={props.previousName}
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
                  disabled={isSubmitting}
                  onClick={() => {
                    props.closeConfirmationWindow()
                  }}
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
              </ModalFooter>
            </Form>
          )}
        </Formik>
      </ModalContent>
    </Modal>
  )
}

export default AdminRenameModal
