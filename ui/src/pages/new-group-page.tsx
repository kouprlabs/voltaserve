// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { useCallback, useState } from 'react'
import { Link, useNavigate, useSearchParams } from 'react-router-dom'
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
import GroupAPI from '@/client/api/group'
import OrganizationAPI from '@/client/api/organization'
import { swrConfig } from '@/client/options'
import OrganizationSelector from '@/components/common/organization-selector'
import { useAppSelector } from '@/store/hook'

type FormValues = {
  name: string
  organizationId: string
}

const NewGroupPage = () => {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const mutate = useAppSelector((state) => state.ui.groups.mutate)
  const [isLoading, setIsLoading] = useState(false)
  const formSchema = Yup.object().shape({
    name: Yup.string().required('Name is required').max(255),
    organizationId: Yup.string().required('Organization is required'),
  })
  const { data: organization } = OrganizationAPI.useGet(
    searchParams.get('org'),
    swrConfig(),
  )

  const handleSubmit = useCallback(
    async (
      { name, organizationId }: FormValues,
      { setSubmitting }: FormikHelpers<FormValues>,
    ) => {
      setSubmitting(true)
      setIsLoading(true)
      try {
        const result = await GroupAPI.create({
          name,
          organizationId,
        })
        await mutate?.()
        setSubmitting(false)
        navigate(`/group/${result.id}/member`)
      } catch {
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
        <title>New Group</title>
      </Helmet>
      <div className={cx('flex', 'flex-col', 'gap-3.5')}>
        <Heading className={cx('text-heading')}>New Group</Heading>
        <Formik
          enableReinitialize={true}
          initialValues={{ name: '', organizationId: organization?.id || '' }}
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
                        isInvalid={Boolean(errors.name && touched.name)}
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
                        isInvalid={Boolean(
                          errors.organizationId && touched.organizationId,
                        )}
                      >
                        <FormLabel>Organization</FormLabel>
                        <OrganizationSelector
                          defaultValue={organization}
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
                  <Button as={Link} to="/group" variant="solid">
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

export default NewGroupPage
