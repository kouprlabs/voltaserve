// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

import { useCallback, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import {
  Button,
  FormControl,
  FormErrorMessage,
  FormLabel,
  Heading,
  Input,
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
import { Helmet } from 'react-helmet-async'
import WorkspaceAPI from '@/client/api/workspace'
import OrganizationSelector from '@/components/common/organization-selector'
import StorageInput from '@/components/common/storage-input'
import { gigabyteToByte } from '@/lib/helpers/convert-storage'
import { useAppSelector } from '@/store/hook'

type FormValues = {
  name: string
  organizationId: string
  storageCapacity: number
}

const NewWorkspacePage = () => {
  const navigate = useNavigate()
  const mutate = useAppSelector((state) => state.ui.workspaces.mutate)
  const [isLoading, setIsLoading] = useState(false)
  const formSchema = Yup.object().shape({
    name: Yup.string().required('Name is required').max(255),
    organizationId: Yup.string().required('Organization is required'),
    storageCapacity: Yup.number()
      .required('Storage capacity is required')
      .positive()
      .integer()
      .min(1, 'Invalid storage usage value'),
  })

  const handleSubmit = useCallback(
    async (
      { name, organizationId, storageCapacity }: FormValues,
      { setSubmitting }: FormikHelpers<FormValues>,
    ) => {
      setSubmitting(true)
      setIsLoading(true)
      try {
        const result = await WorkspaceAPI.create({
          name,
          organizationId,
          storageCapacity,
        })
        mutate?.()
        setSubmitting(false)
        navigate(`/workspace/${result.id}/file/${result.rootId}`)
      } catch (error) {
        setIsLoading(false)
      } finally {
        setSubmitting(false)
      }
    },
    [navigate, mutate],
  )

  return (
    <>
      <Helmet>
        <title>New Workspace</title>
      </Helmet>
      <div className={cx('flex', 'flex-col', 'gap-3.5')}>
        <Heading className={cx('text-heading')}>New Workspace</Heading>
        <Formik
          enableReinitialize={true}
          initialValues={{
            name: '',
            organizationId: '',
            storageCapacity: gigabyteToByte(100),
          }}
          validationSchema={formSchema}
          validateOnBlur={false}
          onSubmit={handleSubmit}
        >
          {({ errors, touched, isSubmitting, setFieldValue }) => (
            <Form>
              <div className={cx('flex', 'flex-col', 'gap-3.5')}>
                <div className={cx('flex', 'flex-col', 'gap-1.5')}>
                  <Field name="name">
                    {({ field }: FieldAttributes<FieldProps>) => (
                      <FormControl
                        maxW="400px"
                        isInvalid={errors.name && touched.name ? true : false}
                      >
                        <FormLabel>Name</FormLabel>
                        <Input {...field} disabled={isSubmitting} autoFocus />
                        <FormErrorMessage>{errors.name}</FormErrorMessage>
                      </FormControl>
                    )}
                  </Field>
                  <Field name="organizationId">
                    {({ field }: FieldAttributes<FieldProps>) => (
                      <FormControl
                        maxW="400px"
                        isInvalid={
                          errors.organizationId && touched.organizationId
                            ? true
                            : false
                        }
                      >
                        <FormLabel>Organization</FormLabel>
                        <OrganizationSelector
                          onConfirm={(value) =>
                            setFieldValue(field.name, value.id)
                          }
                        />
                        <FormErrorMessage>
                          {errors.organizationId}
                        </FormErrorMessage>
                      </FormControl>
                    )}
                  </Field>
                  <Field name="storageCapacity">
                    {(props: FieldAttributes<FieldProps>) => (
                      <FormControl
                        maxW="400px"
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
                </div>
                <div
                  className={cx('flex', 'flex-row', 'items-center', 'gap-1')}
                >
                  <Button
                    type="submit"
                    variant="solid"
                    colorScheme="blue"
                    isDisabled={isSubmitting || isLoading}
                    isLoading={isSubmitting || isLoading}
                  >
                    Create
                  </Button>
                  <Button as={Link} to="/workspace" variant="solid">
                    Cancel
                  </Button>
                </div>
              </div>
            </Form>
          )}
        </Formik>
      </div>
    </>
  )
}

export default NewWorkspacePage
